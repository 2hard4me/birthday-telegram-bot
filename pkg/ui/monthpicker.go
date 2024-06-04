package ui

import (
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/keyboard/inline"
)

func Monthpicker(b *bot.Bot, onSelect inline.OnSelect) *inline.Keyboard {
	kb := inline.New(b).Row()
	for i := 1; i <= 12; i++ {
		kb = kb.Button(time.Month(i).String(), []byte(strconv.Itoa(i)), onSelect)
		if i%4 == 0 {
			kb = kb.Row()
		}
	}
	kb = kb.Button("Cancel", []byte("cancel"), onSelect)
	return kb
}