package main

import (
	"flag"

	"github.com/arvaliullin/wapa/internal/composer/app"
	"github.com/arvaliullin/wapa/internal/utils"

	_ "github.com/lib/pq"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Путь к конфигурации сервиса")
	flag.Parse()

	config, err := app.NewComposerConfig(*configPath)
	utils.FatalOnError(err, "Ошибка загрузки конфигурации сервиса: %v", err)

	service := app.NewComposerService(config)
	service.Run()
}
