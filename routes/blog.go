package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"hypermedlab/backend-myblog/pkgs/forms"
	"hypermedlab/backend-myblog/pkgs/jwt"
	"hypermedlab/backend-myblog/pkgs/middlewares"
	blogSrv "hypermedlab/backend-myblog/services/blog"
	"log"
)

type blog struct {
	srv    blogSrv.Service
	engine *gin.Engine
}

func NewBlogRouter(conn *sqlx.DB, e *gin.Engine) Router {
	return &blog{
		srv:    blogSrv.NewBlogService(conn),
		engine: e,
	}
}

func (b *blog) Init() {
	blogroute := b.engine.Group("/blog")
	blogroute.POST("", middlewares.IPCount, middlewares.AuthToken, b.createBlog)
	blogroute.GET("", middlewares.IPCount, middlewares.AuthToken, b.findBlogsByUser)
	blogroute.GET("/id/:id", middlewares.IPCount, b.findBlogByID)
	blogroute.GET("/title/:title", middlewares.IPCount, b.findBlogsByTitle)
	blogroute.GET("/userid/:userid", middlewares.IPCount, b.findBlogByUserID)
	blogroute.GET("/username/:username", middlewares.IPCount, b.findBlogsByUserName)
	blogroute.PUT("/:id", middlewares.IPCount, middlewares.AuthToken, b.updateBlog)
	blogroute.DELETE("/id/:id", middlewares.IPCount, middlewares.AuthToken, b.deleteBlogByID)
	blogroute.DELETE("/userid/:userid", middlewares.IPCount, middlewares.AuthToken, b.deleteBlogsByUserID)

	commentroute := b.engine.Group("/comment")
	commentroute.POST("", middlewares.IPCount, middlewares.AuthToken, b.createComment)
	commentroute.GET("/id/:id", middlewares.IPCount, b.findCommentByID)
	commentroute.GET("/userid/:userid", middlewares.IPCount, middlewares.AuthToken, b.findCommentsByUserID)
	commentroute.PUT("/:id", middlewares.IPCount, middlewares.AuthToken, b.updateComment)
	commentroute.DELETE("/id/:id", middlewares.IPCount, middlewares.AuthToken, b.deleteCommentByID)
	commentroute.DELETE("/userid/:userid", middlewares.IPCount, middlewares.AuthToken, b.deleteCommentsByUserID)
}

func (b *blog) createBlog(ctx *gin.Context) {
	var form forms.CreateBlog
	err := ctx.Bind(&form)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	v, ok := ctx.Get("user")
	if !ok {
		ctx.String(401, "unAuthorized")
		return
	}
	c := v.(*jwt.CustomClaims)
	form.UserID = c.ID
	form.UserName = c.Name

	blg, err := b.srv.CreateBlog(&form)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.JSON(200, blg)
}

func (b *blog) createComment(ctx *gin.Context) {
	var form forms.CreateComment

	err := ctx.Bind(&form)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	v, ok := ctx.Get("user")
	if !ok {
		ctx.String(401, "unAuthorized")
		return
	}
	claims := v.(*jwt.CustomClaims)
	form.UserID = claims.ID
	form.UserName = claims.Name

	blg, err := b.srv.CreateComment(&form)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.JSON(200, blg)
}

func (b *blog) findBlogByID(ctx *gin.Context) {

	id := ctx.Param("id")

	blg, err := b.srv.FindBlogByID(id)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.JSON(200, blg)
}

func (b *blog) findBlogsByTitle(ctx *gin.Context) {

	title := ctx.Param("title")
	blogs, err := b.srv.FindBlogsByTitle(title)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}
	ctx.JSON(200, blogs)
}

func (b *blog) findBlogsByUser(ctx *gin.Context) {

	v, ok := ctx.Get("user")
	if !ok {
		ctx.String(401, "token expired")
		return
	}
	claims := v.(*jwt.CustomClaims)
	blogs, err := b.srv.FindBlogsByUserID(claims.ID)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.JSON(200, blogs)
}

func (b *blog) findBlogByUserID(ctx *gin.Context) {
	id := ctx.Param("userid")
	blogs, err := b.srv.FindBlogsByUserID(id)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.JSON(200, blogs)
}

func (b *blog) findBlogsByUserName(ctx *gin.Context) {
	username := ctx.Param("username")
	blogs, err := b.srv.FindBlogsByUserName(username)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.JSON(200, blogs)

}

func (b *blog) findCommentByID(ctx *gin.Context) {
	id := ctx.Param("userid")
	comment, err := b.srv.FindCommentByID(id)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}
	ctx.JSON(200, comment)
}

func (b *blog) findCommentsByUserID(ctx *gin.Context) {
	userid := ctx.Param("userid")
	comments, err := b.srv.FindCommentsByUserID(userid)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.JSON(200, comments)
}

func (b *blog) updateBlog(ctx *gin.Context) {
	var form forms.UpdateBlog
	err := ctx.Bind(&form)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	id := ctx.Param("id")
	v, ok := ctx.Get("user")
	if !ok {
		ctx.String(401, "unAuthorized")
		return
	}
	claims := v.(*jwt.CustomClaims)
	form.UserID = claims.ID
	form.BlogID = id

	blg, err := b.srv.UpdateBlog(&form)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.JSON(200, blg)
}

func (b *blog) updateComment(ctx *gin.Context) {
	var form forms.UpdateComment
	err := ctx.Bind(&form)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	id := ctx.Param("id")
	v, ok := ctx.Get("user")
	if !ok {
		ctx.String(401, "unAuthorized")
		return
	}
	claims := v.(*jwt.CustomClaims)
	form.CommentID = id
	form.UserID = claims.ID

	comment, err := b.srv.UpdateComment(&form)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.JSON(200, comment)
}

func (b *blog) deleteBlogByID(ctx *gin.Context) {
	id := ctx.Param("id")
	err := b.srv.DeleteBlogByID(id)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.String(200, "delete success")
}

func (b *blog) deleteBlogsByUserID(ctx *gin.Context) {
	v, ok := ctx.Get("user")
	if !ok {
		ctx.String(401, "unAuthorized")
		return
	}
	claims := v.(*jwt.CustomClaims)

	err := b.srv.DeleteBlogByUserID(claims.ID)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.String(200, "delete success")
}

func (b *blog) deleteCommentByID(ctx *gin.Context) {
	id := ctx.Param("id")
	err := b.srv.DeleteCommentByID(id)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.String(200, "delete success")
}

func (b *blog) deleteCommentsByUserID(ctx *gin.Context) {
	v, ok := ctx.Get("user")
	if !ok {
		ctx.String(401, "unAuthorized")
		return
	}
	claims := v.(*jwt.CustomClaims)

	err := b.srv.DeleteCommentByUserID(claims.ID)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.String(200, "delete success")
}
