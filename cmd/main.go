package main

import (
	"context"
	"github.com/commedesvlados/go-tg-bot/configs"
	event_consumer "github.com/commedesvlados/go-tg-bot/internal/consumer/event-consumer"
	"github.com/commedesvlados/go-tg-bot/internal/events/telegram"
	"github.com/commedesvlados/go-tg-bot/internal/storage/sqlite"
	sqliteClient "github.com/commedesvlados/go-tg-bot/pkg/clients/sqlite"
	telegramClient "github.com/commedesvlados/go-tg-bot/pkg/clients/telegram"
	"log"
)

func main() {
	// config
	cfg := configs.GetConfig()

	// db
	db, err := sqliteClient.NewSqliteClient(cfg.SQlite.StoragePath)
	if err != nil {
		log.Fatalln(err)
	}

	storage := sqlite.NewStorage(db)
	if err := storage.Init(context.Background()); err != nil {
		log.Fatalln(err)
	}

	log.Println("[STORAGE] sqlite is connected and initialized")

	// app
	tgClient := telegramClient.NewClient(cfg.Telegram.TgBotHost, cfg.Telegram.TgBotToken)

	eventProcessor := telegram.NewProcessor(tgClient, storage)

	log.Println("[START] Service started")

	consumer := event_consumer.NewConsumer(eventProcessor, eventProcessor, cfg.Telegram.BatchSize)

	if err := consumer.Start(context.Background()); err != nil {
		log.Fatalln("[STOP] Service is stopped", err)
	}

}
