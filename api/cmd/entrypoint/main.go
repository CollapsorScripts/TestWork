package main

import (
	"api/cmd/config"
	"api/pkg/database"
	"api/pkg/database/models"
	"api/pkg/gateway"
	"api/pkg/logger"
	_ "api/pkg/utilities"
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"time"
)

func preload() error {
	//Инициализация конфигурации
	config.Init()

	//Инициализация базы данных
	if _, err := database.Init(); err != nil {
		return errors.New("Ошибка при инициализации бд: " + err.Error())
	}

	//Миграции
	err := database.GetDB().AutoMigrate(&models.ModelDB{})
	if err != nil {
		logger.Error("Ошибка при получении всех админов: %v", err)
		return err
	}

	return nil
}

func main() {
	//Инициализация логгера
	{
		if err := logger.New(); err != nil {
			log.Fatalf("Ошибка при инициализации логгера: %v", err)
		}
	}
	
	//Инициализация бд + миграции
	if err := preload(); err != nil {
		logger.Error("%v", err)
		return
	}

	apiRouter, err := gateway.New()
	if err != nil {
		logger.Error("Ошибка при создании экземпляра роутера: %v", err)
		return
	}

	srv := apiRouter.PreloadRoutes()

	{
		wait := time.Second * 15

		// Запуск сервера в отдельном потоке
		go func() {
			logger.Info("Сервер запущен на адресе: %s", srv.Addr)
			if err := srv.ListenAndServe(); err != nil {
				log.Println(err)
			}
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		//Блокируем горутину до вызова сигнала Interrupt
		<-c

		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			logger.Error("%s", err.Error())
		}
		logger.Warn("Выключение сервера")
		os.Exit(0)
	}
}
