package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
)

var pathToTemplates = "./cmd/web/templates"

// TemplateData includes data which will be added when rendering a web page
type TemplateData struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float64
	Data          map[string]any
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
	//User          *data.User
}

// render is used to render a web page based on name of web page and data of it
func (app *Config) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) {
	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", pathToTemplates),
		fmt.Sprintf("%s/header.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/navbar.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/footer.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/alerts.partial.gohtml", pathToTemplates),
	}

	var templateSlide []string
	templateSlide = append(templateSlide, fmt.Sprintf("%s/%s", pathToTemplates, t))

	for _, x := range partials {
		templateSlide = append(templateSlide, x)
	}

	if td == nil {
		td = &TemplateData{}
	}

	tmpl, err := template.ParseFiles(templateSlide...)
	if err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, app.AddDefaultData(td, r)); err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// AddDefaultData to initiate default data for a web page
func (app *Config) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")

	if app.IsAuthenticate(r) {
		td.Authenticated = true
		// TODO - get more user information
	}
	td.Now = time.Now()

	return td
}

// IsAuthenticate return TRUE if the user authenticated before in this session
func (app *Config) IsAuthenticate(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "userID")
}
