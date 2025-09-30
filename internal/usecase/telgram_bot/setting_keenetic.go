package telgram_bot

import (
	"context"

	"github.com/Traliaa/KineticVPN-Bot/internal/texts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotService) askRouterURL(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *UserData) {
	if user == nil {
		UserSession[chatID] = &UserData{Step: "router_url"}
		user = UserSession[chatID]
	} else {
		user.Step = "router_url"
	}

	msg := tgbotapi.NewMessage(chatID, texts.GetTranslate(ctx, texts.RU, texts.RouterAddress))
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
