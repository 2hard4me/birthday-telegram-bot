package commands

import (
	"context"
	"fmt"
	"gazprom/pkg/ui"
	"gazprom/pkg/utils"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func RemoveHandler(ctx context.Context, b *bot.Bot, mes *models.Message, bdList []utils.BirthdayInfo) {
	chatID := mes.Chat.ID
	cmd := mes.ReplyToMessage.Text[1:]
	// var kb models.ReplyMarkup
	kb := ui.Birthdaypicker(b, birthdaypickerRemoveHandler, bdList)

	// switch cmd {
	// case "edit":
	// 	// kb = ui.Birthdaypicker(b, birthdaypickerEditHandler, bdList)
	// case "remove":
	// 	kb = ui.Birthdaypicker(b, birthdaypickerDeleteHandler, bdList)
	// }

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:                   chatID,
		Text:                     "Select the birthday to /" + cmd,
		ReplyMarkup:              kb,
		ReplyParameters:          &models.ReplyParameters{
			MessageID: mes.ReplyToMessage.ID,
			AllowSendingWithoutReply: true,
		},
	})
}

func birthdaypickerRemoveHandler(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	chatID := mes.Message.Chat.ID
	var text string

	name := string(data)
	if name == " " {
		return
	}

	bd := utils.BirthdayInfo{ChatID: chatID, Name: name}
	err := utils.RemoveBirthday(&bd)
	if err != nil {
		log.Println(err)
		text = utils.RetryReply("/remove")
	} else {
		text = fmt.Sprintf("<b>Removed Birthday</b> of %s", bd.Name)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	})
}