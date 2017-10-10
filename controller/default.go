package controller

import (
	"fmt"
	"net/http"
	"crypto/md5"

	"maintain/config"
)

func IndexAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := &struct{
			AppEnv   string
			BasePath string
		}{
			AppEnv:   config.AppEnv,
			BasePath: config.BasePath,
		}
		renderTemplate(w, "index", data)
	})
}

func LoginAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := env.SessionManager.SessionStart(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer sess.SessionRelease(w)

		userId := sess.Get("UserID")
		if r.Method == "GET" {
			if (userId == nil) {
				data := &struct{
					AppEnv   string
					BasePath string
				}{
					AppEnv:   config.AppEnv,
					BasePath: config.BasePath,
				}
				renderTemplate(w, "login", data)
				return
			}
			http.Redirect(w, r, config.BasePath, http.StatusFound)
		} else {
			err := r.ParseForm()
			if (err != nil) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			username := r.PostFormValue("username")
			password := r.PostFormValue("password")

			user, err := env.Mtdb.GetNormalUserByUsername(username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			authPart := fmt.Sprintf("%x", md5.Sum([]byte(username+password)))
			authString := fmt.Sprintf("%x", md5.Sum([]byte(authPart+config.AuthSalt)))

			if (user.Password != authString) {
				jr := JsonRet{
					Errcode: 1,
					Errmsg: "用户名或密码不正确",
				}
				renderJson(w, jr)
				return
			}

			sess.Set("UserID", user.ID)

			jr := JsonRet{
				Errcode: 0,
				Errmsg: "success",
			}
			renderJson(w, jr)
		}
	})
}

func LogoutAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, _ := env.SessionManager.SessionStart(w, r)
		userId := sess.Get("UserID")
		if userId != nil {
			sess.Delete("UserID")
		}
		http.Redirect(w, r, config.BasePath+"login", http.StatusFound)
	})
}

func ProfileAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
