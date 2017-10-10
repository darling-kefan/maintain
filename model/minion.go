package model

import (
	"fmt"
	"time"
	"strings"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// minions表
//
// minion名称
// minion外网地址
// minion内网地址
// minion在线状态，0-在线; 1-下线;
type Minion struct {
	ID           int64
	Name         string
	Ipv4Internal string
	Ipv4External string
	Online       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time

	Clusters     []Cluster
}

func (db *Mtdb) AddMinion(m *Minion) (id int64, err error) {
	query := "INSERT INTO `minions`(`name`,`ipv4_internal`,`ipv4_external`,`created_at`) VALUES(?,?,?,?)"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	r, err := execAffectingRows(stmt, m.Name, m.Ipv4Internal, m.Ipv4External, time.Now().Format(TIME_FORMAT))
	if err != nil {
		return 0, err
	}
	lastInsertId, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
	}
	return lastInsertId, nil
}

func (db *Mtdb) UpdateMinion(m *Minion) error {
	if m.ID == 0 {
		return fmt.Errorf("mysql: Project with unassigned ID passed into UpdateMinion")
	}

	query := "UPDATE `minions` SET `name`=?,`ipv4_internal`=?,`ipv4_external`=?,`online`=? WHERE `id` = ?"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = execAffectingRows(stmt, m.Name, m.Ipv4Internal, m.Ipv4External, m.Online, m.ID)
	return err
}

func (db *Mtdb) DeleteMinion(id int64) error {
	if id == 0 {
		return fmt.Errorf("mysql: unassigned id passed into DeleteMinion")
	}

	query := "UPDATE `minions` SET `deleted_at` = ? WHERE `id` = ?"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = execAffectingRows(stmt, time.Now().Format(TIME_FORMAT), id)
	return err
}

// mark: 0-不限; 1-在线; 2-下线
func (db *Mtdb) GetMinions(mark int) ([]*Minion, error) {
	var args []interface{}
	query := "SELECT * FROM `minions` WHERE `deleted_at` IS NULL"
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
	minions, err := scanMinions(rows)
	if err != nil {
		return nil, err
	}

	return minions, nil
}

func (db *Mtdb) GetMinionsByIds(ids []int64) ([]*Minion, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("The param ids: %v is empty", ids)
	}
	query := "SELECT * FROM `minions` WHERE `deleted_at` IS NULL AND `id` IN (?"+strings.Repeat(",?", len(ids)-1)+")"
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
	minions, err := scanMinions(rows)
	if err != nil {
		return nil, err
	}
	return minions, nil
}

func (db *Mtdb) GetMinionsByClusterId(clusterId int64) ([]*Minion, error) {
	if clusterId == 0 {
		return nil, fmt.Errorf("The param clusterId: %v is empty", clusterId)
	}

	query := "SELECT `minions`.`id` AS `id`,`minions`.`name` AS `name`,`minions`.`ipv4_internal` AS `ipv4_internal`,`minions`.`ipv4_external` AS `ipv4_external`,`minions`.`online` AS `online`,`minions`.`created_at` AS `created_at`,`minions`.`updated_at` AS `updated_at`,`minions`.`deleted_at` AS `deleted_at` FROM `minions` LEFT JOIN `cluster_minion` ON `minions`.`id` = `cluster_minion`.`minion_id` WHERE `cluster_minion`.`cluster_id` = ? AND `minions`.`deleted_at` IS NULL"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(clusterId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	minions, err := scanMinions(rows)
	if err != nil {
		return nil, err
	}

	return minions, nil
}

func scanMinions(rows *sql.Rows) ([]*Minion, error) {
	var minions []*Minion
	for rows.Next() {
		var (
			id           int64
			name         sql.NullString
			ipv4Internal sql.NullString
			ipv4External sql.NullString
			online       int
			createdAt    sql.NullString
			updatedAt    sql.NullString
			deletedAt    sql.NullString
		)
		if err := rows.Scan(&id, &name, &ipv4Internal, &ipv4External, &online, &createdAt, &updatedAt, &deletedAt); err != nil {
			return nil, err
		}

		createdAt1, _ := parseTimestamp(createdAt)
		updatedAt1, _ := parseTimestamp(updatedAt)
		deletedAt1, _ := parseTimestamp(deletedAt)
		minion := &Minion{
			ID:           id,
			Name:         name.String,
			Ipv4Internal: ipv4Internal.String,
			Ipv4External: ipv4External.String,
			Online:       online,
			CreatedAt:    createdAt1,
			UpdatedAt:    updatedAt1,
			DeletedAt:    deletedAt1,
		}
		minions = append(minions, minion)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return minions, nil
}
