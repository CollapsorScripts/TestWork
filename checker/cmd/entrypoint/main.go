package main

import (
	"checker/cmd/config"
	"checker/pkg/logger"
	"checker/pkg/worker"
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	config.Init()
	if err := logger.New(); err != nil {
		log.Fatalf("Ошибка при инициализации логгера: %v", err)
		return
	}

	kw, err := worker.New()
	if err != nil {
		logger.Error("Ошибка при создании воркера: %v", err)
		return
	}

	kw.RunPooling()

	//Для правильного завершения приложения
	{
		wait := time.Second * 15
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		//Блокируем горутину до вызова сигнала Interrupt
		<-c

		_, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		kw.CloseKafka()
		os.Exit(0)
	}
}
