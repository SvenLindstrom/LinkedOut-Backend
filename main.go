package main

import (
	"linkedout/databases"
	"linkedout/services/auth"
	"linkedout/services/location"
	"linkedout/services/user"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func loadConf(args []string) {
	_, set := os.LookupEnv("MODE_PROD")
	if !set {
		println("loading dev env")
		err := godotenv.Load(".dev.env")
		if err != nil {
			log.Fatal("failed to load .env file")
		}

		if len(args) > 0 {
			println("loading oauth conf")
			err := godotenv.Load(".oauth.env")
			if err != nil {
				log.Fatal("failed to load .env file")
			}
		}

	}
}

func main() {

	args := os.Args[1:]

	loadConf(args)

	PORT := ":3113"

	r := gin.Default()

	pg := databases.Pg_init()

	authGroup := r.Group("/auth")
	auth.Routes(authGroup, pg)

	api := r.Group("/api")
	api.Use(auth.TokenMiddleware())
	location.Routes(api, pg)
	user.Routes(api, pg)

	r.Run(PORT)
}
