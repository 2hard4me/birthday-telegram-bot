package commands

import (
	"context"
	"fmt"
	"gazprom/pkg/ui"
	"gazprom/pkg/utils"
	"log"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func SearchHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	kb := ui.Searchby(b, searchbyHandler)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           chatID,
		Text:             "Select the search type",
		ReplyMarkup:      kb,
		ReplyParameters:  &models.ReplyParameters{MessageID: update.Message.ID},
	})
}

func searchbyHandler(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	chatID := mes.Message.Chat.ID
	var text string
	var kb models.ReplyMarkup

	switch string(data) {
	case "all":
		Search(ctx, b, mes.Message, &utils.BirthdayInfo{ChatID: chatID}, "all")
		return
	case "name":
		text = "Enter the name of the person to " + mes.Message.ReplyToMessage.Text
		kb = models.ForceReply{ForceReply: true}
	case "date":
		text = "Select the date"
		kb = ui.Datepicker(b, datepickerSearchHandler)
	case "month":
		text = "Select the month"
		kb = ui.Monthpicker(b, monthpickerSearchHandler)
	case "day":
		text = "Select the day"
		kb = ui.Daypicker(b, daypickerSearchHandler)
	case "cancel":
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           chatID,
		Text:             text,
		ReplyMarkup:      kb,
		ReplyParameters:  &models.ReplyParameters{MessageID: mes.Message.ReplyToMessage.ID},
	})
}

func datepickerSearchHandler(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, date time.Time) {
	chatID := mes.Message.Chat.ID

	bd := utils.BirthdayInfo{
		ChatID: chatID,
		Day:    date.Day(),
		Month:  int(date.Month()),
	}
	Search(ctx, b, mes.Message, &bd, "date")
}

func monthpickerSearchHandler(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	chatID := mes.Message.Chat.ID

	if string(data) == "cancel" {
		return
	}

	month, err := strconv.Atoi(string(data))
	if err != nil {
		fmt.Println(err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:                   chatID,
			Text:                     utils.RetryReply("/search"),
			ReplyParameters:          &models.ReplyParameters{
				MessageID: mes.Message.ReplyToMessage.ID,
				AllowSendingWithoutReply: true,
			},
		})
		return
	}

	bd := utils.BirthdayInfo{
		ChatID: chatID,
		Month:  month,
	}
	Search(ctx, b, mes.Message, &bd, "month")
}

func daypickerSearchHandler(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	chatID := mes.Message.Chat.ID
	var day int

	switch string(data) {
	case "cancel":
		b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    chatID,
			MessageID: mes.Message.ID,
		})
		return
	case " ":
		return
	default:
		var err error
		day, err = strconv.Atoi(string(data))
		if err != nil {
			log.Println(err)
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    chatID,
				MessageID: mes.Message.ID,
				Text:      utils.RetryReply("/search"),
			})
			return
		} else {
			b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    chatID,
				MessageID: mes.Message.ID,
			})
		}
	}

	bd := utils.BirthdayInfo{
		ChatID: chatID,
		Day:    day,
	}
	Search(ctx, b, mes.Message, &bd, "day")
}

func Search(ctx context.Context, b *bot.Bot, mes *models.Message, bd *utils.BirthdayInfo, searchby string) {
	chatID := mes.Chat.ID
	cmd := mes.ReplyToMessage.Text[1:]
	var text string

	searchResults, err := utils.SearchBirthday(bd, searchby)
	if err != nil {
		log.Println(err)
		text = utils.RetryReply("/" + cmd)
	} else {
		if len(searchResults) == 0 {
			text = "<b><i>No Birthdays found</i></b>"
		} else if cmd != "search" {
			b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    chatID,
				MessageID: mes.ID,
			})
			RemoveHandler(ctx, b, mes, searchResults)
			return
		} else {
			text = "<b>Search Results:</b>\n" + utils.BirthdayListStr(searchResults)
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	})
}