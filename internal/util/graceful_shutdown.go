package util

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulShutdownServer(server *http.Server) context.CancelFunc {
	// Создаем контекст, который будет отменен при получении сигнала
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := signals()

	// Горутина для обработки сигналов
	go func() {
		defer cancel()
		<-sigChan
		fmt.Println("\nЗавершение работы...")
		server.Shutdown(ctx)
		time.Sleep(2 * time.Second)
	}()
	return cancel
}

func GracefulShutdown() (context.Context, context.CancelFunc) {
	
	ctx, cancel := context.WithCancel(context.Background())

	// Канал для получения сигналов
	sigChan := signals()

	// Горутина для обработки сигналов
	go func() {
		defer cancel()
		<-sigChan
		fmt.Println("\nЗавершение работы...")
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
	return ctx, cancel
}

func signals() chan os.Signal {
	// Канал для получения сигналов
	sigChan := make(chan os.Signal, 1)
	// Регистрируем сигналы, которые хотим перехватывать
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	return sigChan
}