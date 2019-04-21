package main

import (
	blogDB "hypermedlab/backend-myblog/models/blog/db"
	userDB "hypermedlab/backend-myblog/models/user/db"
	"hypermedlab/backend-myblog/pkgs/middlewares"
	"hypermedlab/backend-myblog/pkgs/session"
	"hypermedlab/backend-myblog/routes"
	"hypermedlab/backend-myblog/services/blog"
	"hypermedlab/backend-myblog/services/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/dig"
)

const secret = middlewares.Secret

func buildContainer() *dig.Container {
	container := dig.New()
	container.Provide(routes.NewBlogRouter)
	container.Provide(routes.NewUserRouter)
	container.Provide(blogDB.NewSqlite3)
	container.Provide(userDB.NewSqlite3)
	container.Provide(blog.NewBlogService)
	container.Provide(user.NewService)
	container.Provide(session.NewSessionsStorage)
	container.Provide(middlewares.NewMiddleWare)
	container.Provide(func() (*sqlx.DB, *gin.Engine, string) {
		gin.SetMode(gin.ReleaseMode)
		engine := gin.Default()
		engine.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"*"},
			AllowCredentials: true,
		}))

		conn, err := sqlx.Open("sqlite3", "./data.db")
		if err != nil {
			panic(err.Error())
		}

		return conn, engine, middlewares.Secret

	})

	return container

}

func main() {
	container := buildContainer()
	container.Invoke(func(userroute *routes.User, blogroute *routes.Blog) {
		userroute.Init()
		blogroute.Init()
		userroute.Engine.Run(":8080")
	})

}
