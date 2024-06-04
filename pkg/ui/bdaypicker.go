package ui

import (
	"gazprom/pkg/utils"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/keyboard/inline"
)

func Birthdaypicker(b *bot.Bot, onSelect inline.OnSelect, bdList []utils.BirthdayInfo) *inline.Keyboard {
	kb := inline.New(b)
	for _, bd := range bdList {
		kb = kb.Row().Button(bd.String(), []byte(bd.Name), onSelect)
	}
	kb = kb.Row().Button("Cancel", []byte(" "), onSelect)
	return kb
}