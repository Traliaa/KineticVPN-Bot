package telgram_bot

import "github.com/Traliaa/KineticVPN-Bot/internal/adapter/telegram"

type BotService struct {
	api telegram.BotAdapter
}

func NewBotService() *BotService {
	return &BotService{}
}

func (b *BotService) StartButton() {

}

// ServiceMap для быстрого поиска по ID
var ServiceMap = make(map[string]string)

func init() {
	// Инициализируем map для быстрого доступа
	for _, service := range AvailableServices {
		ServiceMap[service.ID] = service.Name
	}
}

// ServiceDefinition определяет сервис с фиксированным порядком
type ServiceDefinition struct {
	ID   string
	Name string
}

// AvailableServices доступные сервисы для VPN в фиксированном порядке
var AvailableServices = []ServiceDefinition{
	{"youtube", "YouTube"},
	{"instagram", "Instagram"},
	{"facebook", "Facebook"},
	{"whatsapp", "WhatsApp"},
	{"telegram", "Telegram"},
	{"netflix", "Netflix"},
	{"twitter", "Twitter"},
	{"tiktok", "TikTok"},
	{"spotify", "Spotify"},
	{"discord", "Discord"},
	{"vkontakte", "VKontakte"},
	{"twitch", "Twitch"},
	{"reddit", "Reddit"},
	{"linkedin", "LinkedIn"},
	{"pinterest", "Pinterest"},
}
