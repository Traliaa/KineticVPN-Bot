package telegram

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotAdapter interface {
	Start(ctx context.Context)
}

type Bot struct {
	client              *tgbotapi.BotAPI
	config              tgbotapi.UpdateConfig
	commandHandle       func(*tgbotapi.BotAPI, *tgbotapi.Message)
	messageHandle       func(bot *tgbotapi.BotAPI, message *tgbotapi.Message)
	handleCallbackQuery func(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery)
}

func NewClient(token string, commandHandle func(*tgbotapi.BotAPI, *tgbotapi.Message), messageHandle func(bot *tgbotapi.BotAPI, message *tgbotapi.Message), handleCallbackQuery func(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery)) Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	b := Bot{
		client:              bot,
		commandHandle:       commandHandle,
		messageHandle:       messageHandle,
		handleCallbackQuery: handleCallbackQuery,
	}

	// Настраиваем канал обновлений
	b.config = tgbotapi.NewUpdate(0)
	b.config.Timeout = 60

	return b
}

func (b *Bot) Start(ctx context.Context) {

	updates := b.client.GetUpdatesChan(b.config)
	// Обрабатываем входящие обновления
	for update := range updates {
		if update.Message != nil {
			b.messageHandle(b.client, update.Message)
		} else if update.CallbackQuery != nil {
			b.handleCallbackQuery(b.client, update.CallbackQuery)
		}
	}

	// Обрабатываем входящие сообщения
	//for update := range updates {
	//	if update.Message == nil {
	//		continue
	//	}
	//
	//	// Обрабатываем команды
	//	if update.Message.IsCommand() {
	//		b.commandHandle(b.client, update.Message)
	//		continue
	//	}
	//
	//	// Обрабатываем обычные сообщения (опционально)
	//	b.messageHandle(b.client, update.Message)
	//}
}

func (b *Bot) SendMessage(ctx context.Context, ChatID, text string) error {
	//_, err := b.client.SendMessage(ctx, &bot.SendMessageParams{
	//	ChatID: ChatID,
	//	Text:   text,
	//})
	//if err != nil {
	//	return err
	//}
	return nil
}
