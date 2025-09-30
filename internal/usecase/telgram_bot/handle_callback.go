package telgram_bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Traliaa/KineticVPN-Bot/internal/texts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleCallbackQuery обрабатывает нажатия на inline-кнопки
func (b *BotService) HandleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	chatID := callbackQuery.Message.Chat.ID
	data := callbackQuery.Data

	// Инициализируем сессию если её нет
	if _, exists := UserSession[chatID]; !exists {
		UserSession[chatID] = &UserData{Step: "start"}
	}

	user := UserSession[chatID]
	if user == nil {
		// Создаем новую сессию если что-то пошло не так
		UserSession[chatID] = &UserData{Step: "start"}
		user = UserSession[chatID]
	}

	switch {
	case data == "btn_setup_now":
		b.askRouterURL(context.Background(), bot, chatID, user)
	case data == "btn_setup_later":
		sendMainMenu(bot, chatID, "Хорошо! Вы можете настроить VPN в любое время через меню.")
	case strings.HasPrefix(data, "service_"):
		handleServiceSelection(bot, callbackQuery, user)
	case data == "btn_save_settings":
		saveSettings(bot, chatID, user)
	case data == "btn_back_to_main":
		sendWelcomeMessage(bot, chatID, user)
	}

	// Ответ на callback query
	callbackConfig := tgbotapi.NewCallback(callbackQuery.ID, "")
	bot.Send(callbackConfig)
}

func sendWelcomeMessage(bot *tgbotapi.BotAPI, chatID int64, user *UserData) {
	if user == nil {
		// Создаем новую сессию если user nil
		UserSession[chatID] = &UserData{Step: "ask_setup"}
		user = UserSession[chatID]
	} else {
		user.Step = "ask_setup"
	}

	text := texts.GetTranslate(context.Background(), texts.RU, texts.WelcomeMain)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, настроить", "btn_setup_now"),
			tgbotapi.NewInlineKeyboardButtonData("⏰ Позже", "btn_setup_later"),
		),
	)

	bot.Send(msg)
}

func (b *BotService) handleSetupResponse(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *UserData) {
	if user == nil {
		sendWelcomeMessage(bot, message.Chat.ID, nil)
		return
	}

	response := strings.ToLower(message.Text)

	if strings.Contains(response, "да") || strings.Contains(response, "yes") || strings.Contains(response, "настроить") {
		b.askRouterURL(context.Background(), bot, message.Chat.ID, user)
	} else {
		sendMainMenu(bot, message.Chat.ID, "Хорошо! Вы можете настроить VPN в любое время через меню.")
	}
}

func handleRouterURL(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *UserData) {
	if user == nil {
		sendWelcomeMessage(bot, message.Chat.ID, nil)
		return
	}

	user.RouterURL = message.Text
	user.Step = "username"

	text := `👤 **Шаг 2 из 5: Логин роутера**

Введите имя пользователя для доступа к роутеру.

_Обычно это "admin" или указано на наклейке роутера_`

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func handleUsername(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *UserData) {
	if user == nil {
		sendWelcomeMessage(bot, message.Chat.ID, nil)
		return
	}

	user.Username = message.Text
	user.Step = "password"

	text := `🔐 **Шаг 3 из 5: Пароль роутера**

Введите пароль для доступа к роутеру.

_Обычно это "admin" или указан на наклейке роутера_`

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func handlePassword(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *UserData) {
	if user == nil {
		sendWelcomeMessage(bot, message.Chat.ID, nil)
		return
	}

	user.Password = message.Text
	user.Step = "auth_code"

	text := `🔑 **Шаг 4 из 5: Код авторизации VPN**

Введите код для авторизации в VPN сервисе.

_Этот код вам должен был предоставить ваш VPN-провайдер_`

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func handleAuthCode(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *UserData) {
	if user == nil {
		sendWelcomeMessage(bot, message.Chat.ID, nil)
		return
	}

	user.AuthCode = message.Text
	user.Step = "service_selection"

	askServiceSelection(bot, message.Chat.ID, user)
}

func askServiceSelection(bot *tgbotapi.BotAPI, chatID int64, user *UserData) {
	if user == nil {
		sendWelcomeMessage(bot, chatID, nil)
		return
	}

	text := `📱 **Шаг 5 из 5: Выбор сервисов**

Выберите сервисы, которые должны работать через VPN:

_✅ - будет использовать VPN_
_❌ - будет работать напрямую_

Текущий выбор:`

	// Добавляем информацию о выбранных сервисах
	if len(user.SelectedApps) > 0 {
		text += "\n\nВыбрано: " + strings.Join(user.SelectedApps, ", ")
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = createServicesKeyboard(user.SelectedApps)

	bot.Send(msg)
}

func handleServiceSelection(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, user *UserData) {
	if user == nil {
		return
	}

	serviceID := strings.TrimPrefix(callbackQuery.Data, "service_")
	serviceName := ServiceMap[serviceID]

	if serviceName == "" {
		return
	}

	// Переключаем выбор сервиса
	found := false
	for i, app := range user.SelectedApps {
		if app == serviceName {
			// Удаляем если уже есть
			user.SelectedApps = append(user.SelectedApps[:i], user.SelectedApps[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		// Добавляем если нет
		user.SelectedApps = append(user.SelectedApps, serviceName)
	}

	// Обновляем сообщение с новым состоянием
	editMsg := tgbotapi.NewEditMessageTextAndMarkup(
		callbackQuery.Message.Chat.ID,
		callbackQuery.Message.MessageID,
		fmt.Sprintf(`📱 **Шаг 5 из 5: Выбор сервисов**

Выберите сервисы, которые должны работать через VPN:

_✅ - будет использовать VPN_
_❌ - будет работать напрямую_

Текущий выбор: %s`,
			func() string {
				if len(user.SelectedApps) == 0 {
					return "не выбрано"
				}
				return strings.Join(user.SelectedApps, ", ")
			}()),
		createServicesKeyboard(user.SelectedApps),
	)
	editMsg.ParseMode = "Markdown"

	bot.Send(editMsg)
}

func createServicesKeyboard(selectedApps []string) tgbotapi.InlineKeyboardMarkup {
	// Создаем карту выбранных сервисов для быстрого поиска
	selectedMap := make(map[string]bool)
	for _, app := range selectedApps {
		selectedMap[app] = true
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	// Создаем кнопки для сервисов в фиксированном порядке
	var currentRow []tgbotapi.InlineKeyboardButton

	for i, service := range AvailableServices {
		status := "❌"
		if selectedMap[service.Name] {
			status = "✅"
		}

		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", status, service.Name),
			"service_"+service.ID,
		)

		currentRow = append(currentRow, button)

		// Размещаем по 3 кнопки в строке для лучшего отображения
		if (i+1)%3 == 0 || i == len(AvailableServices)-1 {
			rows = append(rows, currentRow)
			currentRow = []tgbotapi.InlineKeyboardButton{}
		}
	}

	// Добавляем кнопки сохранения
	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("💾 Сохранить настройки", "btn_save_settings"),
	})

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("⬅️ На главную", "btn_back_to_main"),
	})

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// Альтернативный вариант с группировкой по категориям
func createServicesKeyboardGrouped(selectedApps []string) tgbotapi.InlineKeyboardMarkup {
	selectedMap := make(map[string]bool)
	for _, app := range selectedApps {
		selectedMap[app] = true
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	// Социальные сети
	socialServices := []ServiceDefinition{
		{"facebook", "Facebook"},
		{"instagram", "Instagram"},
		{"twitter", "Twitter"},
		{"vkontakte", "VKontakte"},
		{"linkedin", "LinkedIn"},
		{"pinterest", "Pinterest"},
	}

	// Медиа и развлечения
	mediaServices := []ServiceDefinition{
		{"youtube", "YouTube"},
		{"netflix", "Netflix"},
		{"tiktok", "TikTok"},
		{"twitch", "Twitch"},
		{"spotify", "Spotify"},
		{"reddit", "Reddit"},
	}

	// Мессенджеры и коммуникация
	communicationServices := []ServiceDefinition{
		{"whatsapp", "WhatsApp"},
		{"telegram", "Telegram"},
		{"discord", "Discord"},
	}

	// Создаем ряды для каждой категории
	rows = append(rows, createServiceRow("📱 Социальные сети", socialServices, selectedMap))
	rows = append(rows, createServiceRow("🎬 Медиа и развлечения", mediaServices, selectedMap))
	rows = append(rows, createServiceRow("💬 Мессенджеры", communicationServices, selectedMap))

	// Добавляем кнопки действий
	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("✅ Выбрать все", "btn_select_all"),
		tgbotapi.NewInlineKeyboardButtonData("❌ Очистить все", "btn_clear_all"),
	})

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("💾 Сохранить настройки", "btn_save_settings"),
	})

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("⬅️ На главную", "btn_back_to_main"),
	})

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func createServiceRow(category string, services []ServiceDefinition, selectedMap map[string]bool) []tgbotapi.InlineKeyboardButton {
	var row []tgbotapi.InlineKeyboardButton

	for _, service := range services {
		status := "❌"
		if selectedMap[service.Name] {
			status = "✅"
		}

		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", status, service.Name),
			"service_"+service.ID,
		)
		row = append(row, button)
	}

	return row
}

func saveSettings(bot *tgbotapi.BotAPI, chatID int64, user *UserData) {
	if user == nil {
		sendWelcomeMessage(bot, chatID, nil)
		return
	}

	// Сохраняем настройки
	user.Step = "completed"

	// Маскируем пароль звездочками
	maskedPassword := strings.Repeat("*", len(user.Password))

	// Формируем сводку настроек
	summary := fmt.Sprintf(`🎉 **Настройка завершена!**

📋 **Сводка настроек:**
• 🔗 Роутер: %s
• 👤 Логин: %s
• 🔐 Пароль: %s
• 🔑 Код VPN: %s
• 📱 Сервисы через VPN: %s

Настройки сохранены! Бот готов к работе.`,
		user.RouterURL,
		user.Username,
		maskedPassword,
		user.AuthCode,
		strings.Join(user.SelectedApps, ", "),
	)

	// Логируем настройки
	userJSON, _ := json.MarshalIndent(user, "", "  ")
	log.Printf("User %d settings: %s", chatID, userJSON)

	msg := tgbotapi.NewMessage(chatID, summary)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = createMainMenuKeyboard()

	bot.Send(msg)
}

func sendHelpMessage(bot *tgbotapi.BotAPI, chatID int64) {
	text := `📖 **Помощь по VPN Боту**

🤔 **Как это работает:**
1. Бот настраивает VPN на вашем роутере
2. Вы выбираете какие сервисы использовать через VPN
3. Остальной трафик идет напрямую

🔧 **Команды:**
/start - начать настройку
/help - эта справка  
/reset - сбросить настройки

💡 **Советы:**
• Данные роутера хранятся только во время настройки
• Вы можете изменить настройки в любое время`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func sendMainMenu(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = createMainMenuKeyboard()
	bot.Send(msg)
}

func createMainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/start"),
			tgbotapi.NewKeyboardButton("/help"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/reset"),
		),
	)
}
