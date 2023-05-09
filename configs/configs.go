package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"sync"
)

type Config struct {
	Telegram struct {
		TgBotHost  string `env:"TELEGRAM_BOT_HOST"`
		TgBotToken string `env:"TELEGRAM_BOT_TOKEN"`
		BatchSize  int    `env:"BATCH_SIZE"`
	}
	SQlite struct {
		StoragePath string `env:"STORAGE_PATH"`
	}
}

var cfg *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading env variables: %s", err.Error())
		}

		log.Println("Read application environment variables")
		cfg = &Config{}
		if err := cleanenv.ReadEnv(cfg); err != nil {
			help, _ := cleanenv.GetDescription(cfg, nil)
			log.Println(help)
			log.Fatalln(err)
		}
	})

	return cfg
}
