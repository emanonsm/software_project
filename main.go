package main

import (
	"database/sql"
	"fmt"
	"github.com/emanonsm/software_project/service"
	_ "github.com/lib/pq"
	"log"
)

// Для подключения к базе
const (
	DRIVER     = "postgres"
	DbHost     = "localhost"
	DbPort     = 5432
	DbUser     = "postgres"
	DbPassword = "postgres"
	DbName     = "testdb"
)

// Начало программы
func main() {

	// Подключение к postgres
	db, _ := sql.Open(
		DRIVER,
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DbHost, DbPort, DbUser, DbPassword, DbName),
	)

	// Обработка неудачного подключения к БД
	if err := db.Ping(); err != nil {
		log.Fatalf("db connection error: %s", err)
	}

	// Создание сервиса -> объявление путей -> поднятие сервера на порту 7000
	service.New(db).InitRoutes().ListenAndServe(":7000")
}
