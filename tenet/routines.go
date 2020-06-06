package tenet

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"github.com/maciej-zuk/tecligo/common"
	"github.com/maciej-zuk/tecligo/temap"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func receiverRoutine(c *Connection) {
	buf := make([]byte, 3)
	readBuffer := make([]byte, 1<<16)
	c.dataReader = bytes.NewReader(readBuffer)
	for c.running {
		ln, err := c.conn.Read(buf)
		if err != nil {
			log.Println("receiverRoutine Read error", err)
			c.Close()
			break
		}
		if ln == 3 {
			var dataLen uint16 = uint16(buf[0])
			dataLen |= uint16(buf[1]) << 8
			var msgType uint8 = uint8(buf[2])
			dataLen -= 3
			ln, err := io.ReadFull(c.conn, readBuffer[:dataLen])
			c.dataReader.Reset(readBuffer)
			if err != nil {
				c.Close()
				break
			}
			if ln == int(dataLen) {
				handler, hasHandler := netMsgHandler[msgType]
				if hasHandler {
					handler(c)
				} else {
					log.Println("receiverRoutine Unknown msg", msgType)
				}
			} else {
				log.Println("receiverRoutine Malformed payload", ln, dataLen)
			}
		} else {
			log.Println("receiverRoutine Malformed preamble")
		}
	}
	c.enterWg.Wait()
	c.exitWg.Done()
}

// sleeping routine
func hearthbeatRoutine(c *Connection) {
	cnt := 0
	for c.running {
		c.WorldInfo.Time += 60
		cnt++
		if cnt == 10 && c.spawned {
			cnt = 0
			c.Send(
				13,
				c.Slot,
				TnByte(64),
				TnByte(16),
				TnByte(8),
				TnByte(0),
				TnByte(0),
				TnVector{16 * float32(c.WorldInfo.SpawnTileX), float32(16 * (c.WorldInfo.MaxTilesY - 100))},
			)
		}
		select {
		case <-c.exitWakeUp:
			break
		case <-time.After(1 * time.Second):
			break
		}
	}
	c.enterWg.Wait()
	c.exitWg.Done()
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	log.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	log.Printf("\tSys = %v MiB", bToMb(m.Sys))
	log.Printf("\tHeapInuse = %v MiB", bToMb(m.HeapInuse))
	log.Printf("\tStackInuse = %v MiB", bToMb(m.StackInuse))
	log.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func signalHandler(c *Connection) {
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	log.Println("Signal")
	c.Close()
}

func cutTiles(c *Connection, path string) {
	log.Println("Cutting tiles")
	w := temap.World{}
	w.Load(path)
	w.LoadTiles()
	w.RecalcUV()
	margin := 8
	_tw, _th := w.GetSize()
	tw := int(_tw)
	th := int(_th)
	var waitGroup sync.WaitGroup
cutTilesOuterLoop:
	for x := 0; x < tw; x += 32 {
		for y := 0; y < th; y += 32 {
			if !c.running {
				log.Println("Stoping cut")
				break cutTilesOuterLoop
			}
			func(x int, y int) {
				defer func() {
					if recover() != nil {
						log.Println("Panic during cutting", x, y)
					}
				}()
				x0 := x - margin
				y0 := y - margin
				x1 := x + 32 + margin
				y1 := y + 32 + margin
				if x0 < 0 {
					x0 = 0
				}
				if y0 < 0 {
					y0 = 0
				}
				if x1 > tw {
					x1 = tw
				}
				if y1 > th {
					y1 = th
				}
				xoff := (x - x0) * 16
				yoff := (y - y0) * 16
				target, _ := sdl.CreateRGBSurfaceWithFormat(0, 512, 512, 32, sdl.PIXELFORMAT_ABGR8888)
				ctx := temap.RenderingContext{uint32(x0), uint32(y0), uint32(x1), uint32(y1), uint32(xoff), uint32(yoff), target}
				w.RenderOnto(&ctx)
				waitGroup.Add(1)
				go func(x int, y int, target *sdl.Surface) {
					defer waitGroup.Done()
					img.SavePNG(target, fmt.Sprintf(common.Settings.TileOutputPath, x/32, y/32))
					target.Free()
				}(x, y, target)
			}(x, y)
		}
	}
	waitGroup.Wait()
	temap.ClearTextureCache()
	log.Println("Cutting tiles - done")
}

// sleeping routine
func mapRoutine(c *Connection) {
	path := common.Settings.MapPath
	getMTime := func() time.Time {
		stat, _ := os.Stat(path)
		return stat.ModTime()
	}
	lastMTime := getMTime()
	for c.running {
		currentMTime := getMTime()
		if currentMTime.After(lastMTime) {
			lastMTime = currentMTime
			cutTiles(c, path)
			if !c.running {
				break
			}
			c.notifyMapUpdate()
		}
		select {
		case <-c.exitWakeUp:
			break
		case <-time.After(1 * time.Hour):
			break
		}
	}
	c.enterWg.Wait()
	c.exitWg.Done()
}
