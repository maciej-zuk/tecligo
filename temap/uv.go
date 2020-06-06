package temap

import (
	"math"
	"math/rand"
)

type uvRule struct {
	Mask, Val uint16
	Uvs       [6]int16
	Blend     uint8
}

var grassRules = []uvRule{
	uvRule{0xaaaa, 0xaaaa, [6]int16{18, 18, 36, 18, 54, 18}, 0},
	uvRule{0x0abf, 0x083f, [6]int16{0, 324, 18, 324, 36, 324}, 0},
	uvRule{0x0abf, 0x023f, [6]int16{54, 324, 72, 324, 90, 324}, 0},
	uvRule{0xa0ef, 0x80cf, [6]int16{0, 342, 18, 342, 36, 342}, 0},
	uvRule{0xa0ef, 0x20cf, [6]int16{54, 342, 72, 342, 90, 342}, 0},
	uvRule{0x22fb, 0x20f3, [6]int16{54, 360, 72, 360, 90, 360}, 0},
	uvRule{0x22fb, 0x02f3, [6]int16{0, 360, 18, 360, 36, 360}, 0},
	uvRule{0x88fe, 0x80fc, [6]int16{0, 378, 18, 378, 36, 378}, 0},
	uvRule{0x88fe, 0x08fc, [6]int16{54, 378, 72, 378, 90, 378}, 0},
	uvRule{0x02bb, 0x0033, [6]int16{90, 270, 108, 270, 126, 270}, 0},
	uvRule{0x08be, 0x003c, [6]int16{144, 270, 162, 270, 180, 270}, 0},
	uvRule{0x20eb, 0x00c3, [6]int16{90, 288, 108, 288, 126, 288}, 0},
	uvRule{0x80ee, 0x00cc, [6]int16{144, 288, 162, 288, 180, 288}, 0},
	uvRule{0x0abf, 0x003f, [6]int16{144, 216, 198, 216, 252, 216}, 0},
	uvRule{0xa0ef, 0x00cf, [6]int16{144, 252, 198, 252, 252, 252}, 0},
	uvRule{0x22fb, 0x00f3, [6]int16{126, 234, 180, 234, 234, 234}, 0},
	uvRule{0x88fe, 0x00fc, [6]int16{162, 234, 216, 234, 270, 234}, 0},
	uvRule{0x00af, 0x002a, [6]int16{36, 270, 54, 270, 72, 270}, 0},
	uvRule{0x00af, 0x008a, [6]int16{36, 288, 54, 288, 72, 288}, 0},
	uvRule{0x00fa, 0x00a2, [6]int16{0, 270, 0, 288, 0, 306}, 0},
	uvRule{0x00fa, 0x00a8, [6]int16{18, 270, 18, 288, 18, 306}, 0},
	uvRule{0x00ff, 0x00ea, [6]int16{198, 288, 216, 288, 234, 288}, 0},
	uvRule{0x00ff, 0x00ba, [6]int16{198, 270, 216, 270, 234, 270}, 0},
	uvRule{0x00ff, 0x00ae, [6]int16{198, 306, 216, 306, 234, 306}, 0},
	uvRule{0x00ff, 0x00ab, [6]int16{144, 306, 162, 306, 180, 306}, 0},
	uvRule{0xaaaa, 0x2aaa, [6]int16{54, 108, 54, 144, 54, 180}, 0},
	uvRule{0xaaaa, 0x8aaa, [6]int16{36, 108, 36, 144, 36, 180}, 0},
	uvRule{0xaaaa, 0xa2aa, [6]int16{54, 90, 54, 126, 54, 162}, 0},
	uvRule{0xaaaa, 0xa8aa, [6]int16{36, 90, 36, 126, 36, 162}, 0},
	uvRule{0x00af, 0x002b, [6]int16{0, 198, 18, 198, 36, 198}, 0},
	uvRule{0x00af, 0x002e, [6]int16{54, 198, 72, 198, 90, 198}, 0},
	uvRule{0x00af, 0x008b, [6]int16{0, 216, 18, 216, 36, 216}, 0},
	uvRule{0x00af, 0x008e, [6]int16{54, 216, 72, 216, 90, 216}, 0},
	uvRule{0x00fa, 0x00b2, [6]int16{72, 144, 72, 162, 72, 180}, 0},
	uvRule{0x00fa, 0x00b8, [6]int16{90, 144, 90, 162, 90, 180}, 0},
	uvRule{0x57ff, 0x02ff, [6]int16{108, 324, 126, 324, 144, 324}, 0},
	uvRule{0x77ff, 0x20ff, [6]int16{108, 342, 126, 342, 144, 342}, 0},
	uvRule{0x7fff, 0x08ff, [6]int16{108, 360, 126, 360, 144, 360}, 0},
	uvRule{0xffff, 0x80ff, [6]int16{108, 378, 126, 378, 144, 378}, 0},
	uvRule{0xffff, 0x00ff, [6]int16{144, 234, 198, 234, 252, 234}, 0},
	uvRule{0x41ff, 0x00ff, [6]int16{36, 306, 54, 306, 72, 306}, 0},
	uvRule{0x14ff, 0x00ff, [6]int16{90, 306, 108, 306, 126, 306}, 0},
	uvRule{0x5fff, 0x0fff, [6]int16{54, 108, 54, 144, 54, 180}, 0},
	uvRule{0xdfff, 0xcfff, [6]int16{36, 108, 36, 144, 36, 180}, 0},
	uvRule{0xf7ff, 0xf3ff, [6]int16{54, 90, 54, 126, 54, 162}, 0},
	uvRule{0xfdff, 0xfcff, [6]int16{36, 90, 36, 126, 36, 162}, 0},
	uvRule{0xa0ff, 0x00ef, [6]int16{108, 18, 126, 18, 144, 18}, 0},
	uvRule{0x0aff, 0x00bf, [6]int16{108, 36, 126, 36, 144, 36}, 0},
	uvRule{0x22ff, 0x00fb, [6]int16{198, 0, 198, 18, 198, 36}, 0},
	uvRule{0x88ff, 0x00fe, [6]int16{180, 0, 180, 18, 180, 36}, 0},
	uvRule{0x20ff, 0x20ef, [6]int16{54, 108, 54, 144, 54, 180}, 0},
	uvRule{0x80ff, 0x80ef, [6]int16{36, 108, 36, 144, 36, 180}, 0},
	uvRule{0x02ff, 0x02bf, [6]int16{54, 90, 54, 126, 54, 162}, 0},
	uvRule{0x08ff, 0x08bf, [6]int16{36, 90, 36, 126, 36, 162}, 0},
	uvRule{0x80ff, 0x80fe, [6]int16{54, 90, 54, 126, 54, 162}, 0},
	uvRule{0x08ff, 0x08fe, [6]int16{54, 108, 54, 144, 54, 180}, 0},
	uvRule{0x20ff, 0x20fb, [6]int16{36, 90, 36, 126, 36, 162}, 0},
	uvRule{0x02ff, 0x02fb, [6]int16{36, 108, 36, 144, 36, 180}, 0},
	uvRule{0x00ff, 0x00bf, [6]int16{18, 18, 36, 18, 54, 18}, 0},
	uvRule{0x00ff, 0x00ef, [6]int16{18, 18, 36, 18, 54, 18}, 0},
	uvRule{0x00ff, 0x00fb, [6]int16{18, 18, 36, 18, 54, 18}, 0},
	uvRule{0x00ff, 0x00fe, [6]int16{18, 18, 36, 18, 54, 18}, 0},
}

var blendRules = []uvRule{
	uvRule{0x00ff, 0x00bf, [6]int16{144, 108, 162, 108, 180, 108}, 8},
	uvRule{0x00ff, 0x00ef, [6]int16{144, 90, 162, 90, 180, 90}, 4},
	uvRule{0x00ff, 0x00fb, [6]int16{162, 126, 162, 144, 162, 162}, 2},
	uvRule{0x00ff, 0x00fe, [6]int16{144, 126, 144, 144, 144, 162}, 1},
	uvRule{0x00ff, 0x00bb, [6]int16{36, 90, 36, 126, 36, 162}, 8 | 2},
	uvRule{0x00ff, 0x00be, [6]int16{54, 90, 54, 126, 54, 162}, 8 | 1},
	uvRule{0x00ff, 0x00eb, [6]int16{36, 108, 36, 144, 36, 180}, 4 | 2},
	uvRule{0x00ff, 0x00ee, [6]int16{54, 108, 54, 144, 54, 180}, 4 | 1},
	uvRule{0x00ff, 0x00fa, [6]int16{180, 126, 180, 144, 180, 162}, 2 | 1},
	uvRule{0x00ff, 0x00af, [6]int16{144, 180, 162, 180, 180, 180}, 8 | 4},
	uvRule{0x00ff, 0x00ba, [6]int16{198, 90, 198, 108, 198, 126}, 8 | 2 | 1},
	uvRule{0x00ff, 0x00ea, [6]int16{216, 144, 216, 162, 216, 180}, 4 | 2 | 1},
	uvRule{0x00ff, 0x00ab, [6]int16{216, 90, 216, 108, 216, 126}, 8 | 4 | 2},
	uvRule{0x00ff, 0x00aa, [6]int16{108, 198, 126, 198, 144, 198}, 8 | 4 | 2 | 1},
	uvRule{0x03ff, 0x02ff, [6]int16{0, 90, 0, 126, 0, 162}, 0},
	uvRule{0x0cff, 0x08ff, [6]int16{18, 90, 18, 126, 18, 162}, 0},
	uvRule{0x30ff, 0x20ff, [6]int16{0, 108, 0, 144, 0, 180}, 0},
	uvRule{0xc0ff, 0x80ff, [6]int16{18, 108, 18, 144, 18, 180}, 0},
}

var noGrassRules = []uvRule{
	uvRule{0x00fb, 0x00b3, [6]int16{72, 144, 72, 162, 72, 180}, 8},
	uvRule{0x00fb, 0x00e3, [6]int16{72, 90, 72, 108, 72, 126}, 4},
	uvRule{0x00fe, 0x00bc, [6]int16{90, 144, 90, 162, 90, 180}, 8},
	uvRule{0x00fe, 0x00ec, [6]int16{90, 90, 90, 108, 90, 126}, 4},
	uvRule{0x00bf, 0x003b, [6]int16{0, 198, 18, 198, 36, 198}, 2},
	uvRule{0x00bf, 0x003e, [6]int16{54, 198, 72, 198, 90, 198}, 1},
	uvRule{0x00ef, 0x00cb, [6]int16{0, 216, 18, 216, 36, 216}, 2},
	uvRule{0x00ef, 0x00ce, [6]int16{54, 216, 72, 216, 90, 216}, 1},
	uvRule{0x00fa, 0x00a0, [6]int16{108, 216, 108, 234, 108, 252}, 8 | 4},
	uvRule{0x00ca, 0x0080, [6]int16{126, 144, 126, 162, 126, 180}, 8},
	uvRule{0x003a, 0x0020, [6]int16{126, 90, 126, 108, 126, 126}, 4},
	uvRule{0x00af, 0x000a, [6]int16{162, 198, 180, 198, 198, 198}, 2 | 1},
	uvRule{0x00ac, 0x0008, [6]int16{0, 252, 18, 252, 36, 252}, 2},
	uvRule{0x00a3, 0x0002, [6]int16{54, 252, 72, 252, 90, 252}, 1},
	uvRule{0x00ea, 0x0080, [6]int16{108, 144, 108, 162, 108, 180}, 8},
	uvRule{0x00ba, 0x0020, [6]int16{108, 90, 108, 108, 108, 126}, 4},
	uvRule{0x00ae, 0x0008, [6]int16{0, 234, 18, 234, 36, 234}, 2},
	uvRule{0x00ab, 0x0002, [6]int16{54, 234, 72, 234, 90, 234}, 1},
	uvRule{0x00bf, 0x002f, [6]int16{234, 0, 252, 0, 270, 0}, 4},
	uvRule{0x00ef, 0x008f, [6]int16{234, 18, 252, 18, 270, 18}, 8},
	uvRule{0x00fb, 0x00f2, [6]int16{234, 36, 252, 36, 270, 36}, 1},
	uvRule{0x00fe, 0x00f8, [6]int16{234, 54, 252, 54, 270, 54}, 2},
}

var uvRules = []uvRule{
	uvRule{0x50ff, 0x00ff, [6]int16{108, 18, 126, 18, 144, 18}, 0},
	uvRule{0x05ff, 0x00ff, [6]int16{108, 36, 126, 36, 144, 36}, 0},
	uvRule{0x44ff, 0x00ff, [6]int16{180, 0, 180, 18, 180, 36}, 0},
	uvRule{0x11ff, 0x00ff, [6]int16{198, 0, 198, 18, 198, 36}, 0},
	uvRule{0x00ff, 0x00ff, [6]int16{18, 18, 36, 18, 54, 18}, 0},
	uvRule{0x007f, 0x003f, [6]int16{18, 0, 36, 0, 54, 0}, 0},
	uvRule{0x00df, 0x00cf, [6]int16{18, 36, 36, 36, 54, 36}, 0},
	uvRule{0x00f7, 0x00f3, [6]int16{0, 0, 0, 18, 0, 36}, 0},
	uvRule{0x00fd, 0x00fc, [6]int16{72, 0, 72, 18, 72, 36}, 0},
	uvRule{0x0077, 0x0033, [6]int16{0, 54, 36, 54, 72, 54}, 0},
	uvRule{0x007d, 0x003c, [6]int16{18, 54, 54, 54, 90, 54}, 0},
	uvRule{0x00d7, 0x00c3, [6]int16{0, 72, 36, 72, 72, 72}, 0},
	uvRule{0x00dd, 0x00cc, [6]int16{18, 72, 54, 72, 90, 72}, 0},
	uvRule{0x00f5, 0x00f0, [6]int16{90, 0, 90, 18, 90, 36}, 0},
	uvRule{0x005f, 0x000f, [6]int16{108, 72, 126, 72, 144, 72}, 0},
	uvRule{0x0075, 0x0030, [6]int16{108, 0, 126, 0, 144, 0}, 0},
	uvRule{0x00d5, 0x00c0, [6]int16{108, 54, 126, 54, 144, 54}, 0},
	uvRule{0x0057, 0x0003, [6]int16{162, 0, 162, 18, 162, 36}, 0},
	uvRule{0x005d, 0x000c, [6]int16{216, 0, 216, 18, 216, 36}, 0},
	uvRule{0x0055, 0x0000, [6]int16{162, 54, 180, 54, 198, 54}, 0},
	uvRule{0x0000, 0x0000, [6]int16{18, 18, 36, 18, 54, 18}, 0},
}

var cactusRules = []uvRule{
	uvRule{0x37b, 0x003, [6]int16{90, 0, 0, 0, 0, 0}, 0},
	uvRule{0x36a, 0x002, [6]int16{72, 0, 0, 0, 0, 0}, 0},
	uvRule{0x319, 0x001, [6]int16{18, 0, 0, 0, 0, 0}, 0},
	uvRule{0x308, 0x000, [6]int16{0, 0, 0, 0, 0, 0}, 0},
	uvRule{0x37b, 0x00b, [6]int16{90, 36, 0, 0, 0, 0}, 0},
	uvRule{0x36a, 0x00a, [6]int16{72, 36, 0, 0, 0, 0}, 0},
	uvRule{0x319, 0x009, [6]int16{18, 36, 0, 0, 0, 0}, 0},
	uvRule{0x380, 0x080, [6]int16{0, 36, 0, 0, 0, 0}, 0},
	uvRule{0x300, 0x000, [6]int16{0, 18, 0, 0, 0, 0}, 0},
	uvRule{0x30d, 0x101, [6]int16{108, 36, 0, 0, 0, 0}, 0},
	uvRule{0x305, 0x101, [6]int16{54, 36, 0, 0, 0, 0}, 0},
	uvRule{0x309, 0x101, [6]int16{54, 0, 0, 0, 0, 0}, 0},
	uvRule{0x301, 0x101, [6]int16{54, 18, 0, 0, 0, 0}, 0},
	uvRule{0x309, 0x100, [6]int16{54, 0, 0, 0, 0, 0}, 0},
	uvRule{0x300, 0x100, [6]int16{54, 18, 0, 0, 0, 0}, 0},
	uvRule{0x30e, 0x202, [6]int16{108, 18, 0, 0, 0, 0}, 0},
	uvRule{0x306, 0x202, [6]int16{36, 36, 0, 0, 0, 0}, 0},
	uvRule{0x30a, 0x202, [6]int16{36, 0, 0, 0, 0, 0}, 0},
	uvRule{0x302, 0x202, [6]int16{36, 18, 0, 0, 0, 0}, 0},
	uvRule{0x30a, 0x200, [6]int16{36, 0, 0, 0, 0, 0}, 0},
	uvRule{0x300, 0x200, [6]int16{36, 18, 0, 0, 0, 0}, 0},
}

var phlebasTiles = [4][3]uint16{
	[3]uint16{2, 4, 2},
	[3]uint16{1, 3, 1},
	[3]uint16{2, 2, 4},
	[3]uint16{1, 1, 3},
}

var lazureTiles = [2][2]uint16{
	[2]uint16{1, 2},
	[2]uint16{3, 4},
}
var wallRandom = [3][3]uint16{
	[3]uint16{2, 0, 0},
	[3]uint16{0, 1, 4},
	[3]uint16{0, 3, 0},
}

var walluvs = [][8]int16{
	[8]int16{324, 108, 360, 108, 396, 108, 216, 216},
	[8]int16{216, 108, 252, 108, 288, 108, 144, 216},
	[8]int16{432, 0, 432, 36, 432, 72, 432, 180},
	[8]int16{36, 144, 108, 144, 180, 144, 108, 216},
	[8]int16{324, 0, 324, 36, 324, 72, 324, 180},
	[8]int16{0, 144, 72, 144, 144, 144, 72, 216},
	[8]int16{216, 144, 252, 144, 288, 144, 180, 216},
	[8]int16{36, 72, 72, 72, 108, 72, 108, 180},
	[8]int16{216, 0, 252, 0, 288, 0, 216, 180},
	[8]int16{180, 0, 180, 36, 180, 72, 180, 180},
	[8]int16{36, 108, 108, 108, 180, 108, 36, 216},
	[8]int16{144, 0, 144, 36, 144, 72, 144, 180},
	[8]int16{0, 108, 72, 108, 144, 108, 0, 216},
	[8]int16{0, 0, 0, 36, 0, 72, 0, 180},
	[8]int16{36, 0, 72, 0, 108, 0, 36, 216},
	[8]int16{36, 36, 72, 36, 108, 36, 72, 180},
	[8]int16{216, 36, 252, 36, 288, 36, 252, 180},
	[8]int16{216, 72, 252, 72, 288, 72, 288, 180},
	[8]int16{360, 0, 360, 36, 360, 72, 360, 180},
	[8]int16{396, 0, 396, 36, 396, 72, 396, 180},
}

// https://github.com/mrkite/TerraFirma/blob/master/uvrules.cpp#L217
func fixWall(world *World, x uint32, y uint32) {
	th := world.tilesHigh
	tw := world.tilesWide
	tile := &world.tiles[world.pt(x, y)]
	var mask uint16 = 0
	if y > 0 {
		top := world.tiles[world.pt(x, y-1)]
		if top.Wall != 0 || (top.Active && top.Type == 54) {
			mask |= 1
		}
	}
	if x > 0 {
		left := world.tiles[world.pt(x-1, y)]
		if left.Wall != 0 || (left.Active && left.Type == 54) {
			mask |= 2
		}
	}
	if x < tw-1 {
		right := world.tiles[world.pt(x+1, y)]
		if right.Wall != 0 || (right.Active && right.Type == 54) {
			mask |= 4
		}
	}
	if y < th-1 {
		bottom := world.tiles[world.pt(x, y+1)]
		if bottom.Wall != 0 || (bottom.Active && bottom.Type == 54) {
			mask |= 8
		}
	}

	lrg := world.info.Walls[tile.Wall].Large
	var st uint16
	if lrg == 1 {
		st = (phlebasTiles[y%4][x%3] - 1) * 2
	} else if lrg == 2 {
		st = (lazureTiles[x%2][y%2] - 1) * 2
	} else {
		st = uint16(math.Floor(rand.Float64()*3) * 2)
	}

	if mask == 15 {
		mask += wallRandom[x%3][y%3]
	}

	tile.Wallu = walluvs[mask][st]
	tile.Wallv = walluvs[mask][st+1]
}

// https://github.com/mrkite/TerraFirma/blob/master/uvrules.cpp#L550
func fixCactus(world *World, x uint32, y uint32) {
	basex := x
	basey := y

	for world.tiles[world.pt(basex, basey)].Active && world.tiles[world.pt(basex, basey)].Type == 80 {
		basey += 1
		if !world.tiles[world.pt(basex, basey)].Active || world.tiles[world.pt(basex, basey)].Type != 80 {
			if world.tiles[world.pt(basex-1, basey)].Active &&
				world.tiles[world.pt(basex-1, basey)].Type == 80 &&
				world.tiles[world.pt(basex-1, basey-1)].Active &&
				world.tiles[world.pt(basex-1, basey-1)].Type == 80 &&
				basex >= x {
				basex -= 1
			}
			if world.tiles[world.pt(basex+1, basey)].Active &&
				world.tiles[world.pt(basex+1, basey)].Type == 80 &&
				world.tiles[world.pt(basex+1, basey-1)].Active &&
				world.tiles[world.pt(basex+1, basey-1)].Type == 80 &&
				basex <= x {
				basex += 1
			}
		}
	}

	th := world.tilesHigh
	tw := world.tilesWide

	var mask uint16 = 0
	if x < tw-1 {
		right := world.tiles[world.pt(x+1, y)]
		if right.Active && right.Type == 80 {
			mask |= 0x01
		}
	}
	if x > 0 {
		left := world.tiles[world.pt(x-1, y)]
		if left.Active && left.Type == 80 {
			mask |= 0x02
		}
		if x > 1 {
			fl := world.tiles[world.pt(x-2, y)]
			if fl.Active && fl.Type == 80 {
				mask |= 0x40
			}
		}
	}

	if y < th {
		bottom := world.tiles[world.pt(x, y+1)]
		if bottom.Active && bottom.Type == 80 {
			mask |= 0x04
		}
		if bottom.Active && world.info.Tiles[bottom.Type].Solid {
			mask |= 0x80
		}
		if x < tw-1 {
			br := world.tiles[world.pt(x+1, y+1)]
			if br.Active && br.Type == 80 {
				mask |= 0x10
			}
		}
		if x > 0 {
			bl := world.tiles[world.pt(x-1, y+1)]
			if bl.Active && bl.Type == 80 {
				mask |= 0x20
			}
		}
	}

	if y > 0 {
		top := world.tiles[world.pt(x, y-1)]
		if top.Active && (top.Type == 80 || top.Type == 227) {
			mask |= 0x08
		}
	}

	if x > basex {
		mask |= 0x200
	}
	if x < basex {
		mask |= 0x100
	}

	for _, rule := range cactusRules {
		if (mask & rule.Mask) == rule.Val {
			tile := &world.tiles[world.pt(x, y)]
			tile.U = rule.Uvs[0]
			tile.V = rule.Uvs[1]
			return
		}
	}
	return
}

// https://github.com/mrkite/TerraFirma/blob/master/uvrules.cpp#L264
func fixTile(world *World, x uint32, y uint32) uint8 {
	var t int = -1
	var l int = -1
	var r int = -1
	var b int = -1
	var tl int = -1
	var tr int = -1
	var bl int = -1
	var br int = -1
	stride := world.tilesWide
	offset := y*stride + x

	tile := world.tiles[world.pt(x, y)]
	c := tile.Type
	if world.info.Tiles[c].Stone {
		c = 1
	}
	if c == 80 {
		fixCactus(world, x, y)
		return 0
	}

	if x > 0 {
		left := world.tiles[offset-1]
		if left.Active && left.Slope != 1 && left.Slope != 3 {
			l = left.Type
			if world.info.Tiles[l].Stone {
				l = 1
			}
		}
		if y > 0 && world.tiles[offset-stride-1].Active {
			tl = world.tiles[offset-stride-1].Type
			if world.info.Tiles[tl].Stone {
				tl = 1
			}
		}
		if y < world.tilesHigh-1 && world.tiles[offset+stride-1].Active {
			bl = world.tiles[offset+stride-1].Type
			if world.info.Tiles[bl].Stone {
				bl = 1
			}
		}
	}
	if x < world.tilesWide-1 {
		right := world.tiles[offset+1]
		if right.Active && right.Slope != 2 && right.Slope != 4 {
			r = right.Type
			if world.info.Tiles[r].Stone {
				r = 1
			}
		}
		if y > 0 && world.tiles[offset-stride+1].Active {
			tr = world.tiles[offset-stride+1].Type
			if world.info.Tiles[tr].Stone {
				tr = 1
			}
		}
		if y < world.tilesHigh-1 && world.tiles[offset+stride+1].Active {
			br = world.tiles[offset+stride+1].Type
			if world.info.Tiles[br].Stone {
				br = 1
			}
		}
	}
	if y > 0 {
		top := world.tiles[offset-stride]
		if top.Active && top.Slope != 3 && top.Slope != 4 {
			t = top.Type
			if world.info.Tiles[t].Stone {
				t = 1
			}
		}
	}
	if y < world.tilesHigh-1 {
		bottom := world.tiles[offset+stride]
		if bottom.Active && bottom.Slope != 1 && bottom.Slope != 2 {
			b = bottom.Type
			if world.info.Tiles[b].Stone {
				b = 1
			}
		}
	}

	// fix slopes
	switch tile.Slope {
	case 1:
		t = -1
		r = -1
		break
	case 2:
		t = -1
		l = -1
		break
	case 3:
		b = -1
		r = -1
		break
	case 4:
		b = -1
		l = -1
		break
	}

	// check blends and merges (blends should be first)
	for _, blend := range world.info.Tiles[c].Blends {
		var dir uint8 = 0
		if blend.HasTile {
			if t == int(blend.Tile) {
				dir |= 8
			}
			if b == int(blend.Tile) {
				dir |= 4
			}
			if l == int(blend.Tile) {
				dir |= 2
			}
			if r == int(blend.Tile) {
				dir |= 1
			}
			if tl == int(blend.Tile) {
				dir |= 0x80
			}
			if tr == int(blend.Tile) {
				dir |= 0x40
			}
			if bl == int(blend.Tile) {
				dir |= 0x20
			}
			if br == int(blend.Tile) {
				dir |= 0x10
			}
		} else {
			if t > -1 && (world.info.Tiles[t].Mask&blend.Mask) != 0 {
				dir |= 8
			}
			if b > -1 && (world.info.Tiles[b].Mask&blend.Mask) != 0 {
				dir |= 4
			}
			if l > -1 && (world.info.Tiles[l].Mask&blend.Mask) != 0 {
				dir |= 2
			}
			if r > -1 && (world.info.Tiles[r].Mask&blend.Mask) != 0 {
				dir |= 1
			}
			if tl > -1 && (world.info.Tiles[tl].Mask&blend.Mask) != 0 {
				dir |= 0x80
			}
			if tr > -1 && (world.info.Tiles[tr].Mask&blend.Mask) != 0 {
				dir |= 0x40
			}
			if bl > -1 && (world.info.Tiles[bl].Mask&blend.Mask) != 0 {
				dir |= 0x20
			}
			if br > -1 && (world.info.Tiles[br].Mask&blend.Mask) != 0 {
				dir |= 0x10
			}
		}
		dir &= blend.Direction

		if (dir&8) != 0 && (!blend.Recursive || (fixTile(world, x, y-1)&4) != 0) {
			if blend.Blend {
				t = -2
			} else {
				t = c
			}
		}
		if (dir&4) != 0 && (!blend.Recursive || (fixTile(world, x, y+1)&8) != 0) {
			if blend.Blend {
				b = -2
			} else {
				b = c
			}
		}
		if (dir&2) != 0 && (!blend.Recursive || (fixTile(world, x-1, y)&1) != 0) {
			if blend.Blend {
				l = -2
			} else {
				l = c
			}
		}
		if (dir&1) != 0 && (!blend.Recursive || (fixTile(world, x+1, y)&2) != 0) {
			if blend.Blend {
				r = -2
			} else {
				r = c
			}
		}
		if dir&0x80 != 0 {
			if blend.Blend {
				tl = -2
			} else {
				tl = c
			}
		}
		if dir&0x40 != 0 {
			if blend.Blend {
				tr = -2
			} else {
				tr = c
			}
		}
		if dir&0x20 != 0 {
			if blend.Blend {
				bl = -2
			} else {
				bl = c
			}
		}
		if dir&0x10 != 0 {
			if blend.Blend {
				br = -2
			} else {
				br = c
			}
		}
	}
	if world.info.Tiles[c].Brick { // brick merges with brick
		if t > -1 && world.info.Tiles[t].Brick {
			t = c
		}
		if b > -1 && world.info.Tiles[b].Brick {
			b = c
		}
		if l > -1 && world.info.Tiles[l].Brick {
			l = c
		}
		if r > -1 && world.info.Tiles[r].Brick {
			r = c
		}
		if tl > -1 && world.info.Tiles[tl].Brick {
			tl = c
		}
		if tr > -1 && world.info.Tiles[tr].Brick {
			tr = c
		}
		if bl > -1 && world.info.Tiles[bl].Brick {
			bl = c
		}
		if br > -1 && world.info.Tiles[br].Brick {
			br = c
		}
	}
	if world.info.Tiles[c].Pile { // pile merges with pile
		if t > -1 && world.info.Tiles[t].Pile {
			t = c
		}
		if b > -1 && world.info.Tiles[b].Pile {
			b = c
		}
		if l > -1 && world.info.Tiles[l].Pile {
			l = c
		}
		if r > -1 && world.info.Tiles[r].Pile {
			r = c
		}
		if tl > -1 && world.info.Tiles[tl].Pile {
			tl = c
		}
		if tr > -1 && world.info.Tiles[tr].Pile {
			tr = c
		}
		if bl > -1 && world.info.Tiles[bl].Pile {
			bl = c
		}
		if br > -1 && world.info.Tiles[br].Pile {
			br = c
		}
	}
	if world.info.Tiles[c].Dirt {
		if t == 0 {
			t = -2
		}
		if b == 0 {
			b = -2
		}
		if l == 0 {
			l = -2
		}
		if r == 0 {
			r = -2
		}
		if tl == 0 {
			tl = -2
		}
		if tr == 0 {
			tr = -2
		}
		if bl == 0 {
			bl = -2
		}
		if br == 0 {
			br = -2
		}
	}
	// everything merges with 357
	if t == 357 {
		t = c
	}
	if b == 357 {
		b = c
	}
	if l == 357 {
		l = c
	}
	if r == 357 {
		r = c
	}
	if tl == 357 {
		tl = c
	}
	if tr == 357 {
		tr = c
	}
	if bl == 357 {
		bl = c
	}
	if br == 357 {
		br = c
	}
	// fix rope
	if c == 213 {
		if t != c {
			if l > -1 && world.info.Tiles[l].Solid {
				l = c
			}
			if r > -1 && world.info.Tiles[r].Solid {
				r = c
			}
		}
	}
	// fix cobweb
	if c == 51 {
		if t > -1 {
			t = c
		}
		if b > -1 {
			b = c
		}
		if l > -1 {
			l = c
		}
		if r > -1 {
			r = c
		}
		if tl > -1 {
			tl = c
		}
		if tr > -1 {
			tr = c
		}
		if bl > -1 {
			bl = c
		}
		if br > -1 {
			br = c
		}
	}

	// slope and half rules
	if (tile.Slope == 1 || tile.Slope == 2) && b > -1 && b != 19 {
		b = c
	}
	if t > -1 {
		top := world.tiles[offset-stride]
		if (top.Slope == 1 || top.Slope == 2) && t != 19 {
			t = c
		}
		if top.Half && t != 19 {
			t = c
		}
	}
	if (tile.Slope == 3 || tile.Slope == 4) && t > -1 && t != 19 {
		t = c
	}
	if b > -1 {
		bottom := world.tiles[offset+stride]
		if (bottom.Slope == 3 || bottom.Slope == 4) && b != 19 {
			b = c
		}
		if bottom.Half {
			b = -1
		}
	}
	if l > -1 {
		left := world.tiles[offset-1]
		if left.Half {
			if tile.Half {
				l = c
			} else if left.Type != c {
				l = -1
			}
		}
	}
	if r > -1 {
		right := world.tiles[offset+1]
		if right.Half {
			if tile.Half {
				r = c
			} else if right.Type != c {
				r = -1
			}
		}
	}
	if tile.Half {
		if l != c {
			l = -1
		}
		if r != c {
			r = -1
		}
		t = -1
	}

	var blend int = 0

	// fix color mismatches
	if !world.info.Tiles[c].Grass {
		if t == -2 && tile.Color != world.tiles[offset-stride].Color {
			blend |= 8
			t = c
		}
		if b == -2 && tile.Color != world.tiles[offset+stride].Color {
			blend |= 4
			b = c
		}
		if l == -2 && tile.Color != world.tiles[offset-1].Color {
			blend |= 2
			l = c
		}
		if r == -2 && tile.Color != world.tiles[offset+1].Color {
			blend |= 1
			r = c
		}
	}

	var mask uint16 = 0
	if t == c {
		mask |= 0xc0
	} else if t == -2 {
		mask |= 0x80
	}
	if b == c {
		mask |= 0x30
	} else if b == -2 {
		mask |= 0x20
	}
	if l == c {
		mask |= 0x0c
	} else if l == -2 {
		mask |= 0x08
	}
	if r == c {
		mask |= 0x03
	} else if r == -2 {
		mask |= 0x02
	}
	if tl == c {
		mask |= 0xc000
	} else if tl == -2 {
		mask |= 0x8000
	}
	if tr == c {
		mask |= 0x3000
	} else if tr == -2 {
		mask |= 0x2000
	}
	if bl == c {
		mask |= 0x0c00
	} else if bl == -2 {
		mask |= 0x0800
	}
	if br == c {
		mask |= 0x0300
	} else if br == -2 {
		mask |= 0x0200
	}

	st := uint16(math.Floor(rand.Float64()*3) * 2)
	if world.info.Tiles[c].Large != 0 {
		st = (phlebasTiles[y%4][x%3] - 1) * 2
	}

	st %= 6

	if world.info.Tiles[c].Grass {
		for _, rule := range grassRules {
			if (mask & rule.Mask) == rule.Val {
				tile := &world.tiles[offset]
				tile.U = rule.Uvs[st]
				tile.V = rule.Uvs[st+1]
				return uint8(int(rule.Blend) | blend)
			}
		}
	}

	if world.info.Tiles[c].Merge || world.info.Tiles[c].Dirt {
		for _, rule := range blendRules {
			if (mask & rule.Mask) == rule.Val {
				tile := &world.tiles[offset]
				tile.U = rule.Uvs[st]
				tile.V = rule.Uvs[st+1]
				if world.info.Tiles[c].Large != 0 && st == 6 {
					tile.V += 90
				}
				return uint8(int(rule.Blend) | blend)
			}
		}
		if !world.info.Tiles[c].Grass {
			for _, rule := range noGrassRules {
				if (mask & rule.Mask) == rule.Val {
					tile := &world.tiles[offset]
					tile.U = rule.Uvs[st]
					tile.V = rule.Uvs[st+1]
					if world.info.Tiles[c].Large != 0 && st == 6 {
						tile.V += 90
					}
					return uint8(int(rule.Blend) | blend)
				}
			}
		}
	}
	// no match, blends become merges
	if world.info.Tiles[c].Grass {
		mask |= (mask & 0xaaaa) >> 1
	}

	for _, rule := range uvRules {
		if (mask & rule.Mask) == rule.Val {
			tile := &world.tiles[offset]
			tile.U = rule.Uvs[st]
			tile.V = rule.Uvs[st+1]
			if world.info.Tiles[c].Large != 0 && st == 6 {
				tile.V += 90
			}
			return uint8(int(rule.Blend) | blend)
		}
	}
	// should never get here.. since there's a catch-all rule in uvRules
	return uint8(blend)
}
