package dbops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"src/api/defs"
	"src/api/utils"
	"strconv"
	"time"
)

/* ------------ User API ------------ */

// AddUserCredential 创建用户
func AddUserCredential(loginName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("insert into users (login_name, pwd) values (?, ?)")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(loginName, pwd)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

// GetUserCredential 获取用户信息
func GetUserCredential(loginName string) (string, error) {
	stmtOut, err := dbConn.Prepare("select pwd from users where login_name = ?")
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	defer stmtOut.Close()

	return pwd, nil
}

func GetUserId(loginName string) (string, error) {
	stmtOut, err := dbConn.Prepare("select user_id from users where login_name = ?")
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	var id string
	err = stmtOut.QueryRow(loginName).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	defer stmtOut.Close()

	return id, nil
}

// DeleteUser 删除用户
func DeleteUser(loginName string, pwd string) error {
	stmtDel, err := dbConn.Prepare("delete from users where login_name = ? and pwd = ?")
	if err != nil {
		log.Printf("Delete User Failed: %s", err)
		return err
	}

	_, err = stmtDel.Exec(loginName, pwd)
	if err != nil {
		return err
	}

	defer stmtDel.Close()

	return nil
}

/* ------------ Video API ------------ */

// AddNewVideo 新建视频资源
func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {
	// 创建uuid
	vid, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}

	t := time.Now()
	// 模板数据不可更改，得到形如 M D y, HH:MM:SS的时间
	ctime := t.Format("Jan 02 2006, 15:04:05")
	stmtIns, err := dbConn.Prepare(`insert into video_info 
    	(id, author_id, name, display_ctime) values (?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}

	_, err = stmtIns.Exec(vid, aid, name, ctime)
	if err != nil {
		return nil, err
	}

	res := &defs.VideoInfo{
		Id:           vid,
		AuthorId:     aid,
		Name:         name,
		DisplayCtime: ctime,
	}
	defer stmtIns.Close()
	return res, nil
}

// GetVideo 根据vid获取视频资源
func GetVideo(vid string) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare("select author_id, name, display_ctime from video_info where id = ?")
	if err != nil {
		return nil, err
	}

	var aid int
	var name string
	var dct string
	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)

	// 出现错误
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// 查询结果为空
	if err == sql.ErrNoRows {
		return nil, nil
	}

	// 查询结果不为空
	res := &defs.VideoInfo{
		Id:           vid,
		AuthorId:     aid,
		Name:         name,
		DisplayCtime: dct,
	}

	defer stmtOut.Close()

	return res, nil
}

func ListVideo(username string, from, to int) ([]*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`select video_info.id, video_info.author_id, video_info.display_ctime
		from video_info inner join users using users.user_id = video_info.author_id
		where users.login_name = ? and video_info.create_time > FROM_UNIXTIME(?) and video_info.create_time <= FROM_UNIXTIME(?)
		order by video_info.create_time desc `)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	var res []*defs.VideoInfo

	rows, err := stmtOut.Query(username, from, to)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, aid, name, ctime string
		if err := rows.Scan(&id, &aid, &name, &ctime); err != nil {
			return res, err
		}

		aidInt, _ := strconv.Atoi(aid)
		c := &defs.VideoInfo{
			Id:           id,
			AuthorId:     aidInt,
			Name:         name,
			DisplayCtime: ctime,
		}
		res = append(res, c)
	}

	defer stmtOut.Close()
	return res, nil
}

// DeleteVideo 根据vid删除视频资源
func DeleteVideo(vid string) error {
	stmtDel, err := dbConn.Prepare("delete from video_info where id = ?")
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}

	defer stmtDel.Close()
	return nil
}

/* ------------ Comments API ------------ */

// AddNewComment 添加评论
func AddNewComment(aid int, vid string, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		return err
	}

	stmtIns, err := dbConn.Prepare(`insert into comments
    		(id, author_id, video_id, content) values (?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(id, aid, vid, content)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

func ListComments(vid string, from, to int) ([]*defs.Comment, error) {
	stmtOut, err := dbConn.Prepare(`select comments.id, user.login_name, comments.content
		from comments inner join user on comments.author_id = user.user_id
		where comments.video_id = ? and comments.time > FROM_UNIXTIME(?) and comments.time <= FROM_UNIXTIME(?)
		order by comments.time desc `)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	var res []*defs.Comment

	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil {
			return res, err
		}

		c := &defs.Comment{
			Id:         id,
			AuthorName: name,
			VideoId:    vid,
			Content:    content,
		}
		res = append(res, c)
	}

	defer stmtOut.Close()
	return res, nil
}
