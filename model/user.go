package model

import (
	"fmt"
	"time"
	"strings"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// users表-管理员表
//
// 用户名
// 密码
// 管理员姓名
// 电话
// 邮箱
// 状态 0-启用; 1-禁用
type User struct {
	ID        int64
	Username  string
	Password  string
	Realname  string
	Phone     string
	Email     string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (db *Mtdb) GetNormalUserByUsername(uname string) (*User, error) {
	query := "SELECT * FROM `users` WHERE `username`=? AND `status` = 0"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(uname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := scanUsers(rows)
	if err != nil {
		return nil, err
	}

	return users[0], nil
}

func (db *Mtdb) GetUsersByIds(ids []int64) ([]*User, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("The param ids: %v is empty", ids)
	}
	query := "SELECT * FROM `users` WHERE `id` IN (?"+strings.Repeat(",?", len(ids)-1)+")"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	args := make([]interface{}, len(ids))
	for k, v := range ids {
		args[k] = v
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := scanUsers(rows)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func scanUsers(rows *sql.Rows) ([]*User, error) {
	var users []*User
	for rows.Next() {
		var (
			id        int64
			username  sql.NullString
			password  sql.NullString
			realname  sql.NullString
			phone     sql.NullString
			email     sql.NullString
			status    int
			createdAt sql.NullString
			updatedAt sql.NullString
		)
		if err := rows.Scan(&id, &username, &password, &realname, &phone, &email, &status, &createdAt, &updatedAt); err != nil {
			return  nil, err
		}

		createdAt1, _ := parseTimestamp(createdAt)
		updatedAt1, _ := parseTimestamp(updatedAt)
		user := &User{
			ID: id,
			Username: username.String,
			Password: password.String,
			Realname: realname.String,
			Phone: phone.String,
			Email: email.String,
			Status: status,
			CreatedAt: createdAt1,
			UpdatedAt: updatedAt1,
		}
		users = append(users, user)
		
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
