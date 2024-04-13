package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
)

type TemplateData struct {
	URL             string
	IsAuthenticated bool
	AuthUser        string
	Flash           string
	Error           string
	CSRFToken       string
}

func (a *application) defaultData(td *TemplateData, r *http.Request) *TemplateData {

	td.URL = a.server.url

	return td
}

func (a *application) render(w http.ResponseWriter, r *http.Request, view string, vars jet.VarMap) error {

	td := &TemplateData{}

	td = a.defaultData(td, r)

	//init session

	a.session = scs.New()
	a.session.Lifetime = 24 * time.Hour
	a.session.Cookie.Persist = true
	a.session.Cookie.Domain = a.server.url
	a.session.Cookie.SameSite = http.SameSiteStrictMode

	tp, err := a.view.GetTemplate(fmt.Sprintf("%s.html", view))
	if err != nil {
		return err
	}

	if err = tp.Execute(w, vars, td); err != nil {
		return err
	}
	return nil

}
