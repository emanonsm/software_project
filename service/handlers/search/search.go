package search

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
)

// Search - вызывается при отображении страницы пациента
func Search(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Закрытие тела реквеста в конце работы функции
	defer r.Body.Close()

	// Получение id из урла и преобразование его к инту
	patientID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(patientID)
	if err != nil {
		http.Error(w, fmt.Sprintf("bad patientID request: %s", patientID), http.StatusBadRequest)
		return
	}

	// Выбор шаблонов, на основе которых будет отрисована страница
	files, err := template.ParseFiles(
		"static/templates/base.tmpl",
		"static/templates/navbar.tmpl",
		"static/templates/search/search.tmpl",
	)

	// Обработка ошибок
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse templates: %s", err), http.StatusInternalServerError)
		return
	}

	// Данные, на основе которых будут отрисовываться шаблоны
	data := struct{
		Title string
		Message string
		ID int
		Name string
		Surname string
		Middlename string
		Age int
		Height int
		Weight int
		Sex string
		Cart string
		Polis string
		Document string
		Snils string
	}{
		Title: "Поиск пациента",
		Message: "ДАННЫЕ ПАЦИЕНТА",
	}

	// Составляем запрос в базу и выполняем его. Если вернулась строка с нормальными данными, то вычитываем её значения в параметры data
	row := db.QueryRow(fmt.Sprintf("SELECT " +
		"pid, name, surname, middlename, age, height, weight, sex, cart, polis, document, snils " +
		"FROM patients WHERE pid = %d", id))
	if row != nil {
		if err := row.Scan(
			&data.ID,
			&data.Name,
			&data.Surname,
			&data.Middlename,
			&data.Age,
			&data.Height,
			&data.Weight,
			&data.Sex,
			&data.Cart,
			&data.Polis,
			&data.Document,
			&data.Snils,
		); err != nil {
			data.Message = "ПАЦИЕНТ НЕ НАЙДЕН"
			if err := files.Execute(w, data); err != nil {
				http.Error(w, fmt.Sprintf("unable to execute templates: %s", err), http.StatusInternalServerError)
				return
			}

			return
		}

		if data.Sex == "m" {
			data.Sex = "Мужской"
		} else {
			data.Sex = "Женский"
		}
	}

	// Отрисовываем шаблон с данными data + обработка ошибок
	if err := files.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("unable to execute templates: %s", err), http.StatusInternalServerError)
		return
	}
}
