package dbops

import (
	"database/sql"
	"log"
	"src/api/defs"
	"strconv"
	"sync"
)

// InsertSession 添加session信息
func InsertSession(sid string, ttl int64, uname string) error {
	ttlstr := strconv.FormatInt(ttl, 10)
	stmtIns, err := dbConn.Prepare("insert into session (id, TTL, login_name) values (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(sid, ttlstr, uname)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

// RetrieveSession 从数据库中取出session
func RetrieveSession(sid string) (*defs.Session, error) {
	ss := &defs.Session{}
	stmtOut, err := dbConn.Prepare("select TTL, login_name from session where id = ?")
	if err != nil {
		return nil, err
	}

	var ttl, uname string
	err = stmtOut.QueryRow(sid).Scan(&ttl, &uname)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if res, err := strconv.ParseInt(ttl, 10, 64); err == nil {
		ss.TTL = res
		ss.UserName = uname
	} else {
		return nil, err
	}

	defer stmtOut.Close()
	return ss, nil
}

func RetrieveAllSessions() (*sync.Map, error) {
	m := &sync.Map{}
	stmtOut, err := dbConn.Prepare("select * from session")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	rows, err := stmtOut.Query()
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	for rows.Next() {
		var id string
		var ttlstr string
		var uname string
		if er := rows.Scan(&id, &ttlstr, &uname); er != nil {
			log.Printf("%s", er)
			break
		}

		if ttl, err1 := strconv.ParseInt(ttlstr, 10, 64); err1 == nil {
			ss := &defs.Session{
				UserName: uname,
				TTL:      ttl,
			}
			m.Store(id, ss)
			log.Printf("Session id: %s, ttl: %d, uname: %s", id, ss.TTL, ss.UserName)
		}
	}
	return m, nil
}

func DeleteSession(sid string) error {
	stmtDel, err := dbConn.Prepare("delete from session where id = ?")
	if err != nil {
		log.Printf("err: %s", err)
	}

	if _, err := stmtDel.Exec(sid); err != nil {
		return err
	}

	return nil
}
