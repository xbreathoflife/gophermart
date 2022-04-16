package main

import (
	"context"
	"flag"
	"github.com/xbreathoflife/gophermart/config"
	"github.com/xbreathoflife/gophermart/internal/app/storage"
	"log"
)

func parseFlags(conf *config.Config) {
	address := flag.String("a", "", "Адрес запуска HTTP-сервера")
	connString := flag.String("d", "", "Строка с адресом подключения к БД")
	serviceAddress := flag.String("r", "", "Aдрес системы расчёта начислений: переменная окружения ОС")
	flag.Parse()

	if *address != "" {
		conf.Address = *address
	}

	if *connString != "" {
		conf.ConnString = *connString
	}

	if *serviceAddress != "" {
		conf.ServiceAddress = *serviceAddress
	}
}

func main() {
	conf := config.Init()
	parseFlags(&conf)

	dbStorage := storage.NewDBStorage(conf.ConnString)
	err := dbStorage.Init(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
