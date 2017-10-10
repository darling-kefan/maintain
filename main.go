package main

import (
	"log"
	"net/http"

	"maintain/config"
	"maintain/controller"
)

// 全局变量
var env *config.Env
var err error

func main() {
	// 初始化数据库连接池, 日志句柄等
	env, err = config.NewInit()
	if err != nil {
		log.Panic(err)
	}

	// 启动静态文件服务
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir(config.WebRoot+"resources/"))))

	// 通用操作
	http.Handle("/", use(controller.IndexAction(env), authenticate))
	http.Handle("/login", controller.LoginAction(env))
	http.Handle("/logout", controller.LogoutAction(env))

	// Project项目相关
	http.Handle("/projects", use(controller.ListProjectsAction(env), authenticate))
	http.Handle("/project/add", use(controller.AddProjectAction(env), authenticate))
	http.Handle("/project/edit/", use(controller.EditProjectAction(env), authenticate))
	http.Handle("/project/save", use(controller.SaveProjectAction(env), authenticate))
	http.Handle("/project/modify", use(controller.ModifyProjectAction(env), authenticate))
	http.Handle("/project/delete/", use(controller.DelProjectAction(env), authenticate))

	// 机器组相关
	http.Handle("/clusters", use(controller.ListClustersAction(env), authenticate))
	http.Handle("/cluster/add", use(controller.AddClusterction(env), authenticate))
	http.Handle("/cluster/edit/", use(controller.EditClusterAction(env), authenticate))
	http.Handle("/cluster/save", use(controller.SaveClusterAction(env), authenticate))
	http.Handle("/cluster/modify", use(controller.ModifyClusterAction(env), authenticate))
	http.Handle("/cluster/delete/", use(controller.DelClusterAction(env), authenticate))
	
	// Minion相关
	http.Handle("/minions", use(controller.ListMinionsAction(env), authenticate))
	http.Handle("/minion/add", use(controller.AddMinionAction(env), authenticate))
	http.Handle("/minion/edit/", use(controller.EditMinionAction(env), authenticate))
	http.Handle("/minion/save", use(controller.SaveMinionAction(env), authenticate))
	http.Handle("/minion/modify", use(controller.ModifyMinionAction(env), authenticate))
	http.Handle("/minion/delete/", use(controller.DelMinionAction(env), authenticate))
	
	// 生产部署相关
	http.Handle("/upgrade", use(controller.UpgradeAction(env), authenticate))
	http.Handle("/deploys", use(controller.DeploysAction(env), authenticate))
	http.Handle("/deployInfo", use(controller.DeployInfoAction(env), authenticate))
	http.Handle("/deployExec", use(controller.DeployExecAction(env), authenticate))
	
	// 监听http请求
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func use(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func authenticate(h http.Handler) http.Handler {
	authFunc := func(w http.ResponseWriter, r *http.Request) {
		sess, err := env.SessionManager.SessionStart(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userID := sess.Get("UserID")

		if userID != nil {
			h.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, config.BasePath+"login", http.StatusFound)
		}
	}

	return http.HandlerFunc(authFunc)
}
