package application

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/sgoldenf/wb_l0/internal/model"
	"github.com/sgoldenf/wb_l0/internal/templates"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		orderID := r.FormValue("order_id")
		http.Redirect(w, r, "/order/?order_id="+orderID, http.StatusSeeOther)
		return
	}
	app.render(w, http.StatusOK, &templates.TemplateData{})
}

func (app *Application) viewOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	data, ok := app.cache.Load(orderID)
	if !ok {
		order, err := app.Orders.ReadOrder(orderID)
		if err != nil {
			if errors.Is(err, model.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}
		data = order.Data
		app.cache.Store(orderID, data)
	}
	app.render(w, http.StatusOK,
		&templates.TemplateData{
			OrderData: data.(string),
		},
	)
}

func (app *Application) render(
	w http.ResponseWriter, status int, data *templates.TemplateData) {
	buf := new(bytes.Buffer)
	err := app.TemplateCache.ExecuteTemplate(buf, "index", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *Application) serverError(w http.ResponseWriter, err error) {
	entry := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, entry)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
