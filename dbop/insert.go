package dbop

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//"strconv"
	. "wbproject/chufangrefresh/structure"
	. "wbproject/chufangrefresh/util"
)

var db1, db2 *sql.DB
var key string

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
		fmt.Println("gettype is ", rf.Type)
		//aid操作
		err := InsertAid(rf.Id, rf.St, rf.Et)
		return err

	case 1:
		fmt.Println("gettype is ", rf.Type)
		//todo .. gid 操作
		err := InsertGid(rf.Id, rf.St, rf.Et)
		return err

	case 2:
		fmt.Println("gettype is ", rf.Type)
		//todo .. uid 操作
		err := InsertUid(rf.Id, rf.St, rf.Et)
		return err

	default:
		fmt.Println(rf.Type, " is of a type I don't know how to handle")

	}
	return nil
}

func Process(db *sql.DB, uid int, date int64) (error, bool) {

	var uploadid int
	row := db.QueryRow("select uploadid from hmp_data_eventqueue where walkdate  = ? and userid= ?", date, uid)
	err := row.Scan(&uploadid)

	if err != nil {
		return err, false
	}

	fmt.Println("in SelectUploadid ", uploadid)

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

func InsertUid(uid int, st int64, et int64) error {

	var arr_userinfo []userinfo

	count := int((et - st) / 86400)

	for i := 0; i <= count; i++ {

		arr_userinfo = append(
			arr_userinfo,
			userinfo{
				userid: uid,
				date:   st + int64(i*86400),
			})
	}

	fmt.Println("array len of userinfo is ", len(arr_userinfo))

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

	Process(db2, arr_userinfo[len(arr_userinfo)-1].userid, arr_userinfo[len(arr_userinfo)-1].date)

	return nil
}

func InsertAid(aid int, st int64, et int64) error {

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

		for i := 0; i <= count; i++ {

			arr_userinfo = append(
				arr_userinfo,
				userinfo{
					userid: uid,
					date:   st + int64(i*86400),
				})
		}
	}

	fmt.Println("array len of userinfo is ", len(arr_userinfo))

	sqlStr := "insert into hmp.hmp_data_eventqueue(sourcetable,userid, walkdate, activeid,timestamp) VALUES "

	vals := []interface{}{}

	for _, uinfo := range arr_userinfo {

		sqlStr += "('hmp_walking_tasks_000',?,?,-1,UNIX_TIMESTAMP()),"
		vals = append(vals, uinfo.userid, uinfo.date)
	}

	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	//format all vals at once
	_, err = db2.Exec(sqlStr, vals...)

	if err != nil {
		return err
	}

	Process(db2, arr_userinfo[len(arr_userinfo)-1].userid, arr_userinfo[len(arr_userinfo)-1].date)

	return nil
}

func InsertGid(gid int, st int64, et int64) error {

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
			fmt.Println("err xx is ", err)
			return err
		}

		count := int((et - st) / 86400)

		for i := 0; i <= count; i++ {

			arr_userinfo = append(
				arr_userinfo,
				userinfo{
					userid: uid,
					date:   st + int64(i*86400),
				})
		}
	}

	fmt.Println("array len of userinfo is ", len(arr_userinfo))

	sqlStr := "insert into hmp.hmp_data_eventqueue(sourcetable,userid, walkdate, activeid,timestamp) VALUES "

	vals := []interface{}{}

	for _, uinfo := range arr_userinfo {

		sqlStr += "('hmp_walking_tasks_000',?,?,-1,UNIX_TIMESTAMP()),"
		vals = append(vals, uinfo.userid, uinfo.date)
	}

	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	//format all vals at once
	_, err = db2.Exec(sqlStr, vals...)

	if err != nil {
		return err
	}

	Process(db2, arr_userinfo[len(arr_userinfo)-1].userid, arr_userinfo[len(arr_userinfo)-1].date)

	return nil
}
