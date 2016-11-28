package main

import (
	"database/sql"
	"github.com/Shakarang/Epirank/config"
	"github.com/Shakarang/Epirank/database"
	"github.com/Shakarang/Epirank/ranking"
	"github.com/Shakarang/Epirank/routes"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

// APIMiddleware will add the db connection to the context
func APIMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("database", db)
		c.Next()
	}
}

func init() {

	// Log as ASCII Formatter
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout, could also be a file.
	log.SetOutput(os.Stdout)

	// Log everything
	log.SetLevel(log.InfoLevel)
}

func main() {

	if db, err := database.Init(config.DatabasePath); err != nil {
		log.Fatal(err)
	} else {

		defer db.Close()
		database.CreateTable(db)

		if err := ranking.InitRanking(db); err != nil {
			log.Error(err)
			os.Exit(-1)
		}

		// Webservice
		router := gin.Default()

		router.Use(APIMiddleware(db))

		router.LoadHTMLGlob("templates/*")
		router.Static("/assets", "./assets")

		router.GET("/", routes.GetStudents)

		router.Run()
	}
}
