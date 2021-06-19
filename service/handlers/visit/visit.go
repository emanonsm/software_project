package visit

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// CreateVisit - функция, исполняемая при переходе на url "/create/visit"
func CreateVisit(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Закрытие тела реквеста в конце работы функции
	defer r.Body.Close()

	// Выбор шаблонов, на основе которых отобразится страница
	files, err := template.ParseFiles(
		"static/templates/base.tmpl",
		"static/templates/navbar.tmpl",
		"static/templates/visit/create_visit.tmpl",
	)

	// Обработка ошибок
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse templates: %s", err), http.StatusInternalServerError)
		return
	}

	// Данные, на основе которых будут отрисовываться шаблоны
	data := struct{Title string; Message string}{
		Title: "Форма прием",
		Message: "ФОРМА ПРИЕМА",
	}

	// Если нажали создать, то идем сюда (http метод POST)
	if r.Method == http.MethodPost {

		// Вычитываем значения из формы
		r.ParseForm()

		// Получаем inputID из формы
		id, err := strconv.Atoi(r.Form.Get("inputID"))
		if err != nil {
			http.Error(w, fmt.Sprintf("wrong id format: %s", r.Form.Get("inputID")), http.StatusBadRequest)
			return
		}

		// Создаем строку запроса в БД, вставляем в неё значения данных из формы
		query := fmt.Sprintf(
			"INSERT INTO visits(patient_id, history, review, diag, recommend, doctor, condition, services, medicine, recipe) " +
				"VALUES(%d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')",
			id,
			r.Form.Get("inputHistory"),
			r.Form.Get("inputReview"),
			r.Form.Get("inputDiag"),
			r.Form.Get("inputRecommend"),
			r.Form.Get("inputDoctor"),
			r.Form.Get("inputCondition"),
			r.Form.Get("inputServices"),
			r.Form.Get("inputMedicine"),
			r.Form.Get("inputRecipe"),
		)

		// Выполнение запроса в БД и обработка ошибки
		data.Message = "ПРИЕМ УСПЕШНО ЗАВЕРШЕН!"
		if _, err := db.Exec(query); err != nil {
			data.Message = "ОШИБКА - НЕКОРРЕКТНЫЕ ДАННЫЕ!"
		}
	}

	// Создание шаблонов и проброс в них переменной с данными data
	if err := files.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("unable to execute templates: %s", err), http.StatusInternalServerError)
		return
	}
}
