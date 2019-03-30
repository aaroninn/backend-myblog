package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"hypermedlab/backend-myblog/pkgs/middlewares"
	"hypermedlab/backend-myblog/routes"
)

const secret = middlewares.Secret

func main() {
	conn, err := sqlx.Open("sqlite3", "./data.db")
	if err != nil {
		panic(err.Error())
	}

	gin.SetMode(gin.ReleaseMode)
	g := gin.Default()
	g.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	newUserRouter(conn, g)
	newBlogRouter(conn, g)

	err = g.Run("0.0.0.0:5000")
	if err != nil {
		panic(err.Error())
	}

}

func newUserRouter(conn *sqlx.DB, g *gin.Engine) {
	userRoute := routes.NewUserRouter(conn, g, secret)
	userRoute.Init()
}

func newBlogRouter(conn *sqlx.DB, g *gin.Engine) {
	blogRoute := routes.NewBlogRouter(conn, g)
	blogRoute.Init()
}
