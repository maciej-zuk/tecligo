package tenet

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	InitLang()
}

var lang map[string]map[string]string
var r1, r2 *regexp.Regexp

func loadInterfaceAsLangJson(jsonData interface{}) {
	cats := jsonData.(map[string]interface{})
	for cat, data := range cats {
		lang[cat] = make(map[string]string)
		vals := data.(map[string]interface{})
		for val, data2 := range vals {
			text := data2.(string)
			lang[cat][val] = text
		}
	}
}

func InitLang() {
	lang = make(map[string]map[string]string)
	files := []string{"Game.json", "Items.json", "lang.json", "Legacy.json", "NPCs.json", "Projectiles.json", "Town.json"}
	var tmp interface{}
	for _, file := range files {
		dat, _ := ioutil.ReadFile("lang/" + file)
		json.Unmarshal(dat, &tmp)
		loadInterfaceAsLangJson(tmp)
	}
	r1 = regexp.MustCompile(`{\$(.*?)}`)
	r2 = regexp.MustCompile(`{(\d+)}`)
}

func formatLang(text string, subs []string) string {
	if !strings.Contains(text, "{") {
		split := strings.Split(text, ".")
		if len(split) == 2 {
			cat, found := lang[split[0]]
			if found {
				val, found := cat[split[1]]
				if found {
					return formatLang(val, []string{})
				}
			}
		}
		return text
	}

	langReplace := func(text string) string {
		split := strings.Split(text[2:len(text)-1], ".")
		if len(split) == 2 {
			cat, found := lang[split[0]]
			if found {
				val, found := cat[split[1]]
				if found {
					return formatLang(val, []string{})
				}
			}
		}
		return text
	}

	subReplacer := func(text string) string {
		n, err := strconv.ParseInt(text[1:len(text)-1], 10, 32)
		if err != nil {
			return text
		} else {
			return formatLang(subs[n], []string{})
		}
	}

	text = r1.ReplaceAllStringFunc(text, langReplace)
	text = r2.ReplaceAllStringFunc(text, subReplacer)
	return text
}
