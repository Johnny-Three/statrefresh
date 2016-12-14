package main

import (
	"flag"
	"fmt"
	"os"
	. "wbproject/chufangrefresh/envbuild"
	. "wbproject/chufangrefresh/server"
)

var version string = "1.0.0PR6"

func main() {

	args := os.Args

	if len(args) == 2 && (args[1] == "-v") {

		fmt.Println("看好了兄弟，现在的版本是【", version, "】，可别弄错了")
		os.Exit(0)
	}

	flag.Parse()

	db1, db2, ip, err := EnvBuild()
	if err != nil {
		panic(err.Error())
	}

	go CheckAndDeleteMap()

	WebServerBase(db1, db2, ip)
}
