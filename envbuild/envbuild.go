package envbuild

import (
	"database/sql"
	"flag"
	//"fmt"
	_ "github.com/go-sql-driver/mysql"
	config "github.com/msbranco/goconfig"
)

var db1, db2 *sql.DB
var config_file_path string

func init() {
	flag.StringVar(&config_file_path, "c", "config file", "Use -c <filepath>")
}

//EnvBuild需要正确的解析文件并且初始化DB和Redis的连接。。
func EnvBuild() (*sql.DB, *sql.DB, error) {

	//get conf
	cf, err := config.ReadConfigFile(config_file_path)

	if err != nil {
		return nil, nil, err
	}

	rdip1, _ := cf.GetString("DBCONN1", "IP")
	rdusr1, _ := cf.GetString("DBCONN1", "USERID")
	rdpwd1, _ := cf.GetString("DBCONN1", "USERPWD")
	rdname1, _ := cf.GetString("DBCONN1", "DBNAME")

	rdip1 = rdusr1 + ":" + rdpwd1 + "@tcp(" + rdip1 + ")/" + rdname1 + "?charset=utf8"

	rdip2, _ := cf.GetString("DBCONN2", "IP")
	rdusr2, _ := cf.GetString("DBCONN2", "USERID")
	rdpwd2, _ := cf.GetString("DBCONN2", "USERPWD")
	rdname2, _ := cf.GetString("DBCONN2", "DBNAME")

	rdip2 = rdusr2 + ":" + rdpwd2 + "@tcp(" + rdip2 + ")/" + rdname2 + "?charset=utf8"

	//open db1
	db1, _ = sql.Open("mysql", rdip1)
	//defer db1.Close()
	db1.SetMaxOpenConns(100)
	db1.SetMaxIdleConns(10)
	db1.Ping()

	//open db2
	db2, _ = sql.Open("mysql", rdip2)
	//defer db1.Close()
	db2.SetMaxOpenConns(100)
	db2.SetMaxIdleConns(10)
	db2.Ping()

	return db1, db2, nil
}
