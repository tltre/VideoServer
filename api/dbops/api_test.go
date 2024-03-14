package dbops

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

/*
	1. 清除表中数据
	2. 运行测试
	3. 观察结果
	4. 清除新增数据
*/

func clearTables() {
	dbConn.Exec("truncate users")
	dbConn.Exec("truncate video_info")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate session")
}

func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}

func TestUserWorkFlow(t *testing.T) {
	t.Run("Add", testAddUserCredential)
	t.Run("Get", testGetUserCredential)
	t.Run("Del", testDeleteUser)
	t.Run("ReGet", testReGetUserCredential)
}

var tempVid string

func TestVideoWorkFlow(t *testing.T) {
	clearTables()
	t.Run("Prepare user", testAddUserCredential)
	t.Run("Add Video", testAddNewVideo)
	t.Run("Find Video", testGetVideo)
	t.Run("Delete Video", testDeleteVideo)
	t.Run("ReGet Video", testReGetVideo)
	clearTables()
}

func TestComment(t *testing.T) {
	clearTables()
	t.Run("New User", testAddUserCredential)
	t.Run("New Video", testAddNewVideo)
	t.Run("Add Comment", testAddNewComment)
	t.Run("Find comments", testListComments)
	clearTables()
}

/* *** User测试函数 *** */

func testAddUserCredential(t *testing.T) {
	err := AddUserCredential("xy", "123")
	if err != nil {
		t.Errorf("Error of Add User: %v", err)
	}
}

func testGetUserCredential(t *testing.T) {
	pwd, err := GetUserCredential("xy")
	if pwd != "123" || err != nil {
		t.Errorf("Error of Get User")
	}
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser("xy", "123")
	if err != nil {
		t.Errorf("Error of Delete User: %v", err)
	}
}

func testReGetUserCredential(t *testing.T) {
	pwd, err := GetUserCredential("xy")
	if err != nil {
		t.Errorf("Error of Get User: %v", err)
	}
	if pwd != "" {
		t.Errorf("Deleting User Test Failed")
	}
}

/* *** Video测试函数 *** */

func testAddNewVideo(t *testing.T) {
	video, err := AddNewVideo(1, "FirstVideo")
	if err != nil {
		t.Errorf("Error of Add Video: %v", err)
	}

	tempVid = video.Id
}

func testGetVideo(t *testing.T) {
	_, err := GetVideo(tempVid)
	if err != nil {
		t.Errorf("Error of Get Video: %v", err)
	}
}

func testDeleteVideo(t *testing.T) {
	err := DeleteVideo(tempVid)
	if err != nil {
		t.Errorf("Error of Delete Video: %v", err)
	}
}

func testReGetVideo(t *testing.T) {
	video, err := GetVideo(tempVid)
	if err != nil || video != nil {
		t.Errorf("Error of ReGet Video")
	}
}

/* *** Comments测试函数 *** */

func testAddNewComment(t *testing.T) {
	err := AddNewComment(1, tempVid, "HelloWorld")
	if err != nil {
		t.Errorf("Error of Add Comment: %v", err)
	}
}

func testListComments(t *testing.T) {
	from := 1514764800
	to, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/100000000, 10))
	comments, err := ListComments(tempVid, from, to)
	if err != nil {
		t.Errorf("Error of List Comment: %v", err)
	}

	for i, ele := range comments {
		fmt.Printf("comments: %d %v\n", i, ele)
	}
}
