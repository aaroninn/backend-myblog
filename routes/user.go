package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"hypermedlab/backend-myblog/pkgs/forms"
	"hypermedlab/backend-myblog/pkgs/jwt"
	"hypermedlab/backend-myblog/pkgs/middlewares"
	userSrv "hypermedlab/backend-myblog/services/user"
	"log"
)

var secret string

type User struct {
	srv    *userSrv.Service
	Engine *gin.Engine
}

func NewUserRouter(conn *sqlx.DB, engine *gin.Engine, se string) *User {
	secret = se
	return &User{
		srv:    userSrv.NewService(conn),
		Engine: engine,
	}

}

func (u *User) Init() {
	userroute := u.Engine.Group("/user")
	userroute.POST("/login", middlewares.IPCount, u.loginHandler)
	userroute.POST("/register", middlewares.IPCount, u.registerHandler)
	userroute.PUT("/password", middlewares.IPCount, middlewares.AuthToken, u.updatePasswordHandler)

	adminroute := u.Engine.Group("/admin")
	adminroute.GET("/users", middlewares.IPCount, middlewares.AdminAuthToken, u.findAllUsersHandler)
	// adminroute.PUT("/user/:id/status/:status", middlewares.IPCount, middlewares.AdminAuthToken, u.changUserStatusHandler)
}

func (u *User) registerHandler(ctx *gin.Context) {
	var form forms.CreateUser
	err := ctx.Bind(&form)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	usr, err := u.srv.RegisterUser(&form)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.JSON(200, usr)
}

func (u *User) loginHandler(ctx *gin.Context) {
	var form forms.LoginForm
	err := ctx.Bind(&form)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	usr, err := u.srv.Login(&form, secret)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.JSON(200, usr)
}

func (u *User) updatePasswordHandler(ctx *gin.Context) {
	var form forms.UpdatePassword
	err := ctx.Bind(&form)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	v, ok := ctx.Get("user")
	if !ok {
		log.Println("get user failed")
		ctx.String(400, "token expired")
		return
	}

	form.UserID = v.(*jwt.CustomClaims).ID

	err = u.srv.UpdatePassword(&form)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.String(200, "update password success")
}

func (u *User) findAllUsersHandler(ctx *gin.Context) {
	users, err := u.srv.FindAllUsers()
	if err != nil {
		log.Println(err)
		ctx.String(500, err.Error())
		return
	}

	ctx.JSON(200, users)
}

// func (u *User) changeUserStatusHandler(ctx *gin.Context) {
// 	id :=
// }
