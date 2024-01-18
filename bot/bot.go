package bot

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Chestnuts4/citrix-update-monitor/config"
	"github.com/prometheus/common/log"
)

// bot 数组
var Bots []Bot

type Bot interface {
	// 发送消息
	SendMsg(msg *config.Msg) error
	// bot名字
	GetBotName() string
	Start(ctx context.Context) error
}

func SendMsgAllBots(m *config.Msg) {
	for _, bot := range Bots {
		bot.SendMsg(m)
	}
}

func Start() {
	tgToken := config.Config.Get("tgbot.token").(string)
	tgProxyStr := config.Config.Get("tgbot.proxy").(string)

	tgbot, err := NewTgbot(tgToken, tgProxyStr)
	if err != nil {
		log.Fatalf("NewTgbot error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lanxinSecret := config.Config.Get("lanxin.secret").(string)
	lanxinWebHook := config.Config.Get("lanxin.webhook").(string)
	lanxinProxy := config.Config.Get("lanxin.proxy").(string)
	lanxinBot, err := NewLangxinBot(lanxinSecret, lanxinWebHook, lanxinProxy)
	if err != nil {
		log.Fatalf("NewTgbot error: %v", err)
	}

	go tgbot.Start(ctx)
	go lanxinBot.Start(ctx)
	Bots = append(Bots, tgbot)
	Bots = append(Bots, lanxinBot)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

}
