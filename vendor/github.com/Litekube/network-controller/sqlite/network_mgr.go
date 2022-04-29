package sqlite

import (
	"errors"
	"time"
)

type NetworkMgr struct {
	Id         int64
	Token      string
	State      int
	BindIp     string
	CreateTime time.Time
	UpdateTime time.Time
}

func (network *NetworkMgr) Insert(u NetworkMgr) error {
	db = GetDb()
	sql := `insert into network_mgr (token, state, bind_ip) values(?,?,?)`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.Token, u.State, u.BindIp)
	return err
}

func (network *NetworkMgr) InsertToken(token string) error {
	db = GetDb()
	sql := `insert into network_mgr (token) values()`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(token)
	return err
}

func (network *NetworkMgr) QueryAll() (bindIps []string, e error) {
	db = GetDb()
	sql := `select bind_ip from network_mgr`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var result = make([]string, 0)
	for rows.Next() {
		var bindIp string
		rows.Scan(&bindIp)
		result = append(result, bindIp)
	}
	return result, nil
}

func (network *NetworkMgr) QueryByToken(token string) (l *NetworkMgr, e error) {
	db = GetDb()
	sql := `select * from network_mgr where token=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(token)
	if err != nil {
		return nil, err
	}
	var result = make([]NetworkMgr, 0)
	for rows.Next() {
		var token, bindIp string
		var id int64
		var state int
		var createTime, updateTime time.Time
		rows.Scan(&id, &token, &state, &bindIp, &createTime, &updateTime)
		result = append(result, NetworkMgr{id, token, state, bindIp, createTime, updateTime})
	}
	if len(result) == 0 {
		return nil, errors.New("fail to find such item")
	}
	return &result[0], nil
}

func (network *NetworkMgr) QueryByIp(ip string) (l *NetworkMgr, e error) {
	db = GetDb()
	sql := `select * from network_mgr where bind_ip=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(ip)
	if err != nil {
		return nil, err
	}
	var result = make([]NetworkMgr, 0)
	for rows.Next() {
		var token, bindIp string
		var id int64
		var state int
		var createTime, updateTime time.Time
		rows.Scan(&id, &token, &state, &bindIp, &createTime, &updateTime)
		result = append(result, NetworkMgr{id, token, state, bindIp, createTime, updateTime})
	}
	if len(result) == 0 {
		return nil, errors.New("fail to find such item")
	}
	return &result[0], nil
}

func (network *NetworkMgr) QueryLogestIdle() (l *NetworkMgr, e error) {
	db = GetDb()
	// valid in sqlite
	sql := `select id,token,state,bind_ip,create_time,update_time,min(update_time) from network_mgr where state=-1`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var result = make([]NetworkMgr, 0)
	for rows.Next() {
		var token, bindIp string
		var id int64
		var state int
		var createTime, updateTime, tmp time.Time
		rows.Scan(&id, &token, &state, &bindIp, &createTime, &updateTime, &tmp)
		result = append(result, NetworkMgr{id, token, state, bindIp, createTime, updateTime})
		//rows.Scan(&id, &token, &bindIp, &updateTime)
		//result = append(result, NetworkMgr{
		//	Id:         id,
		//	Token:      token,
		//	BindIp:     bindIp,
		//	UpdateTime: updateTime,
		//})
	}
	if len(result) == 0 {
		return nil, errors.New("fail to find such item")
	}
	return &result[0], nil
}

func (network *NetworkMgr) UpdateStateByToken(state int, token string) (bool, error) {
	db = GetDb()
	sql := `update network_mgr set state=? where token=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(state, token)
	if err != nil {
		return false, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (network *NetworkMgr) UpdateIpByToken(ip, token string) (bool, error) {
	db = GetDb()
	sql := `update network_mgr set bind_ip=? where token=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(ip, token)
	if err != nil {
		return false, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (network *NetworkMgr) UpdateAllState() (bool, error) {
	db = GetDb()
	// add state!=-1 for no change update_time
	sql := `update network_mgr set state=-1 where state!=-1`
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

func (network *NetworkMgr) DeleteById(id int64) (bool, error) {
	db = GetDb()
	sql := `delete from network_mgr where id=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(id)
	if err != nil {
		return false, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (network *NetworkMgr) DeleteByToken(token string) (bool, error) {
	db = GetDb()
	sql := `delete from network_mgr where token=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(token)
	if err != nil {
		return false, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (network *NetworkMgr) DeleteUnRegisteredIdle(expire int) (bool, error) {
	db = GetDb()
	// valid in sqlite
	sql := `delete from network_mgr where state=-1 and bind_ip="" and julianday('now','localtime')*1440 -julianday(create_time)*1440>? and token!="reserverd"`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(expire)
	if err != nil {
		return false, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}
	return true, nil
}
