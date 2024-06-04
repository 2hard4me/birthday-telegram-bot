package ui

import (
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/datepicker"
)

func Datepicker(b *bot.Bot, onSelect datepicker.OnSelectHandler) *datepicker.DatePicker {
	// set opts for datepicker: only this years' dates
	yearStart := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(time.Now().Year(), 12, 31, 0, 0, 0, 0, time.UTC)
	opts := []datepicker.Option{
		datepicker.From(yearStart),
		datepicker.To(yearEnd),
	}

	return datepicker.New(b, onSelect, opts...)
}