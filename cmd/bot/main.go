package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Traliaa/KineticVPN-Bot/internal/adapter/telegram"
	"github.com/Traliaa/KineticVPN-Bot/internal/config"
	"github.com/Traliaa/KineticVPN-Bot/internal/pg/user_settings"
	"github.com/Traliaa/KineticVPN-Bot/internal/usecase/telgram_bot"
)

const (
	baseURL  = "https://rci.tankhome.netcraze.pro"
	username = "api"
	password = "Demon0203"
	timeout  = 30 * time.Second
)

type KeeneticClient struct {
	client   *http.Client
	baseURL  string
	username string
	password string
}

type APIError struct {
	Operation  string
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: status %d - %s", e.Operation, e.StatusCode, e.Message)
}

func NewKeeneticClient() *KeeneticClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &KeeneticClient{
		client: &http.Client{
			Transport: tr,
			Timeout:   timeout,
		},
		baseURL:  baseURL,
		username: username,
		password: password,
	}
}

// Безопасное извлечение значения из map
func getStringSafe(data map[string]interface{}, keys ...string) string {
	current := data
	for i, key := range keys {
		if val, exists := current[key]; exists {
			if i == len(keys)-1 {
				if str, ok := val.(string); ok {
					return str
				}
				return fmt.Sprintf("%v", val)
			} else {
				if nextMap, ok := val.(map[string]interface{}); ok {
					current = nextMap
				} else {
					break
				}
			}
		} else {
			break
		}
	}
	return "N/A"
}

// Выполнение RCI команды GET
func (kc *KeeneticClient) ExecuteRCI(command string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", kc.baseURL+"/rci/"+command, nil)
	if err != nil {
		return nil, &APIError{Operation: "Create Request", Message: err.Error()}
	}

	req.SetBasicAuth(kc.username, kc.password)

	resp, err := kc.client.Do(req)
	if err != nil {
		return nil, &APIError{Operation: "HTTP Request", Message: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &APIError{
			Operation:  "RCI Command",
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("command failed: %s", string(body)),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &APIError{Operation: "Read Body", Message: err.Error()}
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, &APIError{Operation: "JSON Parse", Message: err.Error()}
	}

	return result, nil
}

// Специализированные методы
func (kc *KeeneticClient) GetSystemInfo() (map[string]interface{}, error) {
	return kc.ExecuteRCI("show/system")
}

func (kc *KeeneticClient) GetInterfaceStatus(interfaceName string) (map[string]interface{}, error) {
	return kc.ExecuteRCI("show/interface/" + interfaceName)
}

func (kc *KeeneticClient) RestartInterface(interfaceName string) error {
	_, err := kc.ExecuteRCI("interface/" + interfaceName + "/restart")
	return err
}

func (kc *KeeneticClient) GetWireGuardStatus() (map[string]interface{}, error) {
	return kc.ExecuteRCI("show/interface/WireGuard1")
}

// Безопасные методы для получения конкретных значений
func (kc *KeeneticClient) GetSystemModel() (string, error) {
	data, err := kc.GetSystemInfo()
	if err != nil {
		return "", err
	}
	return getStringSafe(data, "system", "model"), nil
}

func (kc *KeeneticClient) GetSystemHostname() (string, error) {
	data, err := kc.GetSystemInfo()
	if err != nil {
		return "", err
	}
	return getStringSafe(data, "hostname"), nil
}

func (kc *KeeneticClient) GetWireGuardState() (string, error) {
	data, err := kc.GetWireGuardStatus()
	if err != nil {
		return "", err
	}
	return getStringSafe(data, "state"), nil
}

// Функции для форматирования системной информации
func (kc *KeeneticClient) FormatUptime(seconds string) string {
	sec, err := strconv.Atoi(seconds)
	if err != nil {
		return seconds + " сек"
	}

	days := sec / 86400
	hours := (sec % 86400) / 3600
	minutes := (sec % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%dд %dч %dм", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dч %dм", hours, minutes)
	} else {
		return fmt.Sprintf("%dм", minutes)
	}
}

func (kc *KeeneticClient) FormatMemory(memoryStr string) string {
	// Формат: "295796/524288"
	return memoryStr + " KB"
}

func (kc *KeeneticClient) FormatMemoryPercent(memoryStr string) string {
	// Парсим "295796/524288"
	parts := make([]string, 2)
	for i, part := range splitMemoryString(memoryStr) {
		if i < 2 {
			parts[i] = part
		}
	}

	if len(parts) == 2 {
		used, err1 := strconv.Atoi(parts[0])
		total, err2 := strconv.Atoi(parts[1])
		if err1 == nil && err2 == nil && total > 0 {
			percent := (used * 100) / total
			return fmt.Sprintf("%d%%", percent)
		}
	}
	return "N/A"
}

func splitMemoryString(s string) []string {
	var parts []string
	current := ""
	for _, char := range s {
		if char == '/' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// Основной метод для получения системного статуса
func (kc *KeeneticClient) GetSystemStatus() string {
	data, err := kc.GetSystemInfo()
	if err != nil {
		return "❌ Ошибка получения статуса системы"
	}

	hostname := getStringSafe(data, "hostname")
	uptime := kc.FormatUptime(getStringSafe(data, "uptime"))
	cpuload := getStringSafe(data, "cpuload")
	memory := kc.FormatMemory(getStringSafe(data, "memory"))
	memoryPercent := kc.FormatMemoryPercent(getStringSafe(data, "memory"))
	connfree := getStringSafe(data, "connfree")
	conntotal := getStringSafe(data, "conntotal")

	return fmt.Sprintf(`🖥️ **Статус системы**

📟 Роутер: %s
⏱️ Аптайм: %s
⚡ Загрузка CPU: %s%%
💾 Память: %s (%s)
🔗 Подключения: %s/%s свободно`,
		hostname, uptime, cpuload, memory, memoryPercent, connfree, conntotal)
}

// Детальный системный статус
func (kc *KeeneticClient) GetDetailedSystemStatus() string {
	data, err := kc.GetSystemInfo()
	if err != nil {
		return "❌ Ошибка получения детального статуса"
	}

	hostname := getStringSafe(data, "hostname")
	domainname := getStringSafe(data, "domainname")
	uptime := kc.FormatUptime(getStringSafe(data, "uptime"))
	cpuload := getStringSafe(data, "cpuload")
	memory := getStringSafe(data, "memory")
	memfree := getStringSafe(data, "memfree")
	membuffers := getStringSafe(data, "membuffers")
	memcache := getStringSafe(data, "memcache")
	connfree := getStringSafe(data, "connfree")
	conntotal := getStringSafe(data, "conntotal")
	swap := getStringSafe(data, "swap")

	return fmt.Sprintf(`🖥️ **Детальный статус системы**

📟 Хостнейм: %s
🏠 Домен: %s
⏱️ Аптайм: %s
⚡ Загрузка CPU: %s%%

💾 **Память:**
• Общая: %s
• Свободно: %s KB
• Буферы: %s KB
• Кэш: %s KB

🔗 **Сетевые подключения:**
• Свободно: %s/%s

💽 **SWAP:** %s`,
		hostname, domainname, uptime, cpuload, memory, memfree, membuffers, memcache, connfree, conntotal, swap)
}

// Комбинированный статус системы и VPN
func (kc *KeeneticClient) GetCombinedStatus() string {
	systemStatus := kc.GetSystemStatus()

	vpnState, err := kc.GetWireGuardState()
	vpnInfo := ""
	if err != nil {
		vpnInfo = "\n🛡️ **VPN:** ❌ Не доступен"
	} else {
		vpnInfo = fmt.Sprintf("\n🛡️ **VPN:** %s", vpnState)
	}

	return systemStatus + vpnInfo
}

// Проверка подключения
func (kc *KeeneticClient) CheckConnection() error {
	_, err := kc.GetSystemInfo()
	return err
}

// Функция для отладки - печатает полную структуру JSON
func (kc *KeeneticClient) DebugCommand(command string) {
	fmt.Printf("🔍 Debug command: %s\n", command)

	data, err := kc.ExecuteRCI(command)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	fmt.Printf("📊 Response structure:\n%s\n", string(jsonData))
}

// Методы для Telegram бота
func (kc *KeeneticClient) HandleStatusCommand() string {
	return kc.GetCombinedStatus()
}

func (kc *KeeneticClient) HandleSystemCommand() string {
	return kc.GetSystemStatus()
}

func (kc *KeeneticClient) HandleDetailCommand() string {
	return kc.GetDetailedSystemStatus()
}

func (kc *KeeneticClient) HandleRestartCommand() string {
	if err := kc.RestartInterface("WireGuard1"); err != nil {
		return "❌ Ошибка перезапуска VPN"
	}

	time.Sleep(2 * time.Second)
	newState, _ := kc.GetWireGuardState()

	return fmt.Sprintf("✅ VPN перезапущен\n🛡️ Новый статус: %s", newState)
}

// Пример использования
func main() {

	user_settings.New()
	cfg := config.NewConfig()

	s := telgram_bot.NewBotService()

	bot := telegram.NewClient(cfg.Telegram.Token, s.HandleCommand, s.HandleMessage, s.HandleCallbackQuery)
	go bot.Start(context.Background())
	client := NewKeeneticClient()

	fmt.Println("🔌 Testing connection to Keenetic...")

	// Проверяем подключение
	if err := client.CheckConnection(); err != nil {
		fmt.Printf("❌ Connection failed: %v\n", err)
		return
	}
	//fmt.Println("✅ Successfully connected to Keenetic!")
	//
	//// Тестируем разные форматы статуса
	//fmt.Println("\n" + client.GetSystemStatus())
	//fmt.Println("\n" + client.GetCombinedStatus())
	//
	//// Детальный статус
	//fmt.Println("\n" + client.GetDetailedSystemStatus())
	//
	//// Проверяем VPN отдельно
	//wgState, err := client.GetWireGuardState()
	//if err != nil {
	//	fmt.Printf("\n⚠️ WireGuard status error: %v\n", err)
	//} else {
	//	fmt.Printf("\n🛡️ WireGuard State: %s\n", wgState)
	//}
}
