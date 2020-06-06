package temap

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

// https://github.com/mrkite/TerraFirma/blob/master/worldinfo.cpp#L125
func parseMB(tag string, blend bool, i int) (int, MergeBlend) {
	group := ""
	mb := MergeBlend{}
	mb.HasTile = false
	mb.Direction = 0
	mb.Mask = 0
	mb.Tile = 0
	mb.Blend = blend
	mb.Recursive = false
	lt := len(tag)
	for i < lt {
		c := tag[i]
		i += 1
		if c == ',' {
			break
		}
		if c == '*' {
			mb.Recursive = true
		} else if c == 'v' {
			mb.Direction |= 4
		} else if c == '^' {
			mb.Direction |= 8
		} else if c == '+' {
			mb.Direction |= 8 | 4 | 2 | 1
		} else if c >= '0' && c <= '9' {
			mb.HasTile = true
			mb.Tile *= 10
			mb.Tile += int16(c) - int16('0')
		} else if c >= 'a' && c <= 'z' {
			group += string(c)
		}
	}
	if mb.Direction == 0 {
		mb.Direction = 0xff
	}
	if !mb.HasTile {
		if group == "solid" {
			mb.Mask |= 1
		} else if group == "dirt" {
			mb.Mask |= 4
		} else if group == "brick" {
			mb.Mask |= 128
		} else if group == "moss" {
			mb.Mask |= 256
		}
	}
	return i, mb
}

func coerceToString(i interface{}) string {
	bs, ok := i.(string)
	if !ok {
		tmp, ok := i.(int)
		if ok {
			bs = strconv.Itoa(tmp)
		}
		if !ok {
			bs = ""
		}
	}
	return bs
}

// https://github.com/mrkite/TerraFirma/blob/master/worldinfo.cpp#L171
func (t *TileInfo) load(d tileData) {
	if d.Name.Valid {
		t.Name = d.Name.String
	}
	if d.Color.Valid {
		var c int64
		if d.Color.String[0] == '#' {
			c, _ = strconv.ParseInt(d.Color.String[1:], 16, 64)
		} else {
			c, _ = strconv.ParseInt(d.Color.String, 16, 64)
		}
		t.Color = uint32(c)
	}
	if d.Flags.Valid {
		t.Mask = uint32(d.Flags.Int64)
	}
	t.Mask = t.Mask
	t.Solid = t.Mask&1 != 0
	t.Transparent = t.Mask&2 != 0
	t.Dirt = t.Mask&4 != 0
	t.Stone = t.Mask&8 != 0
	t.Grass = t.Mask&16 != 0
	t.Pile = t.Mask&32 != 0
	t.Flip = t.Mask&64 != 0
	t.Brick = t.Mask&128 != 0
	t.Merge = t.Mask&512 != 0
	if t.Mask&1024 != 0 {
		t.Large = 1
	}
	if d.W.Valid {
		t.Width = int(d.W.Int64)
	} else {
		t.Width = 18
	}
	if d.H.Valid {
		t.Height = int(d.H.Int64)
	} else {
		t.Height = 18
	}
	if d.Skipy.Valid {
		t.Skipy = int(d.Skipy.Int64)
	}
	if d.Toppad.Valid {
		t.Toppad = int(d.Toppad.Int64)
	}
	t.Blends = make([]MergeBlend, 0)
	offset := 0
	bs := coerceToString(d.Blends)
	var mb MergeBlend
	for offset < len(bs) {
		offset, mb = parseMB(bs, true, offset)
		t.Blends = append(t.Blends, mb)
	}
	bs = coerceToString(d.Merge)
	offset = 0
	for offset < len(bs) {
		offset, mb = parseMB(bs, false, offset)
		t.Blends = append(t.Blends, mb)
	}

}

func (t *TileInfo) loadWall(d wallData) {
	if d.Large.Valid {
		t.Large = uint8(d.Large.Int64)
	}
	if d.Name.Valid {
		t.Name = d.Name.String
	}
	if d.Color.Valid {
		var c int64
		if d.Color.String[0] == '#' {
			c, _ = strconv.ParseInt(d.Color.String[1:], 16, 64)
		} else {
			c, _ = strconv.ParseInt(d.Color.String, 16, 64)
		}
		t.Color = uint32(c)
	}
}

func (t *TilesInfo) Load() {
	var tdata []tileData
	var wdata []wallData
	dat, _ := ioutil.ReadFile("./data/tiles.json")
	json.Unmarshal(dat, &tdata)
	dat, _ = ioutil.ReadFile("./data/walls.json")
	json.Unmarshal(dat, &wdata)
	t.WallsCount = int(wdata[len(wdata)-1].Id + 1)
	t.TilesCount = int(tdata[len(tdata)-1].Id + 1)
	t.Tiles = make([]TileInfo, 623)
	t.Walls = make([]TileInfo, 623)
	for _, d := range tdata {
		t.Tiles[d.Id].load(d)
	}
	for _, d := range wdata {
		t.Walls[d.Id].loadWall(d)
	}
}
