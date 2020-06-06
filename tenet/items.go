package tenet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"

	"github.com/sahilm/fuzzy"
)

var (
	giveCmdRex           *regexp.Regexp
	giveSingleCmdRex     *regexp.Regexp
	giveNameCmdRex       *regexp.Regexp
	giveNameSingleCmdRex *regexp.Regexp
)

type itemDbEntry struct {
	Name string `json:"n"`
	Type int    `json:"t"`
}

type itemDbType []itemDbEntry

var itemDb itemDbType

func (e itemDbType) String(i int) string {
	return e[i].Name
}

func (e itemDbType) Len() int {
	return len(e)
}

func init() {
	giveCmdRex = regexp.MustCompile(`give (\d+) (\d+)`)
	giveSingleCmdRex = regexp.MustCompile(`give (\d+)`)
	giveNameCmdRex = regexp.MustCompile(`give \b(\w+)\b (\d+)`)
	giveNameSingleCmdRex = regexp.MustCompile(`give \b(\w+)\b`)
	dat, _ := ioutil.ReadFile("data/items.json")
	json.Unmarshal(dat, &itemDb)
}

func (c *Connection) spawnItemForPlayer(player *Player, itemId uint16, itemStack uint16) {
	c.Send(
		21,
		TnShort(400),
		player.Pos,
		TnVector{0, 1},
		TnShort(itemStack),
		TnByte(0),
		TnByte(0),
		TnShort(itemId),
	)
	c.Send(
		39,
		TnShort(400),
	)
}

func (c *Connection) parseChat(msg TnNetString, author TnByte, color TnColor) {
	var authorName string
	player := &c.Players[author]
	if author == 255 {
		authorName = ""
	} else if player.Active {
		authorName = fmt.Sprintf("<%s> ", string(player.Name))
	} else if author == c.Slot {
		authorName = ""
	} else {
		authorName = "Unknown "
	}
	msgString := string(msg)
	if player.Active {
		groups := giveCmdRex.FindStringSubmatch(msgString)
		if len(groups) == 3 {
			itemId, err1 := strconv.Atoi(groups[1])
			itemStack, err2 := strconv.Atoi(groups[2])
			if err1 == nil && err2 == nil {
				c.spawnItemForPlayer(player, uint16(itemId), uint16(itemStack))
			}
			return
		}
		groups = giveSingleCmdRex.FindStringSubmatch(msgString)
		if len(groups) == 2 {
			itemId, err1 := strconv.Atoi(groups[1])
			if err1 == nil {
				c.spawnItemForPlayer(player, uint16(itemId), 1)
			}
			return
		}
		groups = giveNameCmdRex.FindStringSubmatch(msgString)
		if len(groups) == 3 {
			finds := fuzzy.FindFrom(groups[1], itemDb)
			itemStack, err2 := strconv.Atoi(groups[2])
			if err2 == nil && len(finds) > 0 {
				itemId := itemDb[finds[0].Index].Type
				c.spawnItemForPlayer(player, uint16(itemId), uint16(itemStack))
			}
			return
		}
		groups = giveNameSingleCmdRex.FindStringSubmatch(msgString)
		if len(groups) == 2 {
			finds := fuzzy.FindFrom(groups[1], itemDb)
			if len(finds) > 0 {
				itemId := itemDb[finds[0].Index].Type
				c.spawnItemForPlayer(player, uint16(itemId), 1)
			}
			return
		}
	}
	log.Println(authorName + msgString)
	c.broadcastChat(authorName+msgString, color)
	c.sendDiscordChat(author, player, msgString)
}
