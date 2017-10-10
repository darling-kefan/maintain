package model

import (
	"fmt"
)

type MySQLConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	DbName   string
}

// [[username[:password]@]tcp([host]:port)/db]
func (c MySQLConfig) dataStoreName() string {
	var cred string
	if c.Username != "" {
		cred = c.Username
		if c.Password != "" {
			cred = cred + ":" + c.Password
		}
		cred = cred + "@"
	}

	return fmt.Sprintf("%stcp([%s]:%d)/%s", cred, c.Host, c.Port, c.DbName)
}

