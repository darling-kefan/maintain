package controller

import (
	"fmt"
	"sort"
	"time"
	_ "reflect"
	"strconv"
	"strings"
	"net/url"
	"net/http"

	"maintain/model"
	"maintain/config"
	"maintain/helpers"
)

func DeploysAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projects, err := env.Mtdb.GetProjects(1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
		renderTemplate(w, "deploys", data)
	})
}

func DeployInfoAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryForm, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		idSl, ok := queryForm["projectId"]
		if !ok {
			renderJson(w, JsonRet{Errcode:1, Errmsg:"The param: projectId is not existed."})
			return
		}
		id, err := strconv.ParseInt(idSl[0], 10, 64)
		if err != nil {
			renderJson(w, JsonRet{Errcode:1, Errmsg:err.Error()})
			return
		}

		projects, err := env.Mtdb.GetProjectsByIds([]int64{id})
		if err != nil {
			renderJson(w, JsonRet{Errcode:1, Errmsg:err.Error()})
			return
		}
		if len(projects) == 0 {
			renderJson(w, JsonRet{Errcode:1, Errmsg:err.Error()})
			return
		}

		// 根据项目id获取所有minion
		minions, err := env.Mtdb.GetMinionsByClusterId(projects[0].ClusterId)
		onlineMinion := minions[0].Name
		
		saltClient, err := config.NewSaltClient()
		if err != nil {
			renderJson(w, JsonRet{Errcode:1, Errmsg:err.Error()})
			return
		}
		// 获取当前在线tag
		args := []string{"cd "+projects[0].RootDir+"; git describe"}
		ret, err := saltClient.Cmd(onlineMinion, "cmd.run", args, "glob")
		if err != nil {
			renderJson(w, JsonRet{Errcode:1, Errmsg:err.Error()})
			return
		}
		if !helpers.CheckTag(ret[onlineMinion].(string)) {
			renderJson(w, JsonRet{Errcode:1, Errmsg:ret[onlineMinion].(string)})
			return
		}
		onTag := ret[onlineMinion].(string)
		
		// 获取tag list
		args = []string{"cd "+projects[0].RootDir+"; git fetch origin --tags; git tag"}
		ret, err = saltClient.Cmd(config.ClusterWebRelease, "cmd.run", args, "glob")
		if err != nil {
			renderJson(w, JsonRet{Errcode:1, Errmsg:err.Error()})
			return
		}
		tags := strings.Split(ret[config.ClusterWebRelease].(string), "\n")
		var newTags []string
		for _, t := range tags {
			if helpers.CheckTag(t) {
				newTags = append(newTags, t)
			}
		}
		if len(newTags) == 0 {
			renderJson(w, JsonRet{Errcode:1, Errmsg:ret[config.ClusterWebRelease].(string)})
			return
		}
		sort.Sort(helpers.Tags(newTags))
		latestTag := newTags[0]

		data := &struct{
			Project   *model.Project
			OnTag     string
			LatestTag string
			Tags      []string
		}{
			Project:   projects[0],
			OnTag:     onTag,
			LatestTag: latestTag,
			Tags:      newTags,
		}
		jr := JsonRet{
			Errcode: 0,
			Errmsg:  "success",
			Data:    data,
		}
		renderJson(w, jr)
	})
}

func DeployExecAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if r.PostFormValue("projectId") == "" || r.PostFormValue("toTag") == "" {
			http.Error(w, "Missing parameters.", http.StatusInternalServerError)
			return
		}
		projectId, err := strconv.ParseInt(r.PostFormValue("projectId"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		projects, err := env.Mtdb.GetProjectsByIds([]int64{projectId})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(projects) == 0 {
			http.Error(w, fmt.Sprintf("The project: %d is not existed.", projectId), http.StatusInternalServerError)
			return
		}

		// 根据项目id获取所有minion
		minions, err := env.Mtdb.GetMinionsByClusterId(projects[0].ClusterId)
		var minionNames []string
		for _, v := range minions {
			minionNames = append(minionNames, v.Name)
		}
		
		saltClient, err := config.NewSaltClient()
		if err != nil {
			renderJson(w, JsonRet{Errcode:1, Errmsg:err.Error()})
			return
		}
		// 项目部署
		var stime = time.Now()
		args := []string{"salt://"+projects[0].CmdScript[len(config.SaltFileRoots):], r.PostFormValue("toTag")}

		// 鉴于批量处理minions，会造成码云git报错，错误如下：
		//
		// handShake upstream error The current repository requires restricted access
		// exec request failed on channel 0
		// fatal: The remote end hung up unexpectedly
		//
		// 因此，需依次执行minions并停顿片刻

		// 注释掉一次性执行方式
		/*ret, err := saltClient.Cmd(minionNames, "cmd.script", args, "list")
		if err != nil {
			renderJson(w, JsonRet{Errcode:1, Errmsg:err.Error()})
			return
		}*/

		// 分别执行minion
		ret := make(map[string]interface{}, len(minionNames))
		for _, minionName := range minionNames {
			r, err := saltClient.Cmd(minionName, "cmd.script", args, "glob")
			if err != nil {
				renderJson(w, JsonRet{Errcode:1, Errmsg:err.Error()})
				return
			}
			for k, v := range r {
				ret[k] = v
			}
			// Sleep 1 second
			time.Sleep(1*time.Second)
		}

		upgradeDuration := time.Now().Sub(stime).Nanoseconds() / 1000 / 1000

		// 更新projects.previous_tag, projects.current_tag
		onTag := r.PostFormValue("onTag")
		if onTag == "" {
			onTag = r.PostFormValue("toTag")
		}
		proj := projects[0]
		proj.PreTag = onTag
		proj.CurTag = r.PostFormValue("toTag")
		if err := env.Mtdb.UpdateProject(proj); err != nil {
			renderJson(w, JsonRet{Errcode:1, Errmsg:err.Error()})
			return
		}

		// 保存升级日志
		var minionsSucc, minionsFail []string
		for k, v := range ret {
			vv := v.(map[string]interface{})
			if retcode := vv["retcode"].(float64); retcode == 0 {
				minionsSucc = append(minionsSucc, k)
			} else {
				minionsFail = append(minionsFail, k)
			}
		}
		var upgradeStatus int
		if len(minionsSucc) != len(minions) {
			upgradeStatus = 1
		}
		u := &model.Upgrade{
			ProjectId:   proj.ID,
			TagFrom:     proj.PreTag,
			TagTo:       proj.CurTag,
			MinionsSucc: minionsSucc,
			MinionsFail: minionsFail,
			Duration:    upgradeDuration,
			Status:      upgradeStatus,
		}
		if _, err = env.Mtdb.AddUpgrade(u); err != nil {
			renderJson(w, JsonRet{Errcode:1, Errmsg:err.Error()})
			return
		}

		jr := JsonRet{
			Errcode: 0,
			Errmsg:  "success",
			Data:    ret,
		}
		renderJson(w, jr)
	})
}

func UpgradeAction(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
	})
}
