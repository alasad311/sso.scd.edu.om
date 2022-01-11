package main

import (
	"sso.scd.edu.om/module"
	route "sso.scd.edu.om/routes"
)

func main() {
	//DB connection
	db, dbError := module.DBConnectionSSO()
	if dbError != nil {
		panic(dbError)
	}
	route.SetupRoutes(db)
}
