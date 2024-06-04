package ui

import (
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/keyboard/inline"
)

func Daypicker(b *bot.Bot, onSelect inline.OnSelect) *inline.Keyboard {
	opts := []inline.Option{
		inline.NoDeleteAfterClick(),
	}
	kb := inline.New(b, opts...).Row()
	for i := 1; i <= 31; i++ {
		kb = kb.Button(strconv.Itoa(i), []byte(strconv.Itoa(i)), onSelect)
		if i%7 == 0 {
			kb = kb.Row()
		}
	}
	for i := 0; i < 4; i++ {
		kb = kb.Button(" ", []byte(" "), onSelect)
	}
	kb = kb.Row().Button("Cancel", []byte("cancel"), onSelect)
	return kb
}