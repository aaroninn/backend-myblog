package routes

import (
	"hypermedlab/backend-myblog/pkgs/forms"
	"hypermedlab/backend-myblog/pkgs/jwt"
	"hypermedlab/backend-myblog/pkgs/middlewares"
	userSrv "hypermedlab/backend-myblog/services/user"

	"github.com/gin-gonic/gin"

	"log"
)

var secret string

type User struct {
	srv        *userSrv.Service
	Engine     *gin.Engine
	middleware *middlewares.MiddleWare
}

func NewUserRouter(srv *userSrv.Service, engine *gin.Engine, se string, middleware *middlewares.MiddleWare) *User {
	secret = se
	return &User{
		srv:        srv,
		Engine:     engine,
		middleware: middleware,
	}

}

func (u *User) Init() {
	userroute := u.Engine.Group("/user")
	userroute.POST("/login", middlewares.IPCount, u.loginHandler)
	userroute.POST("/register", middlewares.IPCount, u.registerHandler)
	userroute.PUT("/logout", middlewares.IPCount, u.middleware.AuthToken, u.logOutHandler)
	userroute.PUT("/password", middlewares.IPCount, u.middleware.AuthToken, u.updatePasswordHandler)

	adminroute := u.Engine.Group("/admin")
	adminroute.GET("/users", middlewares.IPCount, u.middleware.AdminAuthToken, u.findAllUsersHandler)
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

func (u *User) logOutHandler(ctx *gin.Context) {
	v, ok := ctx.Get("user")
	if !ok {
		log.Println("get user failed")
		ctx.String(400, "token expired")
		return
	}

	u.srv.LogOut(v.(*jwt.CustomClaims).ID)
	ctx.String(200, "log out success")
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
