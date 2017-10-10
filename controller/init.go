package controller

import (
	"path"
	"time"
	"net/http"
	"html/template"
	"encoding/json"

	"maintain/config"
)

var templateFiles = []string{
	path.Join(config.WebRoot, "views", "index.html"),
	path.Join(config.WebRoot, "views", "login.html"),
	path.Join(config.WebRoot, "views", "projects.html"),
	path.Join(config.WebRoot, "views", "project_add.html"),
	path.Join(config.WebRoot, "views", "project_edit.html"),
	path.Join(config.WebRoot, "views", "deploys.html"),
	path.Join(config.WebRoot, "views", "minions.html"),
	path.Join(config.WebRoot, "views", "minion_add.html"),
	path.Join(config.WebRoot, "views", "minion_edit.html"),
	path.Join(config.WebRoot, "views", "clusters.html"),
	path.Join(config.WebRoot, "views", "cluster_add.html"),
	path.Join(config.WebRoot, "views", "cluster_edit.html"),
	path.Join(config.WebRoot, "views", "widgets", "header.html"),
	path.Join(config.WebRoot, "views", "widgets", "footer.html"),
	path.Join(config.WebRoot, "views", "widgets", "scripts.html"),
	path.Join(config.WebRoot, "views", "widgets", "topbar.html"),
	path.Join(config.WebRoot, "views", "widgets", "left_menu.html"),
}

var funcMap = template.FuncMap{
	"formatDate": func(t time.Time, format string) string { return t.Format(format) },
}

var templates = template.Must(template.New("default").Funcs(funcMap).ParseFiles(templateFiles...))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type JsonRet struct {
	Errcode int         `json:"errcode"`
	Errmsg  string      `json:"errmsg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func renderJson(w http.ResponseWriter, jr JsonRet) {
	jsonRet, err := json.Marshal(jr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRet)
}
