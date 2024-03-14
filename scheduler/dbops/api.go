package dbops

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func AddVideoDelRecord(vid string) error {
	stmtIns, err := dbConn.Prepare("insert into video_del_ref(video_id) values (?)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(vid)
	if err != nil {
		log.Printf("Error when Insert: %v", err)
		return err
	}

	defer stmtIns.Close()
	return nil
}
