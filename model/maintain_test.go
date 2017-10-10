// 此处包名必须为model_test
package model_test

import (
	"testing"

	"maintain/model"
	"maintain/config"
)

func _TestGetProjects(t *testing.T) {
	conf := model.MySQLConfig{
		Username: config.DbUserName,
		Password: config.DbPassword,
		Host: config.DbHost,
		Port: config.DbPort,
		DbName: config.DbName,
	}
	db, err := model.NewMtdb(conf)
	if err != nil {
		t.Errorf("Failed to connect mysql: %v", err)
	}
	defer db.Close()

	projects, err := db.GetProjects(0)
	if err != nil {
		t.Errorf("%v", err)
	}

	for _, p := range projects {
		t.Logf("%v\n", p)
	}
}

func _TestGetProjectsByIds(t *testing.T) {
	conf := model.MySQLConfig{
		Username: config.DbUserName,
		Password: config.DbPassword,
		Host: config.DbHost,
		Port: config.DbPort,
		DbName: config.DbName,
	}
	db, err := model.NewMtdb(conf)
	if err != nil {
		t.Errorf("Failed to connect mysql: %v", err)
	}
	defer db.Close()

	projects, err := db.GetProjectsByIds([]int64{1, 2})
	if err != nil {
		t.Errorf("%v", err)
	}

	for _, p := range projects {
		t.Logf("%v\n", p)
	}
}

func _TestDeleteProject(t *testing.T) {
	conf := model.MySQLConfig{
		Username: config.DbUserName,
		Password: config.DbPassword,
		Host: config.DbHost,
		Port: config.DbPort,
		DbName: config.DbName,
	}
	db, err := model.NewMtdb(conf)
	if err != nil {
		t.Errorf("Failed to connect mysql: %v", err)
	}
	defer db.Close()

	err = db.DeleteProjects([]int64{2})
	if err != nil {
		t.Errorf("%v", err)
	}
}

func _TestAddProject(t *testing.T) {
	conf := model.MySQLConfig{
		Username: config.DbUserName,
		Password: config.DbPassword,
		Host: config.DbHost,
		Port: config.DbPort,
		DbName: config.DbName,
	}
	db, err := model.NewMtdb(conf)
	if err != nil {
		t.Errorf("Failed to connect mysql: %v", err)
	}
	defer db.Close()

	p := &model.Project{
		Name: "v2",
		RootDir: "/opt/adcloud-v2/",
		ClusterId: 1,
		CurTag: "v0.1.23",
		PreTag: "v0.1.24",
		Online: 0,
	}
	id, err := db.AddProject(p)
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Log(id)
}

func _TestUpdateProject(t *testing.T) {
	conf := model.MySQLConfig{
		Username: config.DbUserName,
		Password: config.DbPassword,
		Host: config.DbHost,
		Port: config.DbPort,
		DbName: config.DbName,
	}
	db, err := model.NewMtdb(conf)
	if err != nil {
		t.Errorf("Failed to connect mysql: %v", err)
	}
	defer db.Close()

	p := &model.Project{
		ID: 3,
		Name: "v2",
		RootDir: "/opt/adcloud-v2/",
		ClusterId: 1,
		CurTag: "v0.1.24",
		PreTag: "v0.1.25",
		Online: 1,
	}
	err = db.UpdateProject(p)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func _TestAddCluster(t *testing.T) {
	conf := model.MySQLConfig{
		Username: config.DbUserName,
		Password: config.DbPassword,
		Host: config.DbHost,
		Port: config.DbPort,
		DbName: config.DbName,
	}
	db, err := model.NewMtdb(conf)
	if err != nil {
		t.Errorf("Failed to connect mysql: %v", err)
	}
	defer db.Close()

	p := &model.Cluster{
		Name: "web-v2",
		Online: 0,
	}
	id, err := db.AddCluster(p)
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Log(id)
}

func TestGetClusters(t *testing.T) {
	t.Log(config.DbUserName, config.DbPassword, config.DbHost, config.DbPort, config.DbName)

	conf := model.MySQLConfig{
		Username: config.DbUserName,
		Password: config.DbPassword,
		Host: config.DbHost,
		Port: config.DbPort,
		DbName: config.DbName,
	}
	db, err := model.NewMtdb(conf)
	if err != nil {
		t.Errorf("Failed to connect mysql: %v", err)
	}
	defer db.Close()

	clusters, err := db.GetClusters(0)
	if err != nil {
		t.Errorf("%v", err)
	}

	for _, p := range clusters {
		t.Logf("%v\n", p)
	}
}
