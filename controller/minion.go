package controller

import (
	"strconv"
	"net/http"

	"maintain/config"
	"maintain/model"
)

func ListMinionsAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		minions, err := env.Mtdb.GetMinions(1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := &struct{
			Minions  []*model.Minion
			AppEnv   string
			BasePath string
		}{
			Minions:  minions,
			AppEnv:   config.AppEnv,
			BasePath: config.BasePath,
		}
		renderTemplate(w, "minions", data)
	})
}

func AddMinionAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := &struct{
			AppEnv   string
			BasePath string
		}{
			AppEnv:   config.AppEnv,
			BasePath: config.BasePath,
		}
		renderTemplate(w, "minion_add", data)
	})
}

func SaveMinionAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		m := &model.Minion{
			Name:         r.PostFormValue("name"),
			Ipv4Internal: r.PostFormValue("ipv4Internal"),
			Ipv4External: r.PostFormValue("ipv4External"),
			Online:       0,
		}
		id, err := env.Mtdb.AddMinion(m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jr := JsonRet{
			Errcode: 0,
			Errmsg:  "success",
			Data:    id,
		}
		renderJson(w, jr)
	})
}

func EditMinionAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.URL.Path[len("/minion/edit/"):], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		minions, err := env.Mtdb.GetMinionsByIds([]int64{id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(minions) == 0 {
			http.Redirect(w, r, config.BasePath+"minions", http.StatusFound)
			return
		}
		data := &struct{
			Minion *model.Minion
			AppEnv   string
			BasePath string
		}{
			Minion: minions[0],
			AppEnv:   config.AppEnv,
			BasePath: config.BasePath,
		}
		renderTemplate(w, "minion_edit", data)
	})
}

func ModifyMinionAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := strconv.ParseInt(r.PostFormValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		m := &model.Minion{
			ID:           id,
			Name:         r.PostFormValue("name"),
			Ipv4Internal: r.PostFormValue("ipv4Internal"),
			Ipv4External: r.PostFormValue("ipv4External"),
			Online:       0,
		}
		if err := env.Mtdb.UpdateMinion(m); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jr := JsonRet{
			Errcode: 0,
			Errmsg:  "success",
		}
		renderJson(w, jr)
	})
}

func DelMinionAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.URL.Path[len("/minion/delete/"):], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := env.Mtdb.DeleteMinion(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jr := JsonRet{
			Errcode: 0,
			Errmsg:  "success",
		}
		renderJson(w, jr)
	})
}
