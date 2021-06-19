package root

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

// ShowRoot - функция, выполняемая при переходе по урлу "/"
func ShowRoot(w http.ResponseWriter, r *http.Request, _ *sql.DB) {

	// Закрытие тела реквеста в конце работы функции
	defer r.Body.Close()

	// Выбор шаблонов, на основе которых будет отрисована страница
	files, err := template.ParseFiles(
		"static/templates/base.tmpl",
		"static/templates/root/root.tmpl",
	)

	// Обработка ошибок
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse templates: %s", err), http.StatusInternalServerError)
		return
	}

	// Данные, на основе которых будут отрисовываться шаблоны
	data := struct{Title string}{
		Title: "Главная страница",
	}

	// Отрисовка страницы с данными data + обработка ошибки
	if err := files.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("unable to execute templates: %s", err), http.StatusInternalServerError)
		return
	}
}
