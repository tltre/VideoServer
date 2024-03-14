package defs

/* ------- request ------- */

type UserCredential struct {
	// 打上json Tag，方便后续处理
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}

type NewVideoInfo struct {
	AuthorId int    `json:"author_id"`
	Name     string `json:"name"`
}

type NewComment struct {
	AuthorId int    `json:"author_id"`
	Content  string `json:"content"`
}

/* ------- response ------- */

type SignedUp struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

type SignedIn struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

type UserInfo struct {
	Username string `json:"username"`
	UserId   string `json:"user_id"`
}

type UserSession struct {
	UserName  string `json:"user_name"`
	SessionId string `json:"session_id"`
}

type VideosInfo struct {
	Videos []*VideoInfo `json:"videos"`
}

type CommentList struct {
	Comments []*Comment `json:"comments"`
}

/* ------- Data Model ------- */

type User struct {
	Id   int
	Pwd  string
	Name string
}

type VideoInfo struct {
	Id           string `json:"id"`
	AuthorId     int    `json:"author_id"`
	Name         string `json:"name"`
	DisplayCtime string `json:"display_ctime"`
}

type Comment struct {
	Id         string `json:"id"`
	AuthorName string `json:"author_name"`
	VideoId    string `json:"video_id,"`
	Content    string `json:"content"`
}

type Session struct {
	UserName string `json:"user_name"` // login name
	TTL      int64  `json:"ttl"`
}
