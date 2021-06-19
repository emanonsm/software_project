package service

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Service - Структура, которая хранит в себе подключение к БД и роутер
type Service struct {
	router *mux.Router // роутер
	db     *sql.DB    // БД
}

// New - Создание структуры
func New(db *sql.DB) *Service {
	return &Service{
		router: mux.NewRouter(),
		db: db,
	}
}

// InitRoutes - Метод, объявляющий пути и исполняемые по ним функции (хендлеры)
func (s *Service) InitRoutes() *Service {

	// Файловый сервер
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Пути и хендлеры
	s.router.HandleFunc("/", s.ShowRoot) // Главная страница
	s.router.HandleFunc("/create/patient", s.CreatePatient) // Создание пациента
	s.router.HandleFunc("/create/visit", s.CreateVisit) // Создание посещения
	s.router.HandleFunc("/search/{id:[0-9]+}", s.Search) // Поиск
	s.router.HandleFunc("/search/{id:[0-9]+}/visits", s.Visits) // Просмотр посещений
	s.router.HandleFunc("/search/{id:[0-9]+}/visits/{vid:[0-9]+}", s.ShowVisit) // Просмотр конкретного посещения

	return s
}

// ListenAndServe - Старт сервера
func (s *Service) ListenAndServe(addr string) {
	log.Println("Сервер запущен на порту", addr)

	http.ListenAndServe(addr, s.router)
}


