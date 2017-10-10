package controller

import (
	"strings"
	"strconv"
	"net/http"

	"maintain/config"
	"maintain/model"
)

func ListClustersAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clusters, err := env.Mtdb.GetClusters(1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := &struct{
			Clusters []*model.Cluster
			AppEnv   string
			BasePath string
		}{
			Clusters: clusters,
			AppEnv:   config.AppEnv,
			BasePath: config.BasePath,
		}
		renderTemplate(w, "clusters", data)
	})
}

func AddClusterction(env *config.Env) http.Handler {
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
		renderTemplate(w, "cluster_add", data)
	})
}

func SaveClusterAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		c := &model.Cluster{
			Name: r.PostFormValue("name"),
			Online: 0,
		}
		id, err := env.Mtdb.AddCluster(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, v := range strings.Split(r.PostFormValue("minions"), ",") {
			if vv, err := strconv.ParseInt(v, 10, 64); err == nil {
				if _, err := env.Mtdb.AddClusterMinion(id, vv); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		
		js := JsonRet{
			Errcode: 0,
			Errmsg:  "success",
		}
		renderJson(w, js)
	})
}

func EditClusterAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.URL.Path[len("/cluster/edit/"):], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clusters, err := env.Mtdb.GetClustersByIds([]int64{id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if clusterMinions, err := env.Mtdb.GetMinionsByClusterId(id); err == nil {
			clusters[0].Minions = clusterMinions
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		minions, err := env.Mtdb.GetMinions(1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := &struct{
			Cluster  *model.Cluster
			Minions  []*model.Minion
			AppEnv   string
			BasePath string
		}{
			Cluster:  clusters[0],
			Minions:  minions,
			AppEnv:   config.AppEnv,
			BasePath: config.BasePath,
		}
		renderTemplate(w, "cluster_edit", data)
	})
}

func ModifyClusterAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := strconv.ParseInt(r.PostFormValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		c := &model.Cluster{
			ID: id,
			Name: r.PostFormValue("name"),
			Online: 0,
		}
		if err := env.Mtdb.UpdateCluster(c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := env.Mtdb.DeleteClusterMinions(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, v := range strings.Split(r.PostFormValue("minions"), ",") {
			if vv, err := strconv.ParseInt(v, 10, 64); err == nil {
				if _, err := env.Mtdb.AddClusterMinion(id, vv); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		
		js := JsonRet{
			Errcode: 0,
			Errmsg:  "success",
		}
		renderJson(w, js)
	})
}

func DelClusterAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.URL.Path[len("/cluster/delete/"):], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := env.Mtdb.DeleteCluster(id); err != nil {
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
