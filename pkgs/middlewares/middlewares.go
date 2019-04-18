package middlewares

import (
	"hypermedlab/backend-myblog/pkgs/jwt"
	"hypermedlab/backend-myblog/pkgs/session"

	"github.com/gin-gonic/gin"
)

const Secret = "哈哈没想到吧"

type MiddleWare struct {
	sessions *session.SessionsStorageInMemory
}

func NewMiddleWare(sessions *session.SessionsStorageInMemory) *MiddleWare {
	return &MiddleWare{
		sessions: sessions,
	}
}

//IPCount to
func IPCount(ctx *gin.Context) {
	// c, err := redis.Dial("tcp", "localhost:6379")
	// if err != nil {
	// 	ctx.JSON(500, "Internal Error")
	// 	log.Println("Connect to redis error", err)
	// 	ctx.Abort()
	// 	return
	// }
	// defer c.Close()

	// ip := ctx.Request.RemoteAddr
	// c.Send("INCR", ip)
	// c.Send("EXPIRE", ip, 10)
	// resp, err := redis.Int(c.Do("GET", ip))
	// if err != nil {
	// 	ctx.JSON(500, "Internal Error")
	// 	ctx.Abort()
	// 	return
	// }
	// if resp >= 10 {
	// 	ctx.JSON(403, "too many request")
	// 	ctx.Abort()
	// 	return
	// }
}

func (m *MiddleWare) AuthToken(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")
	claims, err := jwt.ValidateToken(token, Secret)
	if err != nil || claims == nil {
		ctx.JSON(401, err.Error())
		ctx.Abort()
		return
	}

	sess, ok := m.sessions.Get(claims.ID)
	if !ok {
		ctx.String(401, "err token not in session")
		ctx.Abort()
		return
	}

	if sess.GetData().(string) != token {
		ctx.String(401, "err token is expired")
		ctx.Abort()
		return
	}

	ctx.Set("user", claims)
}

func (m *MiddleWare) AdminAuthToken(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")
	claims, err := jwt.ValidateToken(token, Secret)
	if err != nil || claims == nil {
		ctx.JSON(401, err.Error())
		ctx.Abort()
		return
	}

	if claims.ID != "000000" {
		ctx.JSON(401, err.Error())
		ctx.Abort()
		return
	}

	sess, ok := m.sessions.Get(claims.ID)
	if !ok {
		ctx.String(401, "err token not in session")
		ctx.Abort()
		return
	}

	if sess.GetData().(string) != token {
		ctx.String(401, "err token is expired")
		ctx.Abort()
		return
	}

	ctx.Set("user", claims)
}
