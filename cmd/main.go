package main

import (
	"github/commedesvlados/go-tg-bot/configs"
	event_consumer "github/commedesvlados/go-tg-bot/internal/consumer/event-consumer"
	"github/commedesvlados/go-tg-bot/internal/events/telegram"
	"github/commedesvlados/go-tg-bot/internal/storage/files"
	telegramClient "github/commedesvlados/go-tg-bot/pkg/clients/telegram"
	"log"
)

const (
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	cfg := configs.GetConfig()

	tgClient := telegramClient.NewClient(cfg.Telegram.TgBotHost, cfg.Telegram.TgBotToken)

	eventProcessor := telegram.NewProcessor(tgClient, files.NewStorage(storagePath))

	log.Println("[START] Service started")

	consumer := event_consumer.NewConsumer(eventProcessor, eventProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatalln("[STOP] Service is stopped", err)
	}

}
