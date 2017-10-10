package config

import (
	"os"
	"github.com/darling-kefan/go-salt"
	"github.com/darling-kefan/toolkit/cfg"
	"github.com/astaxie/beego/session"

	"maintain/model"
)

var SaltHost string
var	SaltPort string
var	SaltUserName string
var	SaltPassword string
var	SaltDebug bool
var	SaltSslSkipVerify bool

var	SaltFileRoots string
var	ClusterWebRelease string
var	ClusterWebOnline  string

var	DbUserName string
var DbPassword string
var	DbHost string
var	DbPort int
var	DbName string

var	WebRoot string
var	AuthSalt string
var BasePath string
var AppEnv string

func init() {
	if pwd, err := os.Getwd(); err == nil {
		cfg.Init(cfg.FileProvider{
			Filename: pwd+"/.env",
		})
	} else {
		panic(err.Error())
	}

	SaltHost = cfg.MustString("SALT_HOST")
	SaltPort = cfg.MustString("SALT_PORT")
	SaltUserName = cfg.MustString("SALT_USERNAME")
	SaltPassword = cfg.MustString("SALT_PASSWORD")
	SaltDebug = cfg.MustBool("SALT_DEBUG")
	SaltSslSkipVerify = cfg.MustBool("SALT_SSLSKIPVERIFY")

	SaltFileRoots = cfg.MustString("SALT_FILE_ROOTS")
	ClusterWebRelease = cfg.MustString("CLUSTER_WEB_RELEASE")
	ClusterWebOnline = cfg.MustString("CLUSTER_WEB_ONLINE")

	DbUserName  = cfg.MustString("DB_USERNAME")
	DbPassword = cfg.MustString("DB_PASSWORD")
	DbHost = cfg.MustString("DB_HOST")
	DbPort = cfg.MustInt("DB_PORT")
	DbName = cfg.MustString("DB_NAME")

	WebRoot = cfg.MustString("WEB_ROOT")
	AuthSalt = cfg.MustString("AUTH_SALT")
	BasePath = cfg.MustString("BASE_PATH")
	AppEnv = cfg.MustString("APP_ENV")
}

type Env struct {
	SessionManager    *session.Manager //Session Manager
	Mtdb              *model.Mtdb      //Maintain database connection pool
}

func NewSaltClient() (*salt.Client, error) {
	saltConf := salt.Config{
		Host:          SaltHost,
		Port:          SaltPort,
		Username:      SaltUserName,
		Password:      SaltPassword,
		Debug:         SaltDebug,
		SSLSkipVerify: SaltSslSkipVerify,
	}

	// saltClient := &salt.Client{}
	saltClient, err := salt.NewClient(saltConf)
	if err != nil {
		return nil, err
	}
	return saltClient, nil
}

func NewInit() (*Env, error) {
	var err error
	var env *Env

	mtdbConf := model.MySQLConfig{
		Username:  DbUserName,
		Password:  DbPassword,
		Host:      DbHost,
		Port:      DbPort,
		DbName:    DbName,
	}
	mtdb, err := model.NewMtdb(mtdbConf)
	if err != nil {
		return env, err
	}

	sessionManagerConf := session.ManagerConfig{
		CookieName: "gosessionid",
		CookieLifeTime: 3600,
		Gclifetime: 3600,
		Maxlifetime: 3600,
		EnableSetCookie: true,
	}
	sessionManager, err := session.NewManager("memory", &sessionManagerConf)
	if err != nil {
		return env, err
	}
	go sessionManager.GC()
	
	env = &Env{
		SessionManager: sessionManager,
		Mtdb:           mtdb,
	}
	return env, nil
}
