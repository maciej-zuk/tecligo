package tenet

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/maciej-zuk/tecligo/common"

	"github.com/bwmarrin/discordgo"
)

var chatRegex *regexp.Regexp

func init() {
	chatRegex = regexp.MustCompile(`\[([a-zA-Z]{1,10})(\/([a-zA-Z0-9]+))?:(.+?)\]`)
}

func discordBotRoutine(c *Connection) {
	defer c.exitWg.Done()
	dg, err := discordgo.New("Bot " + common.Settings.DiscordBotToken)
	if err != nil {
		log.Println("Error creating Discord session,", err)
		return
	}

	err = dg.Open()
	if err != nil {
		log.Println("Error opening Discord connection,", err)
		return
	}

	c.dg = dg
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		messageCreate(c, s, m)
	})

	if common.Session.DiscordBotWebhookId == "" {
		log.Println("No webhook, creating one")
		webhook, err := dg.WebhookCreate(common.Settings.DiscordBotChannelId, "Terrarian Webhook", "")
		if err != nil {
			log.Println("Error creating Discord webhook", err)
			return
		} else {
			common.Session.DiscordBotWebhookId = webhook.ID
			common.Session.DiscordBotWebhookToken = webhook.Token
			common.Session.Save()
		}
	}

	c.enterWg.Wait()
	c.dg = nil
	dg.Close()
}

func messageCreate(c *Connection, s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.WebhookID != "" {
		return
	}
	c.hub.botConnection.Send(
		82,
		TnShort(1),
		TnString("Say"),
		TnString(fmt.Sprintf("[c/00FFFF:%s]: [c/FFFF00:%s]", m.Author.Username, m.Content)),
	)
}

func chatReplace(text string) string {
	groups := chatRegex.FindStringSubmatch(text)
	if len(groups) == 5 {
		return groups[4]
	} else {
		return text
	}
}

func parseChat(text string) string {
	return chatRegex.ReplaceAllStringFunc(text, chatReplace)
}

func (c *Connection) sendDiscordChat(slot TnByte, player *Player, msg string) {
	if c.dg == nil {
		return
	}
	parsedMsg := parseChat(msg)
	if slot == 255 || slot == c.Slot {
		if strings.Contains(msg, "[c/00FFFF:") {
			return
		}
		_, err := c.dg.ChannelMessageSend(common.Settings.DiscordBotChannelId, parsedMsg)
		if err != nil {
			log.Println("Error Discord channel message send,", err)
		}
	} else {
		_, err := c.dg.WebhookExecute(
			common.Session.DiscordBotWebhookId,
			common.Session.DiscordBotWebhookToken,
			true,
			&discordgo.WebhookParams{
				Content:   parsedMsg,
				Username:  string(player.Name),
				AvatarURL: fmt.Sprintf("%splayer_%d.png?v=%d", common.Settings.ServerPath, slot, player.ImgVersion),
			},
		)
		if err != nil {
			log.Println("Error executing Discord webhook,", err)
		}
	}
}
