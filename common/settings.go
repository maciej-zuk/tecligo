package common

import (
	"encoding/json"
	"io/ioutil"
)

type TeSettings struct {
	BasePath       string `json:"basePath"`
	ServerPath     string `json:"serverPath"`
	SocketPort     int    `json:"socketPort"`
	SocketEndpoint string `json:"socketEndpoint"`
	TerrariaHost   string `json:"terrariaHost"`
	TerrariaPort   int    `json:"terrariaPort"`
	TerrariaPass   string `json:"terrariaPass"`
	MapPath        string `json:"mapPath"`
	MapRendering   bool   `json:"mapRendering"`
	MapRegion      struct {
		X0 int `json:"x0"`
		Y0 int `json:"y0"`
		X1 int `json:"x1"`
		Y1 int `json:"y1"`
	} `json:"mapRegion"`
	TileOutputPath      string `json:"tileOutputPath"`
	DiscordBot          bool   `json:"discordBot"`
	DiscordBotToken     string `json:"discordBotToken"`
	DiscordBotChannelId string `json:"discordBotChannelId"`
}

type TeSession struct {
	DiscordBotWebhookId    string
	DiscordBotWebhookToken string
}

var Settings TeSettings
var Session TeSession

func init() {
	LoadSession()
	dat, err := ioutil.ReadFile("./settings.json")
	ok := true
	if err == nil {
		err = json.Unmarshal(dat, &Settings)
		if err != nil {
			ok = false
		}
	} else {
		ok = false
	}
	if !ok {
		panic("Unable to read and parse settings.json")
	}
}

func LoadSession() {
	dat, err := ioutil.ReadFile("./session.json")
	if err == nil {
		err = json.Unmarshal(dat, &Session)
	}
}

func (session *TeSession) Save() {
	file, _ := json.Marshal(Session)
	_ = ioutil.WriteFile("./session.json", file, 0644)
}
