package ui

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/keyboard/inline"
)

func Searchby(b *bot.Bot, onSelect inline.OnSelect) *inline.Keyboard {
	return inline.New(b).
		Button("List All", []byte("all"), onSelect).
		Row().
		Button("By Name", []byte("name"), onSelect).
		Button("By Date", []byte("date"), onSelect).
		Row().
		Button("By Month", []byte("month"), onSelect).
		Button("By Day", []byte("day"), onSelect).
		Row().
		Button("Cancel", []byte("cancel"), onSelect)
}