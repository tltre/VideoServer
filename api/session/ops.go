package session

import (
	"src/api/dbops"
	"src/api/defs"
	"src/api/utils"
	"sync"
	"time"
)

var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

func NowInMilli() int64 {
	return time.Now().UnixNano() / 100000
}

// DeleteExpiredSession 在内存中与数据库中均删除过期session
func DeleteExpiredSession(sid string) error {
	sessionMap.Delete(sid)
	err := dbops.DeleteSession(sid)
	if err != nil {
		return err
	}
	return nil
}

func LoadSessionsFromDB() {
	r, err := dbops.RetrieveAllSessions()
	if err != nil {
		return
	}

	// 将查询得到的session记录依次写入全局map中
	r.Range(func(k, v interface{}) bool {
		ss := v.(*defs.Session)
		sessionMap.Store(k, ss)
		return true
	})
}

func GenerateNewSessionId(username string) string {
	// 获取一个新的uuid
	id, _ := utils.NewUUID()
	// 使用毫秒
	ct := NowInMilli()
	// 设置有效时间为30分钟
	ttl := ct + 30*60*1000

	// 保存在内存中
	ss := &defs.Session{
		UserName: username,
		TTL:      ttl,
	}
	sessionMap.Store(id, ss)

	// 保存在数据库中
	dbops.InsertSession(id, ttl, username)

	return id
}

func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	if ok {
		ct := NowInMilli()
		if ss.(*defs.Session).TTL < ct {
			// 删除数据库原session数据
			DeleteExpiredSession(sid)
			return "", true
		}
		return ss.(*defs.Session).UserName, false
	}
	return "", true
}
