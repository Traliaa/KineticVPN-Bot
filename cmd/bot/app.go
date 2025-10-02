package main

import (
	"context"
	"log"

	"github.com/Traliaa/KineticVPN-Bot/internal/adapter/telegram"
	"github.com/Traliaa/KineticVPN-Bot/internal/app"
	"github.com/Traliaa/KineticVPN-Bot/internal/pg/user_settings"
	"github.com/Traliaa/KineticVPN-Bot/internal/prepare"
	"github.com/Traliaa/KineticVPN-Bot/internal/usecase/telgram_bot"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func mustNewApp() (*app.App, context.Context) {
	ctx := context.Background()
	a := app.NewApp()
	db, pool, err := prepare.MustNewPg(ctx, a.GetConfig())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.Conn()

	user_settings.New()

	s := telgram_bot.NewBotService()

	bot := telegram.NewClient(a.GetConfig().Telegram.Token, s.HandleCommand, s.HandleMessage, s.HandleCallbackQuery)

	a.SetBot(bot)
	//
	//client := NewKeeneticClient()
	//
	//fmt.Println("🔌 Testing connection to Keenetic...")

	//// Проверяем подключение
	//if err := client.CheckConnection(); err != nil {
	//	fmt.Printf("❌ Connection failed: %v\n", err)
	//	//return
	//}
	//fmt.Println("✅ Successfully connected to Keenetic!")

	a.SetRiver(prepare.MustNewRiver(ctx, pool))
	a.Handlers.Handle("/metrics", promhttp.Handler())

	//go func() {
	//	river.Start(ctx)
	//}()
	//http.AddRouter(a)
	////a.Handlers.Handle("/riverui/", river)

	//mux := chi.NewMux()
	//mux.Handle("/riverui/", river)
	//a.Server.Handler = mux

	//http.AddHandler(a, "/riverui/", river)

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
	return a, ctx

}
