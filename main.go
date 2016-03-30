package main

import (
	"flag"
	. "wbproject/chufangrefresh/envbuild"
	. "wbproject/chufangrefresh/server"
)

func main() {

	flag.Parse()

	db1, db2, ip, err := EnvBuild()
	if err != nil {
		panic(err.Error())
	}

	WebServerBase(db1, db2, ip)
}
