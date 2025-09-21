package main

import (
	"ichibuy/auth/config"
	"ichibuy/auth/db"
	_ "ichibuy/auth/docs"
	"ichibuy/auth/server"
)

// @title           ichibuy/auth API
// @version         1.0
// @description     This is the ichibuy/auth API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      ichibuy-auth.vercel.app

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cfg := config.Load()
	db, err := db.New(cfg.PostgresURI)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := server.New(cfg, db)
	router.Run(":" + cfg.APIPort)
}
