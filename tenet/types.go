package tenet

import (
	"bytes"
	"net"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type (
	Payload interface {
		Serialize(*bytes.Buffer)
		Size() int
	}
	TnByte      byte
	TnShort     int16
	TnInt       int32
	TnFloat     float32
	TnString    string
	TnNetString string

	TnColor struct {
		R, G, B uint8
	}

	TnVector struct {
		X, Y float32
	}

	TnArray struct {
		Count int
		Seed  Payload
	}

	InventorySlot struct {
		Item   TnShort
		Stack  TnShort
		Prefix TnByte
	}

	Player struct {
		Active      bool               `json:"active"`
		Name        TnString           `json:"name"`
		Pos         TnVector           `json:"pos"`
		Life        TnShort            `json:"life"`
		Mana        TnShort            `json:"mana"`
		MaxLife     TnShort            `json:"maxLife"`
		MaxMana     TnShort            `json:"maxMana"`
		Hair        TnByte             `json:"-"`
		HairColor   TnColor            `json:"-"`
		SkinColor   TnColor            `json:"-"`
		EyeColor    TnColor            `json:"-"`
		Inventory   [256]InventorySlot `json:"-"`
		ImgVersion  byte               `json:"imgVersion"`
		FacingRight bool               `json:"facingRight"`
	}

	Connection struct {
		Slot      TnByte
		Players   []Player
		Password  string
		WorldInfo NetWorldInfo

		running     bool
		exitWakeUp  chan bool
		spawned     bool
		conn        net.Conn
		writeBuffer bytes.Buffer
		dataReader  *bytes.Reader
		writeMutex  sync.Mutex
		enterWg     sync.WaitGroup
		exitWg      sync.WaitGroup
		hub         *websocketHub
		dg          *discordgo.Session
	}

	NetWorldInfo struct {
		Time                   TnInt
		Flags                  TnByte
		MoonPhase              TnByte
		MaxTilesX              TnShort
		MaxTilesY              TnShort
		SpawnTileX             TnShort
		SpawnTileY             TnShort
		WorldSurface           TnShort
		RockLayer              TnShort
		WorldID                TnInt
		WorldName              TnString
		UniqueId               [16]TnByte
		WorldGeneratorVersion1 TnInt
		WorldGeneratorVersion2 TnInt
		Moon                   TnByte
		MoonType               TnByte
		TreeBG                 TnByte
		CorruptBG              TnByte
		JungleBG               TnByte
		SnowBG                 TnByte
		HallowBG               TnByte
		CrimsonBG              TnByte
		DesertBG               TnByte
		OceanBG                TnByte
		IceBackStyle           TnByte
		JungleBackStyle        TnByte
		HellBackStyle          TnByte
		WindSpeedSet           TnFloat
		NumClouds              TnByte
		TreeX                  [3]TnInt
		TreeStyle              [4]TnByte
		CaveBackX              [3]TnInt
		CaveBackStyle          [4]TnByte
		MaxRaining             TnFloat
		BitsByte5              TnByte
		BitsByte6              TnByte
		BitsByte7              TnByte
		BitsByte8              TnByte
		BitsByte9              TnByte
		Invasion               TnByte
		InvasionType           TnByte
		LobbyId1               TnInt
		LobbyId2               TnInt
		IntendedSeverity       TnFloat
	}
)
