package service

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"snippet/internal/service/forms"
	"snippet/model"
	templates "snippet/template"
	"time"
)

type AuthApplication interface {
	ServerError(w http.ResponseWriter, err error)
	ClientError(w http.ResponseWriter, status int)
	NotFound(w http.ResponseWriter)
	Render(w http.ResponseWriter, r *http.Request, name string, td *templates.TemplateData)
	Latest() ([]*model.Snippet, error)
	Get(id int) (*model.Snippet, error)
	Insert(title, content, expires string) (int, error)
	Enable(next http.Handler) http.Handler
	Put(r *http.Request, key string, val interface{})
	PopString(r *http.Request, s string) string
	InsertUser(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Remove(r *http.Request, s string)
	IsAutenticated(r *http.Request) bool
}

func (app *Application) PopString(r *http.Request, s string) string {
	return app.Session.PopString(r, s)
}

func (app *Application) Enable(next http.Handler) http.Handler {
	return app.Session.Enable(next)
}

func (app *Application) Put(r *http.Request, key string, val interface{}) {
	app.Session.Put(r, key, val)
}

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

func (app *Application) addDefaultDate(td *templates.TemplateData, r *http.Request) *templates.TemplateData {
	if td == nil {
		td = &templates.TemplateData{
			Form: forms.New(nil),
		}
	}
	td.CurrentYear = time.Now().Year()
	td.Flash = app.PopString(r, "flash")
	td.IsAuthenticated = app.IsAutenticated(r)
	return td
}

func (app *Application) Render(w http.ResponseWriter, r *http.Request, name string, td *templates.TemplateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.ServerError(w, fmt.Errorf("шаблон %s не существует", name))
		return
	}

	err := ts.Execute(w, app.addDefaultDate(td, r))
	if err != nil {
		app.ServerError(w, err)
	}

}

func (app *Application) Latest() ([]*model.Snippet, error) {
	return app.repo.Latest()
}

func (app *Application) Get(id int) (*model.Snippet, error) {
	return app.repo.Get(id)
}

func (app *Application) Insert(title, content, expires string) (int, error) {
	return app.repo.Insert(title, content, expires)
}

func (app *Application) InsertUser(name, email, password string) error {
	return app.repo.InsertUser(name, email, password)
}

func (app *Application) Remove(r *http.Request, key string) {
	app.Session.Remove(r, key)

}

func (app *Application) Authenticate(email, password string) (int, error) {
	id, err := app.repo.Authenticate(email, password)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return 0, err
	}
	return id, err
}

func (app *Application) IsAutenticated(r *http.Request) bool {
	return app.Session.Exists(r, "authenticatedUserID")
}
