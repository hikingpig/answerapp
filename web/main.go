package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func main() {
	tmpl, err := template.ParseGlob("web/templates/*.tmpl.html")
	if err != nil {
		// panic(err)
		http.Handle("/", errHandler(err.Error(), http.StatusInternalServerError))
		return
	}
	http.Handle("/questions/", templateHandler(tmpl, "question.tmpl.html"))
	http.Handle("/", templateHandler(tmpl, "index.tmpl.html"))
	appengine.Main()
}

func templateHandler(tmpl *template.Template, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		if ok := ensureUser(ctx, w, r); !ok {
			return
		}
		err := tmpl.ExecuteTemplate(w, name, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			// panic(err)
		}
	})
}

func errHandler(err string, code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, err, code)
	})
}

func ensureUser(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
	me := user.Current(ctx)
	log.Println("===== user", me)
	if me == nil {
		loginURL, err := user.LoginURL(ctx, r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
		return false
	}
	return true
}
