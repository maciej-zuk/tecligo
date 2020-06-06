package temap

import (
	"github.com/veandco/go-sdl2/sdl"
	"gopkg.in/guregu/null.v4"
)

type (
	RenderingContext struct {
		X0, Y0, X1, Y1, Xoff, Yoff uint32
		Img                        *sdl.Surface
	}

	MergeBlend struct {
		HasTile   bool
		Tile      int16
		Mask      uint32
		Blend     bool
		Recursive bool
		Direction uint8
	}

	TileInfo struct {
		Name                                                             string
		Color                                                            uint32
		LightR, LightG, LightB                                           float64
		Mask                                                             uint32
		Solid, Transparent, Dirt, Stone, Grass, Pile, Flip, Brick, Merge bool
		Blends                                                           []MergeBlend
		Width, Height, Skipy, Toppad                                     int
		U, V, Minu, Maxu, Minv, Maxv                                     int
		IsHilighting                                                     bool
		Large                                                            uint8
	}

	Tile struct {
		Active       bool
		Type         int
		Wall         uint8
		Wallu, Wallv int16
		U, V         int16
		Color        uint8
		WallColor    uint8
		Liquid       uint8
		LiquidType   uint8
		Half         bool
		Inactive     bool
		Slope        uint8
	}

	World struct {
		name                 string
		tilesWide, tilesHigh uint32
		header               WorldHeader
		tiles                []Tile
		info                 TilesInfo
		numSections          uint16
		Version              uint32
		sections             []uint32
		numTiles             uint16
		extra                []bool
		groundLevel          uint32
		npcs                 []NPC
	}

	headerInfo struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		Min    int    `json:"min"`
		Relnum string `json:"relnum"`
		Num    int    `json:"num"`
	}

	// Tile info database
	TilesInfo struct {
		Walls      []TileInfo
		Tiles      []TileInfo
		WallsCount int
		TilesCount int
	}

	tileData struct {
		Id     int         `json:"id"`
		Name   null.String `json:"name"`
		Color  null.String `json:"color"`
		Flags  null.Int    `json:"flags"`
		Merge  interface{} `json:"merge"`
		Blends interface{} `json:"blend"`
		W      null.Int    `json:"w"`
		H      null.Int    `json:"h"`
		Skipy  null.Int    `json:"skipy"`
		Toppad null.Int    `json:"toppad"`
	}

	wallData struct {
		Id    uint16      `json:"id"`
		Name  null.String `json:"name"`
		Color null.String `json:"color"`
		Large null.Int    `json:"large"`
	}

	npcData struct {
		Id   uint32   `json:"id"`
		Name string   `json:"name"`
		Head null.Int `json:"head"`
	}

	NPC struct {
		Name     string
		Title    string
		X        float32
		Y        float32
		Sprite   uint32
		Homeless bool
		Head     uint32
		HomeX    uint32
		HomeY    uint32
	}
)
