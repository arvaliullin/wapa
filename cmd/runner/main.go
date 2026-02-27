package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/arvaliullin/wapa/internal/runner/app"
	"github.com/arvaliullin/wapa/internal/utils"
	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	configPath := flag.String("config", "config.yaml", "Путь к конфигурации сервиса")
	flag.Parse()

	config, err := app.NewRunnerConfig(*configPath)
	utils.FatalOnError(err, "Ошибка загрузки конфигурации сервиса: %v", err)

	service, err := app.NewRunnerService(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка инициализации сервиса: %v\n", err)
		os.Exit(1)
	}

	service.Run(ctx)
}
