package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/go-sql-driver/mysql"

	"maintain/config"
)

var createTablesStatements = []string{
	`CREATE DATABASE IF NOT EXISTS maintain DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`,
	`USE maintain;`,
	
	`CREATE TABLE IF NOT EXISTS projects (
        `+"`id`"+` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
        `+"`name`"+` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '项目名称',
        `+"`root_dir`"+` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '项目根目录',
        `+"`cmd_script`"+` VARCHAR(500) NOT NULL DEFAULT '' COMMENT 'shell文件，用于git部署',
        `+"`cluster_id`"+` INT(11) NOT NULL DEFAULT 0 COMMENT '服务器组id',
        `+"`current_tag`"+` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '当前git tag',
        `+"`previous_tag`"+` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '上一个git tag',
        `+"`online`"+` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '是否在线，0-在线; 1-下线',
        `+"`created_at`"+` TIMESTAMP NULL DEFAULT NULL COMMENT '创建时间',
        `+"`updated_at`"+` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
        `+"`deleted_at`"+` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
        PRIMARY KEY (`+"`id`"+`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci`,

	`CREATE TABLE IF NOT EXISTS clusters (
        `+"`id`"+` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
        `+"`name`"+` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '组名称',
        `+"`online`"+` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '是否在线，0-在线; 1-下线',
        `+"`created_at`"+` TIMESTAMP NULL DEFAULT NULL COMMENT '创建时间',
        `+"`updated_at`"+` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
        `+"`deleted_at`"+` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
        PRIMARY KEY (`+"`id`"+`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci`,

	`CREATE TABLE IF NOT EXISTS minions (
        `+"`id`"+` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
        `+"`name`"+` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '机器名称',
        `+"`ipv4_internal`"+` VARCHAR(15) NOT NULL DEFAULT '' COMMENT '内网ip地址',
        `+"`ipv4_external`"+` VARCHAR(15) NOT NULL DEFAULT '' COMMENT '外网ip地址',
        `+"`online`"+` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '是否在线，0-在线; 1-下线',
        `+"`created_at`"+` TIMESTAMP NULL DEFAULT NULL COMMENT '创建时间',
        `+"`updated_at`"+` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
        `+"`deleted_at`"+` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
        PRIMARY KEY (`+"`id`"+`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci`,

	`CREATE TABLE IF NOT EXISTS cluster_minion (
        `+"`id`"+` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
        `+"`cluster_id`"+` INT(11) NOT NULL DEFAULT 0 COMMENT '服务器组id',
        `+"`minion_id`"+` INT(11) NOT NULL DEFAULT 0 COMMENT 'minion id',
        `+"`created_at`"+` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
        PRIMARY KEY (`+"`id`"+`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci`,
	
	`CREATE TABLE IF NOT EXISTS upgrades (
        `+"`id`"+` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
        `+"`project_id`"+` INT(11) NOT NULL DEFAULT 0 COMMENT '项目id',
        `+"`tag_from`"+` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '升级前tag',
        `+"`tag_to`"+` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '升级后tag',
        `+"`minions_succ`"+` TEXT COMMENT '升级成功的minions',
        `+"`minions_fail`"+` TEXT COMMENT '升级失败的minions',
        `+"`duration`"+` INT(11) NOT NULL DEFAULT 0 COMMENT '升级耗时，单位毫秒',
        `+"`status`"+` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '状态，0-成功; 1-失败',
        `+"`user_id`"+` INT(11) NOT NULL DEFAULT 0 COMMENT '用户id',
        `+"`created_at`"+` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
        PRIMARY KEY (`+"`id`"+`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci`,

	`CREATE TABLE IF NOT EXISTS users (
        `+"`id`"+` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
        `+"`username`"+` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '用户名',
        `+"`password`"+` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '密码',
        `+"`realname`"+` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '姓名',
        `+"`phone`"+` VARCHAR(25) NOT NULL DEFAULT '' COMMENT '手机号码',
        `+"`email`"+` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '邮箱',
        `+"`status`"+` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '状态0-正常;1-禁用',
        `+"`created_at`"+` TIMESTAMP NULL DEFAULT NULL COMMENT '创建时间',
        `+"`updated_at`"+` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
        PRIMARY KEY (`+"`id`"+`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci`,
}

type MySQLConfig struct{
	Username string
	Password string
	Host     string
	Port     int
}

// [[username[:password]@]tcp([host]:port)/db]
func (c MySQLConfig) dataStoreName(databaseName string) string {
	var cred string
	if c.Username != "" {
		cred = c.Username
		if c.Password != "" {
			cred = cred + ":" + c.Password
		}
		cred = cred + "@"
	}

	return fmt.Sprintf("%stcp([%s]:%d)/%s", cred, c.Host, c.Port, databaseName)
}

func (c MySQLConfig) ensureTablesExists() error {
	conn, err := sql.Open("mysql", c.dataStoreName(""))
	if err != nil {
		return fmt.Errorf("mysql: could not get a connection: %v", err)
	}
	defer conn.Close()

	// check the connection
	if conn.Ping() == driver.ErrBadConn {
		return fmt.Errorf("mysql: could not connect to the database. " +
		       "could be bad address, or this address is not whitelisted for access.")
	}

	if _, err := conn.Exec("USE maintain"); err != nil {
		// MySQL error 1049 is "database does not exist"
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1049 {
			return createTables(conn)
		}
	}

	if _, err := conn.Exec("describe projects"); err != nil {
		// MySQL error 1146 is "table does not exist"
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
			return createTables(conn)
		}
		// Unknown error.
		return fmt.Errorf("mysql: could not connect to the database: %v", err)
	}

	if _, err := conn.Exec("describe clusters"); err != nil {
		// MySQL error 1146 is "table does not exist"
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
			return createTables(conn)
		}
		// Unknown error.
		return fmt.Errorf("mysql: could not connect to the database: %v", err)
	}

	if _, err := conn.Exec("describe minions"); err != nil {
		// MySQL error 1146 is "table does not exist"
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
			return createTables(conn)
		}
		// Unknown error.
		return fmt.Errorf("mysql: could not connect to the database: %v", err)
	}

	if _, err := conn.Exec("describe cluster_minion"); err != nil {
		// MySQL error 1146 is "table does not exist"
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
			return createTables(conn)
		}
		// Unknown error.
		return fmt.Errorf("mysql: could not connect to the database: %v", err)
	}

	if _, err := conn.Exec("describe upgrades"); err != nil {
		// MySQL error 1146 is "table does not exist"
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
			return createTables(conn)
		}
		// Unknown error.
		return fmt.Errorf("mysql: could not connect to the database: %v", err)
	}

	if _, err := conn.Exec("describe users"); err != nil {
		// MySQL error 1146 is "table does not exist"
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
			return createTables(conn)
		}
		// Unknown error.
		return fmt.Errorf("mysql: could not connect to the database: %v", err)
	}

	// Add Administrator
	var userId int64
	err = conn.QueryRow("SELECT `id` FROM `users` WHERE `username` = ?", "admin").Scan(&userId)
	switch {
	case err == sql.ErrNoRows:
		_, err := conn.Exec("INSERT INTO `users`(`username`,`password`,`realname`,`phone`,`email`,`status`) VALUES('admin','77a40784aeb6546fe406493337f0664b','administrator','15010240697','tangshouqiang@tvmining.com',0);")
		if err != nil {
			return fmt.Errorf("mysql: %s", err.Error())
		}
	case err != nil:
		return fmt.Errorf("mysql: %s", err.Error())
	}
	
	return nil
}

func createTables(conn *sql.DB) error {
	for _, sql := range createTablesStatements {
		_, err := conn.Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	conf := MySQLConfig{
		Username: config.DbUserName,
		Password: config.DbPassword,
		Host:     config.DbHost,
		Port:     config.DbPort,
	}
	if err := conf.ensureTablesExists(); err != nil {
		fmt.Println(err)
	}
}
