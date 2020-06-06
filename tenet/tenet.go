// Package Terraria net
package tenet

import (
	"fmt"
	"log"
	"net"

	"github.com/maciej-zuk/tecligo/common"
)

func (c *Connection) Connect() {
	var (
		sleepingRoutines     int = 1
		synchronizedRoutines int = 3
	)

	if common.Settings.MapRendering {
		sleepingRoutines++
		synchronizedRoutines++
	}

	if common.Settings.DiscordBot {
		synchronizedRoutines++
	}

	address := fmt.Sprintf("%s:%d", common.Settings.TerrariaHost, common.Settings.TerrariaPort)
	log.Println("Connection start to", address)
	var err error
	c.conn, err = net.Dial("tcp", address)
	if err != nil {
		log.Println("Err:", err)
		return
	}
	log.Println("Connected, start up")
	c.exitWakeUp = make(chan bool, sleepingRoutines)
	c.running = true
	c.spawned = false
	c.Password = common.Settings.TerrariaPass
	c.writeBuffer.Grow(1 << 16)
	c.Players = make([]Player, 1<<8)

	c.enterWg.Add(1)
	c.exitWg.Add(synchronizedRoutines)
	go signalHandler(c)
	go receiverRoutine(c)
	go hearthbeatRoutine(c)
	go websocketRoutine(c)
	if common.Settings.MapRendering {
		go mapRoutine(c)
	}
	if common.Settings.DiscordBot {
		go discordBotRoutine(c)
	}
	log.Println("Spawning")
	c.Send(1, TnString("Terraria228"))
	c.exitWg.Wait()
	log.Println("Exit (tenet)")
}

func (c *Connection) Send(msgType uint8, payloads ...Payload) {
	c.writeMutex.Lock()
	size := 3
	for _, payload := range payloads {
		size += payload.Size()
	}
	c.writeBuffer.Reset()
	c.writeBuffer.WriteByte(byte(size & 0xff))
	c.writeBuffer.WriteByte(byte((size >> 8) & 0xff))
	c.writeBuffer.WriteByte(byte(msgType))
	for _, payload := range payloads {
		payload.Serialize(&c.writeBuffer)
	}
	c.writeBuffer.WriteTo(c.conn)
	c.writeMutex.Unlock()
}

func (c *Connection) Close() {
	if c.running {
		var sleepingRoutines int = 1

		if common.Settings.MapRendering {
			sleepingRoutines++
		}

		log.Println("Closing")
		c.running = false
		c.conn.Close()
		c.enterWg.Done()
		for i := 0; i < sleepingRoutines; i++ {
			c.exitWakeUp <- true
		}
	}
}

func (c *Connection) broadcastWebsocket(msg interface{}) {
	c.hub.broadcast <- msg
}

func (c *Connection) broadcastChat(msg string, color TnColor) {
	c.hub.broadcast <- struct {
		Type  string  `json:"type"`
		Data  string  `json:"data"`
		Color TnColor `json:"color"`
	}{
		"chat",
		msg,
		color,
	}
}

func (c *Connection) notifyPlayerUpdate(slot TnByte) {
	c.broadcastWebsocket(c.getPlayerUpdateMessage(slot))
}

func (c *Connection) notifyWorldUpdate() {
	c.broadcastWebsocket(c.getWorldUpdateMessage())
}

func (c *Connection) notifyMapUpdate() {
	c.broadcastWebsocket(c.getMapUpdateMessage())
}

func (c *Connection) getPlayerUpdateMessage(slot TnByte) interface{} {
	return struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
		Slot TnByte      `json:"slot"`
	}{
		"playerUpdate",
		c.Players[slot],
		slot,
	}
}

func (c *Connection) getWorldUpdateMessage() interface{} {
	return struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		"worldUpdate",
		c.WorldInfo,
	}
}

func (c *Connection) getMapUpdateMessage() interface{} {
	return struct {
		Type string `json:"type"`
	}{
		"mapUpdate",
	}
}
