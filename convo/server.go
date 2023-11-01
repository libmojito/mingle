package convo

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	address string
}

func NewServer(options ...func(*Server)) *Server {
	svr := &Server{}

	for _, o := range options {
		o(svr)
	}

	return svr
}

func WithAddress(addr string) func(*Server) {
	return func(svr *Server) {
		svr.address = addr
	}
}

func (svr *Server) Run() {
	r := gin.Default()

	r.POST("/thread", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"thread_id": "thread_1234567890",
		})
	})

	r.POST("chat", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.JSON(http.StatusOK, gin.H{
			"response": "hello back, this is longer than you think, get ready for it. ok? Aparently it's not long enough. If you can still see this part, it means text wrap is working as expected.",
		})
	})

	r.Run(svr.address)
}
