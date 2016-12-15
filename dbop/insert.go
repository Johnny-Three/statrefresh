package dbop

import (
	"database/sql"
	"fmt"

	"errors"
	. "wbproject/chufangrefresh/logs"
	. "wbproject/chufangrefresh/structure"
	. "wbproject/chufangrefresh/util"

	_ "github.com/go-sql-driver/mysql"
)

var db1, db2 *sql.DB
var key string
var def = 2000

type userinfo struct {
	userid int
	date   int64
}

var Deal_status_map *BeeMap

func init() {

	Deal_status_map = NewBeeMap()
}

func InsertQueue(rf *Refresh, key01 string, db01, db02 *sql.DB) error {

	key, db1, db2 = key01, db01, db02

	switch rf.Type {

	case 0:
		Logger.Info("gettype is ", rf.Type)
		//aid操作
		err := InsertAid(rf.Id, rf.St, rf.Et, key)
		return err

	case 1:
		Logger.Info("gettype is ", rf.Type)
		//todo .. gid 操作
		err := InsertGid(rf.Id, rf.St, rf.Et, key)
		return err

	case 2:
		Logger.Info("gettype is ", rf.Type)
		//todo .. uid 操作
		err := InsertUid(rf.Id, rf.St, rf.Et, key)
		return err

	default:
		Logger.Info(rf.Type, " is of a type I don't know how to handle")

	}
	return nil
}

func Youmeiyouren(db *sql.DB, rf *Refresh) bool {

	var count int
	var row *sql.Row

	switch rf.Type {

	case 0:
		row = db.QueryRow("select count(*) from wanbu_group_user where activeid  = ? ", rf.Id)
		break

	case 1:
		row = db.QueryRow("select count(*) from wanbu_group_user where groupid  = ? ", rf.Id)
		break

	case 2:
		row = db.QueryRow("select count(*) from wanbu_group_user where userid  = ? ", rf.Id)
		break

	default:
		return false
	}

	err := row.Scan(&count)

	if err != nil {
		fmt.Println("error when check wanbu_group_user", err.Error())
		return false
	}

	if count > 0 {
		return true
	}
	return false
}

func Process(db *sql.DB, uid int, date int64, key string) (error, bool) {

	var uploadid int
	row := db.QueryRow("select uploadid from hmp_data_eventqueue where walkdate  = ? and userid= ?", date, uid)
	err := row.Scan(&uploadid)

	if err != nil {
		return err, false
	}

	Logger.Info("in SelectUploadid ", uploadid)
	Logger.Info("deal_status_map set ", key, uploadid)
	Deal_status_map.Set(key, uploadid)

	return nil, true
}

func SelectUploadid(db *sql.DB, id int) bool {

	var count int

	row := db.QueryRow("select count(*) from hmp_data_eventqueue where uploadid  = ? ", id)
	err := row.Scan(&count)

	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func InsertUid(uid int, st int64, et int64, key string) error {

	var arr_userinfo []userinfo

	count := int((et - st) / 86400)

	if count > 730 {

		return errors.New("刷新数据量超过了两年，不予刷新")
	}

	for i := 0; i <= count; i++ {

		arr_userinfo = append(
			arr_userinfo,
			userinfo{
				userid: uid,
				date:   st + int64(i*86400),
			})
	}

	Logger.Infof("用户记录总数为【%d】", len(arr_userinfo))

	sqlStr := "insert into hmp.hmp_data_eventqueue(sourcetable,userid, walkdate, activeid,timestamp) VALUES "

	vals := []interface{}{}

	for _, uinfo := range arr_userinfo {

		sqlStr += "('hmp_walking_tasks_000',?,?,-1,UNIX_TIMESTAMP()),"
		vals = append(vals, uinfo.userid, uinfo.date)
	}

	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	//format all vals at once
	_, err := db2.Exec(sqlStr, vals...)

	if err != nil {
		return err
	}

	Process(db2, arr_userinfo[len(arr_userinfo)-1].userid, arr_userinfo[len(arr_userinfo)-1].date, key)

	return nil
}

func InsertAid(aid int, st int64, et int64, key string) error {

	var uid int

	qs := `select userid from wanbu_group_user where activeid  = ?`

	//fmt.Println("in insertgid", db1)
	rows, err := db1.Query(qs, aid)
	if err != nil {
		return err
	}
	defer rows.Close()
	var arr_userinfo []userinfo

	for rows.Next() {

		err := rows.Scan(&uid)

		if err != nil {
			fmt.Println("err xx is ", err)
			return err
		}

		count := int((et - st) / 86400)

		if count > 730 {

			return errors.New("刷新数据量超过了两年，不予刷新")

		}

		for i := 0; i <= count; i++ {

			arr_userinfo = append(
				arr_userinfo,
				userinfo{
					userid: uid,
					date:   st + int64(i*86400),
				})
		}
	}

	Logger.Infof("用户记录总数为【%d", len(arr_userinfo))

	stepth := len(arr_userinfo) / def
	Logger.Infof("分【%d】次插入hmp_data_eventqueue表，每次%d条", stepth, def)

	for i := 0; i < stepth; i++ {

		sqlStr := "insert into hmp.hmp_data_eventqueue(sourcetable,userid, walkdate, activeid,timestamp) VALUES "
		vals := []interface{}{}

		for j := i * def; j < (i+1)*def; j++ {

			sqlStr += "('hmp_walking_tasks_000',?,?,-1,UNIX_TIMESTAMP()),"
			vals = append(vals, arr_userinfo[j].userid, arr_userinfo[j].date)

		}
		//trim the last ,
		sqlStr = sqlStr[0 : len(sqlStr)-1]
		//format all vals at once
		_, err = db2.Exec(sqlStr, vals...)

		if err != nil {
			return err
		}
		Logger.Infof("总[%d]条数据,总[%d]批,第[%d]批处理完毕,此批[%d]记录", len(arr_userinfo), stepth, i, def)
	}

	yu := len(arr_userinfo) % def

	//模除部分处理
	if yu != 0 {

		sqlStr := "insert into hmp.hmp_data_eventqueue(sourcetable,userid, walkdate, activeid,timestamp) VALUES "
		vals := []interface{}{}

		for j := stepth * def; j < len(arr_userinfo); j++ {

			sqlStr += "('hmp_walking_tasks_000',?,?,-1,UNIX_TIMESTAMP()),"
			vals = append(vals, arr_userinfo[j].userid, arr_userinfo[j].date)
		}

		//trim the last ,
		sqlStr = sqlStr[0 : len(sqlStr)-1]
		//format all vals at once
		_, err = db2.Exec(sqlStr, vals...)

		if err != nil {
			return err
		}

		Logger.Infof("总[%d]条数据,总[%d]批,第[%d]批处理完毕,此批[%d]记录", len(arr_userinfo), stepth, stepth,
			len(arr_userinfo[stepth*def:]))
	}

	Process(db2, arr_userinfo[len(arr_userinfo)-1].userid, arr_userinfo[len(arr_userinfo)-1].date, key)

	Logger.Infof("没看错，他走到了这里，这意味着一个Refresh请求成功写入了DB")

	return nil
}

func InsertGid(gid int, st int64, et int64, key string) error {

	var uid int

	qs := `select userid from wanbu_group_user where groupid  = ?`

	//fmt.Println("in insertgid", db1)
	rows, err := db1.Query(qs, gid)
	if err != nil {
		return err
	}
	defer rows.Close()
	var arr_userinfo []userinfo
	for rows.Next() {

		err := rows.Scan(&uid)

		if err != nil {
			Logger.Critical("err xx is ", err)
			return err
		}

		count := int((et - st) / 86400)

		if count > 730 {

			return errors.New("刷新数据量超过了两年，不予刷新")

		}

		for i := 0; i <= count; i++ {

			arr_userinfo = append(
				arr_userinfo,
				userinfo{
					userid: uid,
					date:   st + int64(i*86400),
				})
		}
	}

	Logger.Infof("用户记录总数为【%d】", len(arr_userinfo))

	stepth := len(arr_userinfo) / def
	Logger.Infof("分【%d】次插入hmp_data_eventqueue表，每次%d条", stepth, def)

	for i := 0; i < stepth; i++ {

		sqlStr := "insert into hmp.hmp_data_eventqueue(sourcetable,userid, walkdate, activeid,timestamp) VALUES "
		vals := []interface{}{}

		for j := i * def; j < (i+1)*def; j++ {

			sqlStr += "('hmp_walking_tasks_000',?,?,-1,UNIX_TIMESTAMP()),"
			vals = append(vals, arr_userinfo[j].userid, arr_userinfo[j].date)

		}
		//trim the last ,
		sqlStr = sqlStr[0 : len(sqlStr)-1]
		//format all vals at once
		_, err = db2.Exec(sqlStr, vals...)

		if err != nil {
			return err
		}
		Logger.Infof("总[%d]条数据,总[%d]批,第[%d]批处理完毕,此批[%d]记录", len(arr_userinfo), stepth, i, def)
	}

	yu := len(arr_userinfo) % def

	//模除部分处理
	if yu != 0 {

		sqlStr := "insert into hmp.hmp_data_eventqueue(sourcetable,userid, walkdate, activeid,timestamp) VALUES "
		vals := []interface{}{}

		for j := stepth * def; j < len(arr_userinfo); j++ {

			sqlStr += "('hmp_walking_tasks_000',?,?,-1,UNIX_TIMESTAMP()),"
			vals = append(vals, arr_userinfo[j].userid, arr_userinfo[j].date)
		}

		//trim the last ,
		sqlStr = sqlStr[0 : len(sqlStr)-1]
		//format all vals at once
		_, err = db2.Exec(sqlStr, vals...)

		if err != nil {
			return err
		}

		Logger.Infof("总[%d]条数据,总[%d]批,第[%d]批处理完毕,此批[%d]记录", len(arr_userinfo), stepth, stepth,
			len(arr_userinfo[stepth*def:]))
	}

	Process(db2, arr_userinfo[len(arr_userinfo)-1].userid, arr_userinfo[len(arr_userinfo)-1].date, key)

	Logger.Infof("没看错，他走到了这里，这意味着一个Refresh请求成功写入了DB")

	return nil
}
