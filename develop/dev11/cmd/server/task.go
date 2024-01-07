package main

import (
	"WB_LVL2/develop/dev11/internal/api"
	"WB_LVL2/develop/dev11/internal/storage"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем.
В рамках задания необходимо работать строго со
стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для
	сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для
	парсинга и валидации параметров методов
	/create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого
	из методов API, используя вспомогательные
	функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API:
POST /create_event
POST /update_event
POST /delete_event
GET /events_for_day
GET /events_for_week
GET /events_for_month
Параметры передаются в виде www-url-form-encoded
(т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString,
в POST через тело запроса.
В результате каждого запроса должен возвращаться
JSON документ содержащий либо {"result": "..."}
в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен
	возвращать HTTP 503. В случае ошибки входных данных
	(невалидный int например) сервер должен возвращать HTTP 400.
	В случае остальных ошибок сервер должен возвращать HTTP 500.
	Web-сервер должен запускаться на порту указанном в конфиге
	и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

const cfgPath = "./config.json"

type config struct {
	Port string `json:"port"`
}

func main() {
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка чтения файла конфигурации: %v\n", err)
		os.Exit(1)
	}

	var cfg config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка расшифровки json: %v\n", err)
		os.Exit(1)
	}
	s := storage.NewStorage()
	s.AddUser(1)
	h := api.NewHandler(s)

	server := new(api.Server)

	go func() {
		err := server.Run(cfg.Port, h.InitRoutes())
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка запуска сервера: %v\n", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	err = server.Shutdown(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка остановки сервера: %v\n", err)
	}
	log.Println("сервер успешно остановлен")
}
