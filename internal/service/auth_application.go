package service

import (
	"html/template"
	"log"
	"os"
	"snippet/internal/repository"

	"github.com/golangcollege/sessions"
)

type Application struct {
	repo          repository.IAuthSnippetModel
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	templateCache map[string]*template.Template
	Session       *sessions.Session
}

func NewApplication(repo repository.IAuthSnippetModel, templateCache *map[string]*template.Template, session sessions.Session) *Application {
	return &Application{
		repo:          repo,
		ErrorLog:      log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		InfoLog:       log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		templateCache: *templateCache,
		Session:       &session,
	}
}
