package tenet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/maciej-zuk/tecligo/temap"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type t2sMap struct {
	T int `json:"t"`
	S int `json:"s"`
}

var t2hs map[int]int
var t2bs map[int]int
var t2ls map[int]int

func init() {
	t2hs = make(map[int]int)
	t2ls = make(map[int]int)
	t2bs = make(map[int]int)
	var slotmap []t2sMap
	dat, _ := ioutil.ReadFile("./data/headslots.json")
	json.Unmarshal(dat, &slotmap)
	for _, e := range slotmap {
		t2hs[e.T] = e.S
	}
	dat, _ = ioutil.ReadFile("./data/bodyslot.json")
	json.Unmarshal(dat, &slotmap)
	for _, e := range slotmap {
		t2bs[e.T] = e.S
	}
	dat, _ = ioutil.ReadFile("./data/legslot.json")
	json.Unmarshal(dat, &slotmap)
	for _, e := range slotmap {
		t2ls[e.T] = e.S
	}
}

func savePlayerHeadIcon(slot TnByte, player *Player) {
	helmetItemId := 0
	if player.Inventory[69].Item != 0 {
		helmetItemId = int(player.Inventory[69].Item)
	} else if player.Inventory[59].Item != 0 {
		helmetItemId = int(player.Inventory[59].Item)
	}

	bodyItemId := 0
	if player.Inventory[70].Item != 0 {
		bodyItemId = int(player.Inventory[70].Item)
	} else if player.Inventory[60].Item != 0 {
		bodyItemId = int(player.Inventory[60].Item)
	}

	legItemId := 0
	if player.Inventory[71].Item != 0 {
		legItemId = int(player.Inventory[71].Item)
	} else if player.Inventory[61].Item != 0 {
		legItemId = int(player.Inventory[61].Item)
	}

	headSlot, hasHelmet := t2hs[helmetItemId]
	bodySlot, hasBody := t2bs[bodyItemId]
	legsSlot, hasLegs := t2ls[legItemId]

	head := temap.GetTextureCrop(temap.Tex_Player0|0, 40, 56)
	eye1 := temap.GetTextureCrop(temap.Tex_Player0|1, 40, 56)
	eye2 := temap.GetTextureCrop(temap.Tex_Player0|2, 40, 56)

	out, _ := sdl.CreateRGBSurfaceWithFormat(0, 40, 56, 32, sdl.PIXELFORMAT_ABGR8888)
	head.SetColorMod(player.SkinColor.R, player.SkinColor.G, player.SkinColor.B)
	head.Blit(&head.ClipRect, out, &out.ClipRect)
	eye1.Blit(&eye1.ClipRect, out, &out.ClipRect)
	eye2.SetColorMod(player.EyeColor.R, player.EyeColor.G, player.EyeColor.B)
	eye2.Blit(&eye2.ClipRect, out, &out.ClipRect)

	for i := 3; i <= 13; i++ {
		tex := temap.GetTextureCrop(temap.Tex_Player0|i, 40, 56)
		tex.Blit(&tex.ClipRect, out, &out.ClipRect)
	}

	if hasHelmet {
		helmet := temap.GetTextureCrop(temap.Tex_ArmorHead|headSlot, 40, 56)
		helmet.Blit(&helmet.ClipRect, out, &out.ClipRect)
	} else {
		hair := temap.GetTextureCrop(temap.Tex_PlayerHair|int(player.Hair), 40, 56)
		hair.SetColorMod(player.HairColor.R, player.HairColor.G, player.HairColor.B)
		hair.Blit(&hair.ClipRect, out, &out.ClipRect)
	}
	if hasLegs {
		legs := temap.GetTextureCrop(temap.Tex_ArmorLegs|legsSlot, 40, 56)
		legs.Blit(&legs.ClipRect, out, &out.ClipRect)
	}
	if hasBody {
		body := temap.GetTextureCrop(temap.Tex_ArmorBody|bodySlot, 40, 56)
		body.Blit(&body.ClipRect, out, &out.ClipRect)
	}

	img.SavePNG(out, fmt.Sprintf("map/player_%d.png", slot))
	out.Free()
}
