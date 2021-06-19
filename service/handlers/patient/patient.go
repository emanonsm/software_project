package patient

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// CreatePatient - создание пациента, url "/create/patient"
func CreatePatient(w http.ResponseWriter, r *http.Request, db *sql.DB)  {

	// Закрытие тела реквеста в конце работы функции
	defer r.Body.Close()

	// Выбор шаблонов, на основе которых будет отрисована страница
	files, err := template.ParseFiles(
		"static/templates/base.tmpl",
		"static/templates/navbar.tmpl",
		"static/templates/patient/create_patient.tmpl",
	)

	// Обработка ошибок
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse templates: %s", err), http.StatusInternalServerError)
		return
	}

	// Данные, на основе которых будут отрисовываться шаблоны
	data := struct{Title string; Message string}{
		Title: "Внести пациента",
		Message: "ФОРМА СОЗДАНИЯ ПАЦИЕНТА",
	}

	// Если нажали на кнопку создать, то идем сюда (http метод POST)
	if r.Method == http.MethodPost {

		// Вычитываем значения из HTML формы
		r.ParseForm()

		// Получаем возраст (inputAge) из формы и преобразуем к инту
		age, err := strconv.Atoi(r.Form.Get("inputAge"))
		if err != nil {
			http.Error(w, fmt.Sprintf("wrong age format: %s", r.Form.Get("inputAge")), http.StatusBadRequest)
			return
		}

		// Получаем рост (inputHeight) из формы и преобразуем к инту
		height, err := strconv.Atoi(r.Form.Get("inputHeight"))
		if err != nil {
			http.Error(w, fmt.Sprintf("wrong height format: %s", r.Form.Get("inputHeight")), http.StatusBadRequest)
			return
		}

		// Получаем вес (inputWeight) из формы и преобразуем к инту
		weight, err := strconv.Atoi(r.Form.Get("inputWeight"))
		if err != nil {
			http.Error(w, fmt.Sprintf("wrong weight format: %s", r.Form.Get("inputWeight")), http.StatusBadRequest)
			return
		}

		// Создаем строку запроса в БД
		query := fmt.Sprintf(
			"INSERT INTO patients(name, surname, middlename, sex, cart, polis, document, snils, age, height, weight) " +
				"VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, %d)",
			r.Form.Get("inputName"),
			r.Form.Get("inputSurname"),
			r.Form.Get("inputMiddlename"),
			r.Form.Get("sex"),
			r.Form.Get("inputCart"),
			r.Form.Get("inputPolis"),
			r.Form.Get("inputDocument"),
			r.Form.Get("inputSnils"),
			age,
			height,
			weight,
		)

		// Выполняем запрос в базе данных, при необходимости обрабатываем ошибку
		data.Message = "ПАЦИЕНТ УСПЕШНО ДОБАВЛЕН!"
		if _, err := db.Exec(query); err != nil {
			data.Message = "ПРОИЗОШЛА ОШИБКА, ПОВТОРИТЕ ЗАПРОС"
		}
	}

	// Отрисовываем шаблон с данными data + обработка ошибок
	if err := files.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("unable to execute templates: %s", err), http.StatusInternalServerError)
		return
	}
}
