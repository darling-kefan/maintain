package controller

import (
	_ "fmt"
	"strconv"
	"net/http"
	_ "html/template"
	_ "io/ioutil"

	"maintain/config"
	"maintain/model"
)

func ListProjectsAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projects, err := env.Mtdb.GetProjects(1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clusters, err := env.Mtdb.GetClusters(0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		clustersMap := make(map[int64]*model.Cluster, len(clusters))
		for _, c := range clusters {
			clustersMap[c.ID] = c
		}
		for _, v := range projects {
			if cl, ok := clustersMap[v.ClusterId]; ok {
				v.Cluster = cl
			}
		}

		data := &struct{
			Projects []*model.Project
			AppEnv   string
			BasePath string
		}{
			Projects: projects,
			AppEnv:   config.AppEnv,
			BasePath: config.BasePath,
		}
		renderTemplate(w, "projects", data)
	})
}

func AddProjectAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// @todo 机器组
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
		renderTemplate(w, "project_add", data)
	})
}

func SaveProjectAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		clusterId, err := strconv.ParseInt(r.PostFormValue("clusterId"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p := &model.Project{
			Name:      r.PostFormValue("name"),
			RootDir:   r.PostFormValue("rootDir"),
			CmdScript: r.PostFormValue("cmdScript"),
			ClusterId: clusterId,
			CurTag:    r.PostFormValue("curTag"),
			PreTag:    r.PostFormValue("preTag"),
		}
		id, err := env.Mtdb.AddProject(p)
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

func EditProjectAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.URL.Path[len("/project/edit/"):], 10, 64)
		if err != nil {
			http.Error(w, "The param project_id is unvalid.", http.StatusInternalServerError)
			return
		}
		projects, err := env.Mtdb.GetProjectsByIds([]int64{id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(projects) == 0 {
			http.Redirect(w, r, config.BasePath+"projects", http.StatusFound)
			return
		}
		clusters, err := env.Mtdb.GetClusters(1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := &struct{
			Project  *model.Project
			Clusters []*model.Cluster
			AppEnv   string
			BasePath string
		}{
			Project:  projects[0],
			Clusters: clusters,
			AppEnv:   config.AppEnv,
			BasePath: config.BasePath,
		}
		// @机器组
		renderTemplate(w, "project_edit", data)
	})
}

func ModifyProjectAction(env *config.Env) http.Handler {
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
		clusterId, err := strconv.ParseInt(r.PostFormValue("clusterId"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p := &model.Project{
			ID:        id,
			Name:      r.PostFormValue("name"),
			RootDir:   r.PostFormValue("rootDir"),
			CmdScript: r.PostFormValue("cmdScript"),
			ClusterId: clusterId,
			CurTag:    r.PostFormValue("curTag"),
			PreTag:    r.PostFormValue("preTag"),
			Online:    0,
		}
		err = env.Mtdb.UpdateProject(p)
		if err != nil {
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

func DelProjectAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//http.Error(w, "***************test Error********************", http.StatusInternalServerError)
		//return

		id, err := strconv.ParseInt(r.URL.Path[len("/project/delete/"):], 10, 64)
		if err != nil {
			http.Error(w, "The param project_id is unvalid.", http.StatusInternalServerError)
			return
		}
		if err := env.Mtdb.DeleteProject(id); err != nil {
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
