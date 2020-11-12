package main

import (
	"aitu.com/snippetbox/pkg/forms"
	"aitu.com/snippetbox/pkg/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl.html", &templateData{
		Snippets: snippet,
	})

}


func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}


	flash := app.session.PopString(r, "flash")
	app.render(w, r, "show.page.tmpl.html", &templateData{
		Flash:flash,
		Snippet: snippet,
	})




}



func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl.html", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//form := forms.New(r.PostForm)
	//form.Required("title", "content", "expires")
	//form.MaxLength("title", 100)
	//form.PermittedValues("expires", "2021-12-15", "2020-11-16", "2020-11-11")

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires","created","profits")
	form.MaxLength("title", 100)

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl.html", &templateData{Form: form})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"),form.Get("created"), form.Get("expires"),form.Get("profits"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Company successfully added!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

