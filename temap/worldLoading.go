package temap

import (
	"encoding/json"
	"io/ioutil"
)

// https://github.com/mrkite/TerraFirma/blob/master/worldheader.cpp#L37
func (w *World) loadHeader(handle *Handle) {
	handle.Seek(int64(w.sections[0]), 0)
	var worldHeader WorldHeader
	worldHeader.Name = handle.S()
	worldHeader.Seed = handle.S()
	worldHeader.GenVersion = handle.U64()
	handle.Skip(16)
	worldHeader.WorldID = handle.U32()
	worldHeader.LeftWorld = handle.U32()
	worldHeader.RightWorld = handle.U32()
	worldHeader.TopWorld = handle.U32()
	worldHeader.BottomWorld = handle.U32()
	worldHeader.MaxTilesY = handle.U32()
	worldHeader.MaxTilesX = handle.U32()
	worldHeader.GameMode = handle.U32()
	worldHeader.DrunkWorld = handle.Bool()
	worldHeader.CreationTime = handle.U64()
	worldHeader.MoonType = handle.U8()
	worldHeader.TreeX[0] = handle.U32()
	worldHeader.TreeX[1] = handle.U32()
	worldHeader.TreeX[2] = handle.U32()
	worldHeader.TreeStyle[0] = handle.U32()
	worldHeader.TreeStyle[1] = handle.U32()
	worldHeader.TreeStyle[2] = handle.U32()
	worldHeader.TreeStyle[3] = handle.U32()
	worldHeader.CaveBackX[0] = handle.U32()
	worldHeader.CaveBackX[1] = handle.U32()
	worldHeader.CaveBackX[2] = handle.U32()
	worldHeader.CaveBackStyle[0] = handle.U32()
	worldHeader.CaveBackStyle[1] = handle.U32()
	worldHeader.CaveBackStyle[2] = handle.U32()
	worldHeader.CaveBackStyle[3] = handle.U32()
	worldHeader.IceBackStyle = handle.U32()
	worldHeader.JungleBackStyle = handle.U32()
	worldHeader.HellBackStyle = handle.U32()
	worldHeader.SpawnTileX = handle.U32()
	worldHeader.SpawnTileY = handle.U32()
	worldHeader.WorldSurface = handle.D()
	worldHeader.RockLayer = handle.D()
	handle.D()
	handle.Bool()
	handle.U32()
	handle.Bool()
	handle.Bool()
	worldHeader.DungeonX = handle.U32()
	worldHeader.DungeonY = handle.U32()
	worldHeader.Crimson = handle.Bool()
	worldHeader.DownedBoss1 = handle.Bool()
	worldHeader.DownedBoss2 = handle.Bool()
	worldHeader.DownedBoss3 = handle.Bool()
	worldHeader.DownedQueenBee = handle.Bool()
	worldHeader.DownedMechBoss1 = handle.Bool()
	worldHeader.DownedMechBoss2 = handle.Bool()
	worldHeader.DownedMechBoss3 = handle.Bool()
	worldHeader.DownedMechBossAny = handle.Bool()
	worldHeader.DownedPlantBoss = handle.Bool()
	worldHeader.DownedGolemBoss = handle.Bool()
	worldHeader.DownedSlimeKing = handle.Bool()
	worldHeader.SavedGoblin = handle.Bool()
	worldHeader.SavedWizard = handle.Bool()
	worldHeader.SavedMech = handle.Bool()
	worldHeader.DownedGoblins = handle.Bool()
	worldHeader.DownedClown = handle.Bool()
	worldHeader.DownedFrost = handle.Bool()
	worldHeader.DownedPirates = handle.Bool()
	worldHeader.ShadowOrbSmashed = handle.Bool()
	worldHeader.SpawnMeteor = handle.Bool()
	worldHeader.ShadowOrbCount = handle.U8()
	worldHeader.AltarCount = handle.U32()
	worldHeader.HardMode = handle.Bool()
	worldHeader.InvasionDelay = handle.U32()
	worldHeader.InvasionSize = handle.U32()
	worldHeader.InvasionType = handle.U32()
	worldHeader.InvasionX = handle.D()
	worldHeader.SlimeRainTime = handle.D()
	worldHeader.SundialCooldown = handle.U8()
	handle.Bool()
	handle.U32()
	handle.F()
	worldHeader.SavedOreTiersCobalt = handle.U32()
	worldHeader.SavedOreTiersMythril = handle.U32()
	worldHeader.SavedOreTiersAdamantite = handle.U32()
	for num := 0; num < 8; num++ {
		worldHeader.Styles[num] = handle.U8()
	}
	worldHeader.CloudBGActive = handle.U32()
	worldHeader.NumClouds = handle.U16()
	worldHeader.WindSpeedTarget = handle.F()
	n0 := handle.U32()
	for num := n0; num > 0; num-- {
		handle.S()
	}
	worldHeader.SavedAngler = handle.Bool()
	worldHeader.AnglerQuest = handle.U32()
	worldHeader.SavedStylist = handle.Bool()
	worldHeader.SavedTaxCollector = handle.Bool()
	worldHeader.SavedGolfer = handle.Bool()
	worldHeader.InvasionSizeStart = handle.U32()
	handle.U32()
	for i := handle.U16(); i > 0; i-- {
		handle.U32()
	}
	worldHeader.FastForwardTime = handle.Bool()
	worldHeader.DownedFishron = handle.Bool()
	worldHeader.DownedMartians = handle.Bool()
	worldHeader.DownedAncientCultist = handle.Bool()
	worldHeader.DownedMoonlord = handle.Bool()
	worldHeader.DownedHalloweenKing = handle.Bool()
	worldHeader.DownedHalloweenTree = handle.Bool()
	worldHeader.DownedChristmasIceQueen = handle.Bool()
	worldHeader.DownedChristmasSantank = handle.Bool()
	worldHeader.DownedChristmasTree = handle.Bool()
	worldHeader.DownedTowerSolar = handle.Bool()
	worldHeader.DownedTowerVortex = handle.Bool()
	worldHeader.DownedTowerNebula = handle.Bool()
	worldHeader.DownedTowerStardust = handle.Bool()
	worldHeader.TowerActiveSolar = handle.Bool()
	worldHeader.TowerActiveVortex = handle.Bool()
	worldHeader.TowerActiveNebula = handle.Bool()
	worldHeader.TowerActiveStardust = handle.Bool()
	worldHeader.LunarApocalypseIsUp = handle.Bool()
	handle.Bool()
	handle.Bool()
	handle.U32()
	for j := handle.U32(); j > 0; j-- {
		handle.U32()
	}
	handle.Bool()
	handle.U32()
	handle.F()
	handle.F()
	worldHeader.SavedBartender = handle.Bool()
	handle.Bool()
	handle.Bool()
	handle.Bool()
	handle.U8()
	handle.U8()
	handle.U8()
	handle.U8()
	handle.U8()
	worldHeader.CombatBookWasUsed = handle.Bool()
	handle.U32()
	handle.Bool()
	handle.Bool()
	handle.Bool()
	for k := handle.U32(); k > 0; k-- {
		handle.U32()
	}
	worldHeader.ForceHalloweenForToday = handle.Bool()
	worldHeader.ForceXMasForToday = handle.Bool()
	worldHeader.SavedOreTiersCopper = handle.U32()
	worldHeader.SavedOreTiersIron = handle.U32()
	worldHeader.SavedOreTiersSilver = handle.U32()
	worldHeader.SavedOreTiersGold = handle.U32()
	worldHeader.BoughtCat = handle.Bool()
	worldHeader.BoughtDog = handle.Bool()
	worldHeader.BoughtBunny = handle.Bool()
	worldHeader.DownedEmpressOfLight = handle.Bool()
	worldHeader.DownedQueenSlime = handle.Bool()

	w.header = worldHeader
	w.groundLevel = uint32(w.header.WorldSurface)
	w.tilesWide = w.header.MaxTilesX
	w.tilesHigh = w.header.MaxTilesY
}

// https://github.com/mrkite/TerraFirma/blob/master/world.cpp#L565
func loadTile(handle *Handle, extra []bool) (rle uint32, tile Tile) {
	tile = Tile{}
	f1 := handle.U8()
	f2 := uint8(0)
	f3 := uint8(0)
	if f1&1 != 0 {
		f2 = handle.U8()
		if f2&1 != 0 {
			f3 = handle.U8()
		}
	}
	tile.Active = f1&2 != 0
	tile.U = -1
	tile.V = -1
	var t int
	if tile.Active {
		t = int(handle.U8())
		if f1&0x20 != 0 {
			t |= int(handle.U8()) << 8
		}
		if extra[t] {
			tile.U = int16(handle.U16())
			tile.V = int16(handle.U16())
		}
		if f3&0x8 != 0 {
			tile.Color = handle.U8()
		}
	} else {
		t = 0
	}
	tile.Type = t
	tile.Wall = 0
	if f1&4 != 0 {
		tile.Wall = handle.U8()
		if f3&0x10 != 0 {
			tile.WallColor = handle.U8()
		}
		tile.Wallu = -1
		tile.Wallv = -1
	}
	if f1&0x18 != 0 {
		tile.Liquid = handle.U8()
		tile.LiquidType = 0
		if (f1 & 0x18) == 0x10 {
			tile.LiquidType = 1
		} else if (f1 & 0x18) == 0x18 {
			tile.LiquidType = 2
		}
	} else {
		tile.Liquid = 0
	}
	slop := (f2 >> 4) & 7
	if slop == 1 {
		tile.Half = true
	}
	if slop > 1 {
		tile.Slope = slop - 1
	}
	if f3&4 != 0 {
		tile.Inactive = true
	}
	f6 := f1 >> 6
	rle = 0
	if f6 == 1 {
		rle = uint32(handle.U8())
	}
	if f6 == 2 {
		rle = uint32(handle.U16())
	}
	return rle, tile
}

// https://github.com/mrkite/TerraFirma/blob/master/world.cpp#L135
func (w *World) LoadTiles() {
	tw := w.tilesWide
	w.LoadTilesRegion(0, tw)
}

func (w *World) LoadTilesRegion(x0 uint32, x1 uint32) {
	th := w.tilesHigh
	handle := Handle{}
	handle.LoadFile(w.name)
	defer handle.Close()
	handle.Seek(int64(w.sections[1]), 0)
	w.tiles = make([]Tile, w.tilesWide*w.tilesHigh)
	for x := uint32(0); x < x1; x++ {
		for y := uint32(0); y < th; y++ {
			rle, tile := loadTile(&handle, w.extra)
			w.tiles[w.pt(x, y)] = tile
			if x > x0 {
				for ry := uint32(0); ry < rle; ry++ {
					w.tiles[w.pt(x, y+ry+1)] = tile
				}
			}
			y += rle
		}
	}
}

func (w *World) RecalcUV() {
	w.RecalcUVRegion(uint32(0), uint32(0), w.tilesWide, w.tilesHigh)
}

func (w *World) RecalcUVRegion(x0 uint32, y0 uint32, x1 uint32, y1 uint32) {
	for x := x0; x < x1; x++ {
		for y := y0; y < y1; y++ {
			tileData := &w.tiles[w.pt(x, y)]
			if tileData.Wallu < 0 {
				fixWall(w, x, y)
			}
			if !tileData.Active {
				continue
			}
			if tileData.U < 0 {
				fixTile(w, x, y)
			}
		}
	}
}

// https://github.com/mrkite/TerraFirma/blob/master/world.cpp#L49
func (w *World) Load(name string) {
	w.name = name
	w.info.Load()
	handle := Handle{}
	handle.LoadFile(name)
	defer handle.Close()
	w.Version = handle.U32()
	handle.Skip(8) // format
	handle.Skip(4) // revision
	handle.Skip(8) // flags
	w.numSections = handle.U16()
	w.sections = make([]uint32, w.numSections)
	for x := uint16(0); x < w.numSections; x++ {
		w.sections[x] = handle.U32()
	}
	w.numTiles = handle.U16()
	var mask uint8 = 0x80
	var bits uint8 = 0
	w.extra = make([]bool, w.numTiles)
	for x := uint16(0); x < w.numTiles; x++ {
		if mask == 0x80 {
			bits = handle.U8()
			mask = 1
		} else {
			mask <<= 1
		}
		w.extra[x] = bits&mask != 0
	}
	w.loadHeader(&handle)
	w.loadNPCs(&handle)
}

func (w *World) loadNPCs(handle *Handle) {
	var ndata []npcData
	var npcs []NPC = make([]NPC, 0, 24)
	dat, _ := ioutil.ReadFile("./data/npcs.json")
	json.Unmarshal(dat, &ndata)
	handle.Seek(int64(w.sections[4]), 0)
	var byId map[uint32]int = make(map[uint32]int)
	for n, npc := range ndata {
		byId[npc.Id] = n
	}
	for handle.U8() != 0 {
		npc := NPC{}
		if w.Version >= 190 {
			npc.Sprite = handle.U32()
			n, exists := byId[npc.Sprite]
			if exists {
				_npc := ndata[n]
				if _npc.Head.Valid {
					npc.Head = uint32(_npc.Head.Int64)
				} else {
					npc.Head = uint32(0)
				}
				npc.Title = _npc.Name
			}
		}
		npc.Name = handle.S()
		npc.X = handle.F()
		npc.Y = handle.F()
		npc.Homeless = handle.U8() != 0
		npc.HomeX = handle.U32()
		npc.HomeY = handle.U32()
		npcs = append(npcs, npc)
	}
	if w.Version >= 140 {
		for handle.U8() != 0 {
			npc := NPC{}
			if w.Version >= 190 {
				npc.Sprite = handle.U32()
				n, exists := byId[npc.Sprite]
				if exists {
					npc.Title = ndata[n].Name
				}
			}
			npc.X = handle.F()
			npc.Y = handle.F()
			npc.Homeless = true
			npcs = append(npcs, npc)
		}
	}
	w.npcs = npcs
}
