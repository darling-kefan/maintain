package model

import (
	"fmt"
	"time"
	"strings"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// projects表
//
// 项目名称
// 项目目录
// 服务器组
// 当前Tag
// 上一个Tag
// 是否在线 0-在线; 1-下线;
type Project struct {
	ID          int64
	Name        string
	RootDir     string
	CmdScript   string
	ClusterId   int64
	CurTag      string
	PreTag      string
	Online      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time

	Cluster     *Cluster
}

func (db *Mtdb) AddProject(p *Project) (id int64, err error) {
	query := "INSERT INTO `projects`(`name`,`root_dir`,`cmd_script`,`cluster_id`,`current_tag`,`previous_tag`,`created_at`) VALUES(?,?,?,?,?,?,?)"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	r, err := execAffectingRows(stmt, p.Name, p.RootDir, p.CmdScript, p.ClusterId, p.CurTag, p.PreTag, time.Now().Format(TIME_FORMAT))
	if err != nil {
		return 0, err
	}
	lastInsertId, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
	}
	return lastInsertId, nil
}

func (db *Mtdb) UpdateProject(p *Project) error {
	if p.ID == 0 {
		return fmt.Errorf("mysql: Project with unassigned ID passed into UpdateProject")
	}

	query := "UPDATE `projects` SET `name`=?,`root_dir`=?,`cmd_script`=?,`cluster_id`=?,`current_tag`=?,`previous_tag`=?,`online`=? " +
		"WHERE `id` = ?"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = execAffectingRows(stmt, p.Name, p.RootDir, p.CmdScript, p.ClusterId, p.CurTag, p.PreTag, p.Online, p.ID)
	return err
}

func (db *Mtdb) DeleteProject(id int64) error {
	if id == 0 {
		return fmt.Errorf("mysql: unassigned id passed into DeleteProject")
	}

	// query := "DELETE FROM `projects` WHERE `id` = ?"
	query := "UPDATE `projects` SET `deleted_at` = ? WHERE `id` = ?"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = execAffectingRows(stmt, time.Now().Format(TIME_FORMAT), id)
	return err
}

func (db *Mtdb) DeleteProjects(ids []int64) error {
	if len(ids) == 0 {
		return fmt.Errorf("mysql: unassigned ids passed into DeleteProjects")
	}

	// query := "DELETE FROM `projects` WHERE `id` IN (?"+strings.Repeat(",?", len(ids)-1)+")"
	query := "UPDATE `projects` SET `deleted_at` = ? where `id` (?"+strings.Repeat(",?", len(ids)-1)+")"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	args := make([]interface{}, len(ids)+1)
	args = append(args, time.Now().Format(TIME_FORMAT))
	for _, v := range ids {
		args = append(args, v)
	}
	_, err = execAffectingRows(stmt, args...)
	return err
}

// mark: 0-不限; 1-在线; 2-下线
func (db *Mtdb) GetProjects(mark int) ([]*Project, error) {
	var args []interface{}
	query := "SELECT * FROM `projects` WHERE `deleted_at` IS NULL"
	if mark == 1 || mark == 2 {
		query = query + " AND `online` = ?"
		args = append(args, mark-1)
	}

	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects, err := scanProjects(rows)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (db *Mtdb) GetProjectsByIds(ids []int64) ([]*Project, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("The param ids: %v is empty", ids)
	}
	query := "SELECT * FROM `projects` WHERE `deleted_at` IS NULL AND `id` IN (?"+strings.Repeat(",?", len(ids)-1)+")"
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

	projects, err := scanProjects(rows)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func scanProjects(rows *sql.Rows) ([]*Project, error) {
	var projects []*Project
	for rows.Next() {
		var (
			id           int64
			name         sql.NullString
			rootDir      sql.NullString
			cmdScript    sql.NullString
			clusterId    sql.NullInt64
			currentTag   sql.NullString
			previousTag  sql.NullString
			online       int
			createdAt    sql.NullString
			updatedAt    sql.NullString
			deletedAt    sql.NullString
		)
		if err := rows.Scan(&id, &name, &rootDir, &cmdScript, &clusterId, &currentTag, &previousTag, &online, &createdAt, &updatedAt, &deletedAt); err != nil {
			return nil, err
		}

		createdAt1, _ := parseTimestamp(createdAt)
		updatedAt1, _ := parseTimestamp(updatedAt)
		deletedAt1, _ := parseTimestamp(deletedAt)
		project := &Project{
			ID:        id,
			Name:      name.String,
			RootDir:   rootDir.String,
			CmdScript: cmdScript.String,
			ClusterId: clusterId.Int64,
			CurTag:    currentTag.String,
			PreTag:    previousTag.String,
			Online:    online,
			CreatedAt: createdAt1,
			UpdatedAt: updatedAt1,
			DeletedAt: deletedAt1,
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

