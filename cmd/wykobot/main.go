package main

import (
	"flag"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xaxes/vikop-gorace/internal"
	"github.com/xaxes/vikop-gorace/internal/log"
	"github.com/xaxes/vikop-gorace/internal/tg"
	"github.com/xaxes/vikop-gorace/pkg/dedup"
	"github.com/xaxes/vikop-gorace/pkg/wykop"
	"go.uber.org/zap"
)

func newTgBot(l *zap.Logger, cfg internal.TelegramConfig) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}
	l.Debug("connect to telegram", zap.String("username", bot.Self.UserName))

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	return bot, nil
}

func newCfg(l *zap.Logger, cfgPath string, dbPath *string) (internal.Config, error) {
	cfg, err := internal.LoadConfig(cfgPath)
	if err != nil {
		return internal.Config{}, err
	}
	if dbPath != nil {
		cfg.BadgerDirpath = *dbPath
	}

	return cfg, nil
}

func fetchAndSend(l *zap.Logger, deduper dedup.Deduper, cfg internal.Config, bot *tgbotapi.BotAPI) {
	w := wykop.New(l, cfg.Wykop.WykopUserKey, deduper)
	var entries []wykop.Entry

	for _, period := range []wykop.Period{wykop.Six, wykop.Twelve, wykop.TwentyFour} {
		l.Info("get entries", zap.Int("period", int(period)))
		es, err := w.Hot(period, cfg.Wykop.MaxNoOfEntries, cfg.Wykop.MaxPage, cfg.Wykop.TagBlacklist, cfg.Wykop.MinVotes)
		if err != nil {
			l.Error("get entries", zap.Int("period", int(period)), zap.Error(err))
		}
		entries = append(entries, es...)
	}

	for _, entry := range entries {
		msg := tgbotapi.NewMessage(cfg.Telegram.Channel, tg.FormatEntry(entry))
		msg.ParseMode = "MarkdownV2"
		if _, err := bot.Send(msg); err != nil {
			l.Error("send message", zap.Error(err))
			continue
		}
		time.Sleep(cfg.Telegram.MsgDelay)
	}

	l.Info("done", zap.Int("msg_count", len(entries)))
}

func main() {
	cfgPath := flag.String("config", "./config.yaml", "Config path")
	dbPath := flag.String("db", "/tmp/", "Persistent storage for service's state")
	mode := flag.String("mode", "dev", "Logging mode. Options: prod, dev")
	debug := flag.Bool("debug", false, "Log debug level")
	flag.Parse()
	l := log.NewZap(*mode, *debug)

	cfg, err := newCfg(l, *cfgPath, dbPath)
	if err != nil {
		l.Fatal("load config", zap.String("path", *cfgPath), zap.Error(err))
	}
	l.Debug("load config", zap.String("config", cfg.String()))

	bot, err := newTgBot(l, cfg.Telegram)
	if err != nil {
		l.Fatal("init telegram bot", zap.Error(err))
	}

	deduper, err := dedup.NewBadger(cfg.BadgerDirpath)
	if err != nil {
		l.Fatal("init badger deduper", zap.Error(err))
	}

	fetchAndSend(l, deduper, cfg, bot)
}
