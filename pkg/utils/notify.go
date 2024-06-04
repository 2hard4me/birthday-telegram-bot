package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-telegram/bot"
)

var (
	tz *time.Location
)

func StartNotifier(ctx context.Context, b *bot.Bot) {
	tz, _ = time.LoadLocation("Asia/Kolkata")
	s := gocron.NewScheduler(tz)
	s.Every(1).Day().At("12:00").Do(notifyAllBefore, ctx, b, -1)
	s.Every(1).Day().At("00:00").Do(notifyAllBefore, ctx, b, 0)
	s.Every(1).Day().At("20:00").Do(notifyAllBefore, ctx, b, 1)
	s.Every(1).Day().At("10:00").Do(notifyAllBefore, ctx, b, 7)
	s.StartAsync()
}

func notifyBefore(ctx context.Context, b *bot.Bot, bd BirthdayInfo, before int) {
	var text string

	switch before {
	case 0:
		text = "Today"
	case 1:
		text = "Tomorrow"
	default:
		text = fmt.Sprintf("%d days from now", before)
	}

	text += fmt.Sprintf(" is %s's birthday 🎂 (%s)", bd.Name, bd.Date())

	if before == -1 {
		text = fmt.Sprintf("Did you wish %s a happy birthday?\nIf not wish now itself.", bd.Name)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: bd.ChatID,
		Text:   text,
	})
}

func notifyAllBefore(ctx context.Context, b *bot.Bot, before int) {
	var bdDate time.Time
	now := time.Now().In(tz)
	if before == -1 {
		bdDate = now
	} else {
		bdDate = now.AddDate(0, 0, before)
	}

	bds, err := GetBirthdaysOn(bdDate.Day(), int(bdDate.Month()))
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, bd := range bds {
		notifyBefore(ctx, b, bd, before)
	}
}