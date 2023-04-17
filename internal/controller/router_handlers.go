package controller

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (h *Handler) Router() http.Handler {
	standardMiddleware := alice.New(h.recoverPanic, logRequest, secureHeaders)
	dynamicMiddleware := alice.New(h.services.Enable)
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(h.home()))
	mux.Get("/snippet/create", dynamicMiddleware.Append(h.requireAuthentication).ThenFunc(h.createSnippetForm()))
	mux.Post("/snippet/create", dynamicMiddleware.Append(h.requireAuthentication).ThenFunc(h.createSnippet()))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(h.showSnippet()))

	// new handlers
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(h.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(h.singupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(h.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(h.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(h.requireAuthentication).ThenFunc(h.logoutUser))
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}

/*
func (h *Handler) Router() http.Handler {
	standardMiddleware := alice.New(h.recoverPanic, logRequest, secureHeaders)
	mux := http.NewServeMux()
	mux.HandleFunc("/", (h.home()))
	mux.HandleFunc("/snippet", h.showSnippet())
	mux.HandleFunc("/snippet/create", h.createSnippetForm())
	mux.HandleFunc("/snippet/create/", h.createSnippet())
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
*/
