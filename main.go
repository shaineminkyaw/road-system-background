package main

import (
	"github.com/shaineminkyaw/road-system-background/config"
	"github.com/shaineminkyaw/road-system-background/ds"
)

func main() {
	//

	conf := config.Init()
	db := ds.ConnectToDB(conf.SQL.Host, conf.SQL.Port, conf.SQL.DB, conf.SQL.User, conf.SQL.Password)
	ds.NewRBAC(db)

}
