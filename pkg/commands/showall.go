package commands

import (
	"context"
	"gazprom/pkg/utils"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func ShowallHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	var text string

	bd := utils.BirthdayInfo{ChatID: chatID}
	bdList, err := utils.SearchBirthday(&bd, "all")
	if err != nil {
		text = utils.RetryReply("/showall")
	} else {
		if len(bdList) == 0 {
			text = "<b><i>No Birthdays found</i></b>"
		} else {
			text = "<b>All Birthdays:</b>\n" + utils.BirthdayListStr(bdList)
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	})
}