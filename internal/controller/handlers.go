package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/begenov/snippet-box/internal/e"

	templates "github.com/begenov/snippet-box/template"

	"github.com/begenov/snippet-box/internal/service/forms"
)

func (h *Handler) home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			log.Printf("%v\n", e.Wrap(homeError, errors.New("status not found")))
			h.services.NotFound(w)
			return
		}

		s, err := h.services.Latest()
		if err != nil {
			h.services.ServerError(w, err)
			return
		}
		h.services.Render(w, r, "home.page.html", &templates.TemplateData{
			Snippets: s,
		})

	}
}

func (h *Handler) createSnippetForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("no")

		h.services.Render(w, r, "create.page.html", &templates.TemplateData{
			Form: forms.New(nil),
		})
	}
}

func (h *Handler) createSnippet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("%v\n", e.Wrap(createSnippetError, errors.New("method not allowed")))
			h.services.ClientError(w, http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			h.services.ClientError(w, http.StatusBadRequest)
			return
		}
		form := forms.New(r.PostForm)
		form.Required("title", "content", "expires")
		form.MaxLength("title", 100)
		form.PermittedValues("expires", "365", "7", "1")
		if !form.Valid() {
			h.services.Render(w, r, "create.page.html", &templates.TemplateData{
				Form: form,
			})
			return
		}

		id, err := h.services.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
		if err != nil {
			h.services.ServerError(w, err)
			return
		}
		h.services.Put(r, "flash", "Snippet successfully created!")
		http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
	}
}

func (h *Handler) showSnippet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get(":id"))
		fmt.Printf("r.URL.Query().Get(\"id\"): %v\n", r.URL.Query().Get(":id"))
		if err != nil || id < 1 {
			fmt.Println(err)
			log.Printf("%v\n", e.Wrap(chowSnippetError, errors.New("not found")))
			h.services.NotFound(w)
			return
		}
		s, err := h.services.Get(id)
		if err != nil {
			if errors.Is(err, ErrNoRecord) {
				h.services.NotFound(w)
			} else {
				h.services.ServerError(w, err)
			}
			return
		}

		h.services.Render(w, r, "show.page.html", &templates.TemplateData{
			Snippet: s,
		})
	}
}

func (h *Handler) signupUserForm(w http.ResponseWriter, r *http.Request) {
	h.services.Render(w, r, "singup.page.html", &templates.TemplateData{
		Form: forms.New(nil),
	})
}

func (h *Handler) singupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.services.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MinLength("password", 10)
	form.MatchesPattern("email", forms.EmailRx)
	if !form.Valid() {
		h.services.Render(w, r, "singup.page.html", &templates.TemplateData{
			Form: form,
		})
	}
	err = h.services.InsertUser(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, fmt.Errorf("error model duplicate email")) {
			form.Errors.Add("email", "ADdress is already in use")
			h.services.Render(w, r, "singup.page.html", &templates.TemplateData{
				Form: form,
			})
		} else {
			h.services.ServerError(w, err)
		}
		return
	}
	h.services.Put(r, "flash", "Your singup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (h *Handler) loginUserForm(w http.ResponseWriter, r *http.Request) {
	h.services.Render(w, r, "login.page.html", &templates.TemplateData{
		Form: forms.New(nil),
	})
}

var err1 = fmt.Errorf("error invalid creations")

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.services.ClientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	id, err := h.services.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if !errors.Is(err, err1) {
			form.Errors.Add("generic", "Email or Password is incorect")
			h.services.Render(w, r, "login.page.html", &templates.TemplateData{
				Form: form,
			})
		} else {
			h.services.ServerError(w, err)
		}
		return
	}
	h.services.Put(r, "authenticatedUserID", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)

}

func (h *Handler) logoutUser(w http.ResponseWriter, r *http.Request) {
	h.services.Remove(r, "authenticatedUserID")
	h.services.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
