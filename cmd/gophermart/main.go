package main

import (
	"flag"
	"github.com/xbreathoflife/gophermart/config"
	"github.com/xbreathoflife/gophermart/internal/app/server"
	"github.com/xbreathoflife/gophermart/internal/app/storage"
	"log"
	"net/http"
)

func parseFlags(conf *config.Config) {
	address := flag.String("a", "", "Адрес запуска HTTP-сервера")
	connString := flag.String("d", "", "Строка с адресом подключения к БД")
	serviceAddress := flag.String("r", "", "Адрес системы расчёта начислений: переменная окружения ОС")
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

	gophermartServer := server.NewGothServer(dbStorage, conf.ServiceAddress)
	r := gophermartServer.ServerHandler()

	log.Fatal(http.ListenAndServe(conf.Address, r))
}
