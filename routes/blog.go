package routes

import (
	"hypermedlab/backend-myblog/pkgs/forms"
	"hypermedlab/backend-myblog/pkgs/jwt"
	"hypermedlab/backend-myblog/pkgs/middlewares"
	blogSrv "hypermedlab/backend-myblog/services/blog"

	"github.com/gin-gonic/gin"

	"log"
)

type Blog struct {
	srv    *blogSrv.Service
	Engine *gin.Engine
}

func NewBlogRouter(srv *blogSrv.Service, e *gin.Engine) *Blog {
	return &Blog{
		srv:    srv,
		Engine: e,
	}
}

func (b *Blog) Init() {
	blogroute := b.Engine.Group("/blog")
	blogroute.POST("", middlewares.IPCount, middlewares.AuthToken, b.createBlog)
	blogroute.GET("", middlewares.IPCount, middlewares.AuthToken, b.findBlogsByUser)
	blogroute.GET("/id/:id", middlewares.IPCount, b.findBlogByID)
	blogroute.GET("/title/:title", middlewares.IPCount, b.findBlogsByTitle)
	blogroute.GET("/userid/:userid", middlewares.IPCount, b.findBlogByUserID)
	blogroute.GET("/username/:username", middlewares.IPCount, b.findBlogsByUserName)
	blogroute.PUT("/:id", middlewares.IPCount, middlewares.AuthToken, b.updateBlog)
	blogroute.DELETE("/id/:id", middlewares.IPCount, middlewares.AuthToken, b.deleteBlogByID)
	blogroute.DELETE("/userid/:userid", middlewares.IPCount, middlewares.AuthToken, b.deleteBlogsByUserID)

	commentroute := b.Engine.Group("/comment")
	commentroute.POST("", middlewares.IPCount, middlewares.AuthToken, b.createComment)
	commentroute.GET("/id/:id", middlewares.IPCount, b.findCommentByID)
	commentroute.GET("/userid/:userid", middlewares.IPCount, middlewares.AuthToken, b.findCommentsByUserID)
	commentroute.PUT("/:id", middlewares.IPCount, middlewares.AuthToken, b.updateComment)
	commentroute.DELETE("/id/:id", middlewares.IPCount, middlewares.AuthToken, b.deleteCommentByID)
	commentroute.DELETE("/userid/:userid", middlewares.IPCount, middlewares.AuthToken, b.deleteCommentsByUserID)
}

func (b *Blog) createBlog(ctx *gin.Context) {
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

func (b *Blog) createComment(ctx *gin.Context) {
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

func (b *Blog) findBlogByID(ctx *gin.Context) {

	id := ctx.Param("id")

	blg, err := b.srv.FindBlogByID(id)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.JSON(200, blg)
}

func (b *Blog) findBlogsByTitle(ctx *gin.Context) {

	title := ctx.Param("title")
	blogs, err := b.srv.FindBlogsByTitle(title)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}
	ctx.JSON(200, blogs)
}

func (b *Blog) findBlogsByUser(ctx *gin.Context) {

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

func (b *Blog) findBlogByUserID(ctx *gin.Context) {
	id := ctx.Param("userid")
	blogs, err := b.srv.FindBlogsByUserID(id)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.JSON(200, blogs)
}

func (b *Blog) findBlogsByUserName(ctx *gin.Context) {
	username := ctx.Param("username")
	blogs, err := b.srv.FindBlogsByUserName(username)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.JSON(200, blogs)

}

func (b *Blog) findCommentByID(ctx *gin.Context) {
	id := ctx.Param("userid")
	comment, err := b.srv.FindCommentByID(id)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}
	ctx.JSON(200, comment)
}

func (b *Blog) findCommentsByUserID(ctx *gin.Context) {
	userid := ctx.Param("userid")
	comments, err := b.srv.FindCommentsByUserID(userid)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.JSON(200, comments)
}

func (b *Blog) updateBlog(ctx *gin.Context) {
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

func (b *Blog) updateComment(ctx *gin.Context) {
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

func (b *Blog) deleteBlogByID(ctx *gin.Context) {
	id := ctx.Param("id")
	err := b.srv.DeleteBlogByID(id)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.String(200, "delete success")
}

func (b *Blog) deleteBlogsByUserID(ctx *gin.Context) {
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

func (b *Blog) deleteCommentByID(ctx *gin.Context) {
	id := ctx.Param("id")
	err := b.srv.DeleteCommentByID(id)
	if err != nil {
		log.Println(err)
		ctx.String(400, err.Error())
		return
	}

	ctx.String(200, "delete success")
}

func (b *Blog) deleteCommentsByUserID(ctx *gin.Context) {
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
