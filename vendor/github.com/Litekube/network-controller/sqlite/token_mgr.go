package sqlite

import (
	"errors"
	"fmt"
	"time"
)

type TokenMgr struct {
	Id         int64
	Token      string
	ExpireTime time.Time
	CreateTime time.Time
	UpdateTime time.Time
}

func (tm *TokenMgr) Insert(t TokenMgr, expireTime int32) error {
	db = GetDb()
	var sql string
	if expireTime < 0 {
		sql = fmt.Sprintf(`insert into token_mgr (token,expire_time) values(?,-1)`)
	} else {
		sql = fmt.Sprintf(`insert into token_mgr (token,expire_time) values(?,datetime(CURRENT_TIMESTAMP, 'localtime','+%d minute'))`, expireTime)
	}
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(t.Token)
	return err
}

func (tm *TokenMgr) QueryByToken(token string) (l *TokenMgr, e error) {
	db = GetDb()
	sql := `select * from token_mgr where token=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(token)
	if err != nil {
		return nil, err
	}
	var result = make([]TokenMgr, 0)
	for rows.Next() {
		var token string
		var id int64
		var expireTime, createTime, updateTime time.Time
		rows.Scan(&id, &token, &createTime, &updateTime)
		result = append(result, TokenMgr{id, token, expireTime, createTime, updateTime})
	}
	if len(result) == 0 {
		return nil, errors.New("fail to find such item")
	}
	return &result[0], nil
}

func (tm *TokenMgr) DeleteExpireToken() (bool, error) {
	db = GetDb()
	// valid in sqlite
	sql := `delete from token_mgr where julianday('now','localtime')*1440 >julianday(expire_time)*1440 and expire_time!=-1`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec()
	if err != nil {
		return false, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}
	return true, nil
}
