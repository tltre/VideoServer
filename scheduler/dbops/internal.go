package dbops

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

const VIDEO_PATH = ".\\videos\\"

func DeleteVideo(vid string) error {
	err := os.Remove(VIDEO_PATH + vid + ".mp4")
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Delete Video Error: %v", err)
		return err
	}
	return nil
}

func ReadVideoDeletionRecord(count int) ([]string, error) {
	stmtOut, err := dbConn.Prepare("select video_id from video_del_ref limit ?")
	var ids []string
	if err != nil {
		return ids, err
	}

	res, err := stmtOut.Query(count)
	if err != nil {
		log.Printf("Query VideoDeletionRecord Error: %v", err)
		return ids, err
	}

	for res.Next() {
		var id string
		if err := res.Scan(&id); err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}

	defer stmtOut.Close()
	return ids, nil
}

func DelVideoRecord(vid string) error {
	stmtDel, err := dbConn.Prepare("delete from video_del_ref where video_id = ?")
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		log.Printf("Error when Delete Record: %v", err)
		return err
	}

	defer stmtDel.Close()
	return nil
}
