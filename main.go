package main

import (
	"linkedout/databases"
	"linkedout/services/auth"
	"linkedout/services/location"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func loadConf() {
	_, set := os.LookupEnv("MODE_PROD")
	if !set {
		println("loading dev env")
		err := godotenv.Load(".dev.env")
		if err != nil {
			log.Fatal("failed to load .env file")
		}
	}
}

func main() {

	loadConf()

	PORT := ":3113"

	r := gin.Default()

	pg := databases.Pg_init()

	authGroup := r.Group("/auth")

	auth.Routes(authGroup, pg)

	api := r.Group("/api")

	api.Use(auth.TokenMiddleware())

	location.Routes(api, pg)

	r.Run(PORT)
}
