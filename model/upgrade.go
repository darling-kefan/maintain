package model

import (
	"fmt"
	"time"
	"strings"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// upgrades表-项目升级记录
//
// 项目ID
// 升级前Tag
// 升级后Tag
// 升级成功minion
// 升级失败minion
// 升级状态 0-成功; 1-失败
// 操作用户ID
type Upgrade struct {
	ID          int64
	ProjectId   int64
	TagFrom     string
	TagTo       string
	MinionsSucc []string
	MinionsFail []string
	Duration    int64
	Status      int
	UserID      int64
	CreatedAt   time.Time
}

func (db *Mtdb) AddUpgrade(u *Upgrade) (id int64, err error) {
	query := "INSERT INTO `upgrades`(`project_id`,`tag_from`,`tag_to`,`minions_succ`,`minions_fail`,`duration`,`status`,`user_id`) VALUES(?,?,?,?,?,?,?,?)"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	minionsSucc := strings.Join(u.MinionsSucc, ",")
	minionsFail := strings.Join(u.MinionsFail, ",")
	r, err := execAffectingRows(stmt, u.ProjectId, u.TagFrom, u.TagTo, minionsSucc, minionsFail, u.Duration, u.Status, u.UserID)
	if err != nil {
		return 0, err
	}
	lastInsertId, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
	}
	return lastInsertId, nil
}

func (db *Mtdb) GetUpgradesByProjectId(projectId int64) ([]*Upgrade, error) {
	if projectId == 0 {
		return nil, fmt.Errorf("The param projectId: %v is empty", projectId)
	}

	query := "SELECT * FROM `upgrades` WHERE `project_id` = ?"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	upgrades, err := scanUpgrades(rows)
	if err != nil {
		return nil, err
	}
	return upgrades, nil
}

func (db *Mtdb) GetUpgradesByIds(ids []int64) ([]*Upgrade, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("The param ids: %v is empty", ids)
	}
	query := "SELECT * FROM `upgrades` WHERE `id` IN (?"+strings.Repeat(",?", len(ids)-1)+")"
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
	upgrades, err := scanUpgrades(rows)
	if err != nil {
		return nil, err
	}
	return upgrades, nil
}

func scanUpgrades(rows *sql.Rows) ([]*Upgrade, error) {
	var upgrades []*Upgrade
	for rows.Next() {
		var (
			id           int64
			projectId    int64
			tagFrom      sql.NullString
			tagTo        sql.NullString
			minionsSucc  sql.NullString
			minionsFail  sql.NullString
			duration     int64
			status       int
			userID       int64
			createdAt    sql.NullString
		)
		if err := rows.Scan(&id, &projectId, &tagFrom, &tagTo, &minionsSucc, &minionsFail, &duration, &status, &userID, &createdAt); err != nil {
			return nil, err
		}

		var msucc,mfail []string
		for _, v := range strings.Split(minionsSucc.String, ",") {
			msucc = append(msucc, v)
		}
		for _, v := range strings.Split(minionsFail.String, ",") {
			mfail = append(mfail, v)
		}
		
		createdAt1, _ := parseTimestamp(createdAt)
		upgrade := &Upgrade{
			ID:           id,
			ProjectId:    projectId,
			TagFrom:      tagFrom.String,
			TagTo:        tagTo.String,
			MinionsSucc:  msucc,
			MinionsFail:  mfail,
			Duration:     duration,
			Status:       status,
			UserID:       userID,
			CreatedAt:    createdAt1,
		}
		upgrades = append(upgrades, upgrade)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return upgrades, nil
}
