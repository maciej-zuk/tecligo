package temap

import (
	"strconv"
	"sync"
	"unsafe"

	"github.com/maciej-zuk/tecligo/common"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

// import (
//     "fmt"
//     "github.com/veandco/go-sdl2/img"
//     "os"
//     "path/filepath"
//     "strings"
// )

// func DumpAllTextures() {
//     filepath.Walk(basepath, func(path string, info os.FileInfo, err error) error {
//         if !info.IsDir() && strings.HasSuffix(path, ".xnb") {
//             basename := filepath.Base(path)
//             filename := strings.Split(basename, ".")[0]
//             t, err := load(filename)
//             if !err {
//                 img.SavePNG(t.tex, fmt.Sprintf("dump/%s.png", filename))
//                 t.tex.Free()
//             }
//         }
//         return nil
//     })
// }

const (
	Tex_WallOutline  int = 0x0000
	Tex_Shroom       int = 0x0001
	Tex_Banner       int = 0x0002
	Tex_Actuator     int = 0x0003
	Tex_Background   int = 0x10000
	Tex_Underworld   int = 0x20000
	Tex_Wall         int = 0x30000
	Tex_Tile         int = 0x40000
	Tex_Liquid       int = 0x50000
	Tex_NPC          int = 0x60000
	Tex_NPCHead      int = 0x70000
	Tex_ArmorHead    int = 0x80000
	Tex_ArmorBody    int = 0x90000
	Tex_ArmorFemale  int = 0xa0000
	Tex_ArmorLegs    int = 0xb0000
	Tex_TreeTops     int = 0xc0000
	Tex_TreeBranches int = 0xd0000
	Tex_Xmas         int = 0xe0000
	Tex_Wood         int = 0xf0000
	Tex_Cactus       int = 0x100000
	Tex_Wire         int = 0x110000
	Tex_Item         int = 0x120000
	Tex_Player0      int = 0x130000
	Tex_PlayerHair   int = 0x140000
)

type texData struct {
	tex  *sdl.Surface
	data []byte
}

var (
	textures      map[int]texData
	emptyTexture  *sdl.Surface
	helperTexture *sdl.Surface
	basepath      string
	cacheMux      sync.Mutex
)

func genHelperTexture() {
	helperTexture, _ = sdl.CreateRGBSurfaceWithFormat(0, 4*16, 1*16, 32, sdl.PIXELFORMAT_ABGR8888)
	renderer, _ := sdl.CreateSoftwareRenderer(helperTexture)
	defer renderer.Destroy()
	gfx.BoxColor(renderer, 0*16, 0, 1*16+15, 16*3+15, sdl.Color{7, 60, 190, 200})
	gfx.BoxColor(renderer, 1*16, 0, 2*16+15, 16*3+15, sdl.Color{202, 36, 15, 200})
	gfx.BoxColor(renderer, 2*16, 0, 3*16+15, 16*3+15, sdl.Color{255, 149, 1, 200})
}

func InitTextures() {
	basepath = common.Settings.BasePath
	textures = make(map[int]texData)
	emptyTexture, _ = sdl.CreateRGBSurfaceWithFormat(0, 16, 16, 32, sdl.PIXELFORMAT_ABGR8888)
	genHelperTexture()
}

func DestroyTextures() {
	emptyTexture.Free()
	helperTexture.Free()
	ClearTextureCache()
}

func ClearTextureCache() {
	cacheMux.Lock()
	defer cacheMux.Unlock()

	for k, t := range textures {
		if t.tex != nil {
			t.tex.Free()
		}
		t.data = nil
		delete(textures, k)
	}
}

func GetTexture(ttype int) *sdl.Surface {
	return GetTextureCrop(ttype, 0, 0)
}

func GetTextureCrop(ttype int, cropw int, croph int) *sdl.Surface {
	cacheMux.Lock()
	defer cacheMux.Unlock()

	texture, hasTexture := textures[ttype]
	if hasTexture {
		return texture.tex
	}

	mask := ttype & 0xff0000
	num := ttype & 0x00ffff
	if mask == 0 {
		mask = ttype
	}
	var name string
	switch mask {
	case Tex_WallOutline:
		name = "Wall_Outline"
	case Tex_Shroom:
		name = "Shroom_Tops"
	case Tex_Banner:
		name = "House_Banner_1"
	case Tex_Actuator:
		name = "Actuator"
	case Tex_Background:
		name = "Background_" + strconv.Itoa(num)
	case Tex_Underworld:
		name = "Backgrounds/Underworld " + strconv.Itoa(num)
	case Tex_Wall:
		name = "Wall_" + strconv.Itoa(num)
	case Tex_Tile:
		name = "Tiles_" + strconv.Itoa(num)
	case Tex_Liquid:
		name = "Liquid_" + strconv.Itoa(num)
	case Tex_NPC:
		name = "NPC_" + strconv.Itoa(num)
	case Tex_NPCHead:
		name = "NPC_Head_" + strconv.Itoa(num)
	case Tex_ArmorHead:
		name = "Armor_Head_" + strconv.Itoa(num)
	case Tex_ArmorBody:
		name = "Armor_Body_" + strconv.Itoa(num)
	case Tex_ArmorFemale:
		name = "Female_Body_" + strconv.Itoa(num)
	case Tex_ArmorLegs:
		name = "Armor_Legs_" + strconv.Itoa(num)
	case Tex_TreeTops:
		name = "Tree_Tops_" + strconv.Itoa(num)
	case Tex_TreeBranches:
		name = "Tree_Branches_" + strconv.Itoa(num)
	case Tex_Xmas:
		name = "Xmas_" + strconv.Itoa(num)
	case Tex_Wood:
		name = "Tiles_5_" + strconv.Itoa(num)
	case Tex_Cactus:
		name := "Tiles_80"
		switch num {
		case 1:
			name = "Evil_Cactus"
		case 2:
			name = "Good_Cactus"
		case 3:
			name = "Crimson_Cactus"
		}
		name = name
	case Tex_Wire:
		name = "WiresNew"
	case Tex_Item:
		name = "Item_" + strconv.Itoa(num)
	case Tex_Player0:
		name = "Player_0_" + strconv.Itoa(num)
	case Tex_PlayerHair:
		name = "Player_Hair_" + strconv.Itoa(num)
	}
	texture, err := load(name)
	if !err {
		if cropw > 0 && croph > 0 {
			tmp, _ := sdl.CreateRGBSurfaceWithFormat(0, int32(cropw), int32(croph), 32, sdl.PIXELFORMAT_ABGR8888)
			texture.tex.Blit(&texture.tex.ClipRect, tmp, &tmp.ClipRect)
			texture.tex.Free()
			texture.data = []byte{}
			texture.tex = tmp
		}
		textures[ttype] = texture
	}
	// img.SavePNG(texture.tex, name+".png")
	return texture.tex
}

func load(name string) (texData, bool) {
	handle := Handle{}
	defer handle.Close()
	if handle.LoadFile(basepath + name + ".xnb") {
		return texData{emptyTexture, []byte{}}, true
	}
	handle.U32()
	version := handle.U16()
	compressed := (version & 0x8000) != 0
	ln := handle.U32()
	var rawdata []byte
	if compressed {
		tmp := handle.Read(int(ln))
		rawdata = LZXDecompress(tmp)
	} else {
		rawdata = handle.Read(int(ln))
	}
	tex := Handle{}
	defer tex.Close()
	tex.LoadBytes(rawdata)
	numReaders := 0
	bits := 0
	for {
		b7 := int(tex.U8())
		numReaders |= (b7 & 0x7f) << bits
		bits += 7
		if b7&0x80 == 0 {
			break
		}
	}
	for x := 0; x < numReaders; x++ {
		tex.S()
		tex.U32()
	}
	for tex.U8()&0x80 != 0 {
	}
	for tex.U8()&0x80 != 0 {
	}
	tex.U32()
	width := int32(tex.U32())
	height := int32(tex.U32())
	tex.U32()
	tex.U32()
	rd := tex.Read(int(width * 4 * height))
	img, _ := sdl.CreateRGBSurfaceWithFormatFrom(unsafe.Pointer(&rd[0]), width, height, 32, int32(width*4), sdl.PIXELFORMAT_ABGR8888)
	return texData{img, rd}, false
}
