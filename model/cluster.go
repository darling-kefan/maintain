package model

import (
	"fmt"
	"time"
	"strings"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// clusters机器组表
//
// 组名称
// 在线状态 0-在线; 1-下线;
type Cluster struct {
	ID        int64
	Name      string
	Online    int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	Minions   []*Minion
}

func (db *Mtdb) AddCluster(c *Cluster) (id int64, err error) {
	query := "INSERT INTO `clusters`(`name`, `online`, `created_at`) values(?, ?, ?)"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	r, err := execAffectingRows(stmt, c.Name, c.Online, time.Now().Format(TIME_FORMAT))
	if err != nil {
		return 0, err
	}
	lastInsertId, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
	}
	return lastInsertId, nil
}

func (db *Mtdb) UpdateCluster(c *Cluster) error {
	if c.ID == 0 {
		return fmt.Errorf("mysql: Cluster with unassigned ID passed into UpdateCluster")
	}

	query := "UPDATE `clusters` SET `name`=?, `online`=? WHERE `id` = ?"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = execAffectingRows(stmt, c.Name, c.Online, c.ID)
	return err
}

func (db *Mtdb) DeleteCluster(id int64) error {
	if id == 0 {
		return fmt.Errorf("mysql: unassigned id passed into DeleteCluster")
	}

	// query := "DELETE FROM `clusters` WHERE `id` = ?"
	query := "UPDATE `clusters` SET `deleted_at` = ? WHERE `id` = ?"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = execAffectingRows(stmt, time.Now().Format(TIME_FORMAT), id)
	return err
}

// mark: 0-不限; 1-在线; 2-下线
func (db *Mtdb) GetClusters(mark int) ([]*Cluster, error) {
	var args []interface{}
	query := "SELECT * FROM `clusters` WHERE `deleted_at` IS NULL"
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

	clusters, err := scanClusters(rows)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}

func (db *Mtdb) GetClustersByIds(ids []int64) ([]*Cluster, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("The param ids: %v is empty", ids)
	}
	query := "SELECT * FROM `clusters` WHERE `deleted_at` IS NULL AND `id` IN (?"+strings.Repeat(",?", len(ids)-1)+")"
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

	clusters, err := scanClusters(rows)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}

func scanClusters(rows *sql.Rows) ([]*Cluster, error) {
	var clusters []*Cluster
	for rows.Next() {
		var (
			id           int64
			name         sql.NullString
			online       int
			createdAt    sql.NullString
			updatedAt    sql.NullString
			deletedAt    sql.NullString
		)
		if err := rows.Scan(&id, &name, &online, &createdAt, &updatedAt, &deletedAt); err != nil {
			return nil, err
		}

		createdAt1, _ := parseTimestamp(createdAt)
		updatedAt1, _ := parseTimestamp(updatedAt)
		deletedAt1, _ := parseTimestamp(deletedAt)
		cluster := &Cluster{
			ID:        id,
			Name:      name.String,
			Online:    online,
			CreatedAt: createdAt1,
			UpdatedAt: updatedAt1,
			DeletedAt: deletedAt1,
		}
		clusters = append(clusters, cluster)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clusters, nil
}





func (db *Mtdb) AddClusterMinion(clusterId, minionId int64) (id int64, err error) {
	query := "INSERT INTO `cluster_minion`(`cluster_id`, `minion_id`) values(?, ?)"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	r, err := execAffectingRows(stmt, clusterId, minionId)
	if err != nil {
		return
	}
	id, err = r.LastInsertId()
	if err != nil {
		return
	}
	return
}

func (db *Mtdb) DeleteClusterMinions(clusterId int64) error {
	if clusterId == 0 {
		return fmt.Errorf("mysql: unassigned clusterId passed into DeleteClusterMinions")
	}

	query := "DELETE FROM `cluster_minion` WHERE `cluster_id` = ?"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = execAffectingRows(stmt, clusterId)
	return err
}
