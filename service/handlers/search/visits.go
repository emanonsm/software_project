package search

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

// Visit - структура посещения
type Visit struct {
	VisitID   int    // ID визита
	History   string // анамнез
	Review    string // обследование
	Diag      string // основной диагноз
	Recommend string // рекомендации
	VisitTS   string // время визита
	Doctor    string // фио лечащего врача
	Condition string // состояние
	Services  string // услуги
	Medicine  string // лекарства
	Recipe    string // рецепт
}

// Visits - функция, исполняемая при переходе на url "/search/{id:[0-9]+}/visits"
func Visits(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Закрытие тела реквеста в конце работы функции
	defer r.Body.Close()

	// Получение id из урла, преобразование к инту
	patientID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(patientID)
	if err != nil {
		http.Error(w, fmt.Sprintf("bad patientID request: %s", patientID), http.StatusBadRequest)
		return
	}

	// Выбор шаблонов, на основе которых отобразится страница
	files, err := template.ParseFiles(
		"static/templates/base.tmpl",
		"static/templates/navbar.tmpl",
		"static/templates/search/visits.tmpl",
	)

	// Обработка ошибок
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse templates: %s", err), http.StatusInternalServerError)
		return
	}

	// Данные, на основе которых будут отрисовываться шаблоны
	data := struct{Title string; Message string; PatientID int; Visits []Visit}{
		Title: "Посещения пациента",
		PatientID: id,
		Message: "ПОСЕЩЕНИЯ ПАЦИЕНТА",
	}

	// Создание строки запроса в базу и выполнение данного запроса, после чего обработка возможных ошибок
	rows, err := db.Query(fmt.Sprintf("SELECT vid, visit_ts FROM visits WHERE patient_id = %d", data.PatientID))
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to find visits: %s", err), http.StatusInternalServerError)
		return
	}

	// Создаем массив посещений. Если вернулся ответ из базы с ненулевым количеством строк, то...
	visits := make([]Visit, 0)
	if rows != nil {

		// Проходимся по каждой строке, вычитываем значения из строк в структуру "visit"
		for rows.Next() {
			visit := Visit{}
			var ts time.Time
			if err := rows.Scan(
				&visit.VisitID,
				&ts,
			); err != nil {
				data.Message = "ПОСЕЩЕНИЯ НЕ НАЙДЕНЫ"
				break
			}

			// Преобразовываем время к нужному формату и запихиваем структуру "visit" в список таких визитов "visits"
			visit.VisitTS = ts.Format("02-01-2006 15:04:05")
			visits = append(visits, visit)
		}

		// Записываем визиты в данные, на основе которых будут отрисовываться шаблоны
		data.Visits = visits
	}

	// Создание шаблонов и проброс в них переменной с данными data
	if err := files.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("unable to execute templates: %s", err), http.StatusInternalServerError)
		return
	}
}

// ShowVisit - функция, исполняемая при переходе на url "/search/{id:[0-9]+}/visits/{vid:[0-9]+}"
func ShowVisit(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Закрытие тела реквеста в конце работы функции
	defer r.Body.Close()

	// Получение vid из урла, преобразование к инту
	visitID := mux.Vars(r)["vid"]
	id, err := strconv.Atoi(visitID)
	if err != nil {
		http.Error(w, fmt.Sprintf("bad visitID request: %s", visitID), http.StatusBadRequest)
		return
	}

	// Получение id из урла. преобразование к инту
	patientID := mux.Vars(r)["id"]
	pid, err := strconv.Atoi(patientID)
	if err != nil {
		http.Error(w, fmt.Sprintf("bad patientID request: %s", patientID), http.StatusBadRequest)
		return
	}

	// Выбор шаблонов, на основе которых будет отрисована страница
	files, err := template.ParseFiles(
		"static/templates/base.tmpl",
		"static/templates/navbar.tmpl",
		"static/templates/search/visit.tmpl",
	)

	// Обработка ошибок
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse templates: %s", err), http.StatusInternalServerError)
		return
	}

	// Данные, на основе которых будут отрисовываться шаблоны
	data := struct{Title string; Message string; PatientID int; Visit Visit}{
		Title: "Информация о приеме",
		Message: "ДАННЫЕ О ПРИЕМЕ",
		PatientID: pid,
	}

	// Создание строки запроса в базу и выполнение данного запроса
	row := db.QueryRow(fmt.Sprintf("SELECT " +
		"vid, history, review, diag, recommend, visit_ts, doctor, condition, services, medicine, recipe " +
		"FROM visits WHERE vid = %d", id))

	// Если вернулась строка с данными, то считываем данные в структуру
	if row != nil {
		visit := Visit{}
		var ts time.Time
		if err := row.Scan(
			&visit.VisitID,
			&visit.History,
			&visit.Review,
			&visit.Diag,
			&visit.Recommend,
			&ts,
			&visit.Doctor,
			&visit.Condition,
			&visit.Services,
			&visit.Medicine,
			&visit.Recipe,
		); err != nil { // Если ошибка, то отрисовываем страницу с ошибкой
			data.Message = "ПРИЕМ НЕ НАЙДЕН"
			if err := files.Execute(w, data); err != nil {
				http.Error(w, fmt.Sprintf("unable to execute templates: %s", err), http.StatusInternalServerError)
				return
			}

			return
		}

		// Если ошибки не было, преобразовываем время к нужному формату и проставляем в data структуру visit с данными из базы
		visit.VisitTS = ts.Format("02-01-2006 15:04:05")
		data.Visit = visit
	}

	// Отрисовываем шаблон с данными data + обработка ошибок
	if err := files.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("unable to execute templates: %s", err), http.StatusInternalServerError)
		return
	}
}
