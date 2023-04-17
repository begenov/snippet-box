package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/golangcollege/sessions"

	"github.com/begenov/snippet-box/internal/controller"
	"github.com/begenov/snippet-box/internal/repository"
	"github.com/begenov/snippet-box/internal/service"
	template "github.com/begenov/snippet-box/template"
)

var addr, dsn, secret *string

// var loggers logger.IAuthLoggers

func init() {
	pkg.Removefile()
	addr = flag.String("addr", ":4000", "Сетевой адрес веб-сервера")
	dsn = flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "Название MySQL источника данных")
	secret = flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
}

func main() {
	flag.Parse()
	cache, err := template.NewTemplatesData("./ui/html/")
	if err != nil {
		log.Fatal(err)
	}
	db, err := repository.Open(*dsn)
	if err != nil {
		log.Fatal(err)
		return
	}
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	repo := repository.NewDB(db)
	services := service.NewApplication(repo, &cache, *session)
	handler := controller.NewHandler(services)

	srv := &http.Server{
		Addr:    *addr,
		Handler: handler.Router(),
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

}
