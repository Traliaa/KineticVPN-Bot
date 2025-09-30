# KineticVPN-Bot
Автоматическое обновление VPN на вашем роутере Кинетик через Telegram. Бот сам следит за соединением и меняет VPN при обрыве. Смените IP-адрес одной командой.



https://help.keenetic.com/hc/ru/articles/11282223272092-%D0%9F%D1%80%D0%B8%D0%BC%D0%B5%D0%BD%D0%B5%D0%BD%D0%B8%D0%B5-%D0%BC%D0%B5%D1%82%D0%BE%D0%B4%D0%BE%D0%B2-API-%D0%BF%D0%BE%D1%81%D1%80%D0%B5%D0%B4%D1%81%D1%82%D0%B2%D0%BE%D0%BC-%D1%81%D0%B5%D1%80%D0%B2%D0%B8%D1%81%D0%B0-HTTP-Proxy


//Работает
// curl -u 'api:Demon0203' http://192.168.2.1:81/rci/show/system
//curl -u 'api:Demon0203' https://rci.tankhome.netcraze.pro/rci/show/system
//curl -u 'api:Demon0203' http://rci.tankhome.netcraze.pro/rci/show/system


	ctx := context.Background()
	b := telegram.NewClient("8250747795:AAFBy_jRtBWmeJkMDCGnLr4LOZjgfZ4dFB0")
	go b.Start(ctx)