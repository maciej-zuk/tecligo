package tenet

import (
	"log"
)

var netMsgHandler map[uint8]func(c *Connection) = make(map[uint8]func(c *Connection))

func init() {
	InitNetMsg()
}

func noop(c *Connection) {
}

func InitNetMsg() {
	netMsgHandler[3] = func(c *Connection) {
		c.Slot = DeSerializeByte(c.dataReader)
		c.Send(
			4,
			c.Slot,
			TnByte(0),         // skin
			TnByte(0),         // hair
			TnString("GoBot"), // name
			TnByte(0),         // hair2
			TnByte(0),         // ?
			TnByte(0),         // ?
			TnByte(0),         // ?
			TnColor{0, 0, 0},  // hair
			TnColor{0, 0, 0},  // skin
			TnColor{0, 0, 0},  // eye
			TnColor{0, 0, 0},  // shirt
			TnColor{0, 0, 0},  // shirt2
			TnColor{0, 0, 0},  // pants
			TnColor{0, 0, 0},  // shoe
			TnByte(0),         // flags
		)
		c.Send(
			16,
			c.Slot,
			TnShort(500),
			TnShort(500),
		)
		c.Send(
			42,
			c.Slot,
			TnShort(5),
			TnShort(5),
		)
		c.Send(6)
	}

	netMsgHandler[7] = func(c *Connection) { // world info
		nw := NetWorldInfo{}
		nw.Time = DeSerializeInt(c.dataReader)
		nw.Flags = DeSerializeByte(c.dataReader)
		nw.MoonPhase = DeSerializeByte(c.dataReader)
		nw.MaxTilesX = DeSerializeShort(c.dataReader)
		nw.MaxTilesY = DeSerializeShort(c.dataReader)
		nw.SpawnTileX = DeSerializeShort(c.dataReader)
		nw.SpawnTileY = DeSerializeShort(c.dataReader)
		nw.WorldSurface = DeSerializeShort(c.dataReader)
		nw.RockLayer = DeSerializeShort(c.dataReader)
		nw.WorldID = DeSerializeInt(c.dataReader)
		nw.WorldName = DeSerializeString(c.dataReader)
		for x := 0; x < 16; x++ {
			nw.UniqueId[x] = DeSerializeByte(c.dataReader)
		}
		nw.WorldGeneratorVersion1 = DeSerializeInt(c.dataReader)
		nw.WorldGeneratorVersion2 = DeSerializeInt(c.dataReader)
		nw.MoonType = DeSerializeByte(c.dataReader)
		nw.TreeBG = DeSerializeByte(c.dataReader)
		nw.CorruptBG = DeSerializeByte(c.dataReader)
		nw.JungleBG = DeSerializeByte(c.dataReader)
		nw.SnowBG = DeSerializeByte(c.dataReader)
		nw.HallowBG = DeSerializeByte(c.dataReader)
		nw.CrimsonBG = DeSerializeByte(c.dataReader)
		nw.DesertBG = DeSerializeByte(c.dataReader)
		nw.OceanBG = DeSerializeByte(c.dataReader)
		nw.IceBackStyle = DeSerializeByte(c.dataReader)
		nw.JungleBackStyle = DeSerializeByte(c.dataReader)
		nw.HellBackStyle = DeSerializeByte(c.dataReader)
		nw.WindSpeedSet = DeSerializeFloat(c.dataReader)
		nw.NumClouds = DeSerializeByte(c.dataReader)
		for x := 0; x < 3; x++ {
			nw.TreeX[x] = DeSerializeInt(c.dataReader)
		}
		for x := 0; x < 4; x++ {
			nw.TreeStyle[x] = DeSerializeByte(c.dataReader)
		}
		for x := 0; x < 3; x++ {
			nw.CaveBackX[x] = DeSerializeInt(c.dataReader)
		}
		for x := 0; x < 4; x++ {
			nw.CaveBackStyle[x] = DeSerializeByte(c.dataReader)
		}
		nw.MaxRaining = DeSerializeFloat(c.dataReader)
		nw.BitsByte5 = DeSerializeByte(c.dataReader)
		nw.BitsByte6 = DeSerializeByte(c.dataReader)
		nw.BitsByte7 = DeSerializeByte(c.dataReader)
		nw.BitsByte8 = DeSerializeByte(c.dataReader)
		nw.BitsByte9 = DeSerializeByte(c.dataReader)
		nw.InvasionType = DeSerializeByte(c.dataReader)
		nw.LobbyId1 = DeSerializeInt(c.dataReader)
		nw.LobbyId2 = DeSerializeInt(c.dataReader)
		nw.IntendedSeverity = DeSerializeFloat(c.dataReader)
		c.WorldInfo = nw
		if c.spawned == false {
			c.Send(8,
				nw.SpawnTileX,
				nw.SpawnTileY,
			)
		}
		c.notifyWorldUpdate()
	}

	netMsgHandler[49] = func(c *Connection) { // spawn request
		if c.spawned == false {
			c.Send(12,
				c.Slot,
				c.WorldInfo.SpawnTileX,
				c.WorldInfo.SpawnTileY,
			)
			c.spawned = true
			log.Println("Spawned!")
			c.Send(
				13,
				c.Slot,
				TnByte(0),
				TnByte(64),
				TnByte(1),
				TnByte(1),
				TnByte(0),
				TnVector{16 * float32(c.WorldInfo.SpawnTileX), float32(16 * (c.WorldInfo.MaxTilesY - 100))},
			)
		}
	}

	netMsgHandler[2] = func(c *Connection) { //kick
		log.Println("Kicked: ", DeSerializeNetString(c.dataReader))
		c.Close()
	}

	netMsgHandler[4] = func(c *Connection) { // player info
		slot := DeSerializeByte(c.dataReader)
		DeSerializeByte(c.dataReader)               // style variant
		hair := DeSerializeByte(c.dataReader)       // hair
		name := DeSerializeString(c.dataReader)     // name
		DeSerializeByte(c.dataReader)               // hair dye
		DeSerializeByte(c.dataReader)               // b1 hide flags
		DeSerializeByte(c.dataReader)               // b2 hide flags
		DeSerializeByte(c.dataReader)               // b3 hide flags
		hairColor := DeSerializeColor(c.dataReader) // hair3
		skinColor := DeSerializeColor(c.dataReader) // skin2
		eyeColor := DeSerializeColor(c.dataReader)  // eye
		DeSerializeColor(c.dataReader)              // shirt
		DeSerializeColor(c.dataReader)              // shirt2
		DeSerializeColor(c.dataReader)              // pants
		DeSerializeColor(c.dataReader)              // shoe
		DeSerializeByte(c.dataReader)               // flags
		c.Players[slot].Name = name
		c.Players[slot].Hair = hair + 1
		c.Players[slot].HairColor = hairColor
		c.Players[slot].SkinColor = skinColor
		c.Players[slot].EyeColor = eyeColor
		savePlayerHeadIcon(slot, &c.Players[slot])
	}

	netMsgHandler[5] = func(c *Connection) { // player inventory
		slot := DeSerializeByte(c.dataReader)
		invSlot := DeSerializeByte(c.dataReader)
		stack := DeSerializeShort(c.dataReader)
		prefix := DeSerializeByte(c.dataReader)
		item := DeSerializeShort(c.dataReader)

		inv := &(c.Players[slot].Inventory[invSlot])
		inv.Item = item
		inv.Prefix = prefix
		inv.Stack = stack
		if invSlot == 59 || invSlot == 69 || invSlot == 60 || invSlot == 70 || invSlot == 61 || invSlot == 71 {
			savePlayerHeadIcon(slot, &c.Players[slot])
			c.Players[slot].ImgVersion++
			c.notifyPlayerUpdate(slot)
		}
		// 59 - helmet
		// 69 - vanity helmet
	}

	netMsgHandler[9] = noop  // status update
	netMsgHandler[10] = noop // compressed tile data
	netMsgHandler[11] = noop // recalc uv
	netMsgHandler[12] = noop // spawn

	netMsgHandler[13] = func(c *Connection) { // update player
		slot := DeSerializeByte(c.dataReader)
		flags1 := DeSerializeByte(c.dataReader) // flags1
		DeSerializeByte(c.dataReader)           // flags2
		DeSerializeByte(c.dataReader)           // flags3
		DeSerializeByte(c.dataReader)           // flags4
		DeSerializeByte(c.dataReader)           // selItem
		pos := DeSerializeVector(c.dataReader)
		c.Players[slot].Pos = pos
		c.Players[slot].FacingRight = flags1&64 != 0
		c.notifyPlayerUpdate(slot)
	}

	netMsgHandler[14] = func(c *Connection) { // set player active
		slot := DeSerializeByte(c.dataReader)
		active := DeSerializeByte(c.dataReader)
		c.Players[slot].Active = active != 0
		c.notifyPlayerUpdate(slot)
	}

	netMsgHandler[16] = func(c *Connection) { // life update
		slot := DeSerializeByte(c.dataReader)
		life := DeSerializeShort(c.dataReader)
		maxLife := DeSerializeShort(c.dataReader)
		c.Players[slot].Life = life
		c.Players[slot].MaxLife = maxLife
	}

	netMsgHandler[17] = noop // modify tile
	netMsgHandler[19] = noop // use doors
	netMsgHandler[20] = noop // modify tile block
	netMsgHandler[21] = noop // item data
	netMsgHandler[33] = noop // chest lock
	netMsgHandler[34] = noop // chest remove
	netMsgHandler[22] = noop // item owner
	netMsgHandler[23] = noop // npc data
	netMsgHandler[27] = noop // update proj
	netMsgHandler[28] = noop // npc strike
	netMsgHandler[29] = noop // destroy proj
	netMsgHandler[30] = noop // npc hostile
	netMsgHandler[35] = noop // heal
	netMsgHandler[36] = noop // set zone

	netMsgHandler[37] = func(c *Connection) { // pass required
		if c.Password == "" {
			log.Println("PASSWORD REQUIRED!")
			c.Close()
		} else {
			c.Send(
				38,
				TnString(c.Password),
			)
		}
	}

	netMsgHandler[39] = func(c *Connection) { // item deown request
		item := DeSerializeShort(c.dataReader)
		c.Send(
			22,
			item,
			TnByte(255),
		)
	}

	netMsgHandler[40] = noop // talk
	netMsgHandler[41] = noop // item animation

	netMsgHandler[42] = func(c *Connection) { // mana update
		slot := DeSerializeByte(c.dataReader)
		mana := DeSerializeShort(c.dataReader)
		maxMana := DeSerializeShort(c.dataReader)
		c.Players[slot].Mana = mana
		c.Players[slot].MaxMana = maxMana
	}

	netMsgHandler[43] = noop // mana add
	netMsgHandler[45] = noop // npc team
	netMsgHandler[47] = noop // sign
	netMsgHandler[48] = noop // liquid change
	netMsgHandler[50] = noop // npc buffs
	netMsgHandler[51] = noop // npc special effect
	netMsgHandler[52] = noop // unlock chest/door
	netMsgHandler[54] = noop // buff
	netMsgHandler[56] = noop // npc data 2
	netMsgHandler[57] = noop // world balance
	netMsgHandler[58] = noop // play harp
	netMsgHandler[59] = noop // switch flip
	netMsgHandler[60] = noop // npc house info
	netMsgHandler[65] = noop // teleport
	netMsgHandler[66] = noop // heal
	netMsgHandler[72] = noop // traveler shop
	netMsgHandler[74] = noop // angler quest status
	netMsgHandler[78] = noop // event progress
	netMsgHandler[79] = noop // item placement
	netMsgHandler[80] = noop // chest stuff

	netMsgHandler[82] = func(c *Connection) { //subsystem msg
		cmd := DeSerializeShort(c.dataReader)
		if cmd == 1 { //chat subsystem
			author := DeSerializeByte(c.dataReader)
			text := DeSerializeNetString(c.dataReader)
			color := DeSerializeColor(c.dataReader)
			c.parseChat(text, author, color)
		}
	}

	netMsgHandler[83] = noop  // kills count
	netMsgHandler[84] = noop  // npc stealth
	netMsgHandler[86] = noop  // trainig dummy stuff?
	netMsgHandler[88] = noop  // projectile stuff?
	netMsgHandler[91] = noop  // npc emoticon bubble
	netMsgHandler[98] = noop  // event progress
	netMsgHandler[101] = noop // pilar status
	netMsgHandler[103] = noop // moonlord countdown
	netMsgHandler[106] = noop // poof of smoke
	netMsgHandler[107] = noop // announcement
	netMsgHandler[112] = noop // tree growth
	netMsgHandler[114] = noop // DD2 WipeEntities
	netMsgHandler[115] = noop // minion target
	netMsgHandler[116] = noop // DD2 SetEnemySpawningOnHold
	netMsgHandler[117] = noop // npc hurt
	netMsgHandler[118] = noop // npc death
	netMsgHandler[125] = noop // tile damage
	netMsgHandler[129] = noop // spawn preparation
	netMsgHandler[131] = noop // immune time
	netMsgHandler[134] = noop // sync luck
	netMsgHandler[135] = noop // player dead
	netMsgHandler[136] = noop // cavern monster type
	netMsgHandler[139] = noop // host player slot
}
