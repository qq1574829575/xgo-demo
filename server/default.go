package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IController interface {
	Router(server *Server)
}

type Server struct {
	*gin.Engine
	// 路由分组
	routerGroup *gin.RouterGroup
}

func Init() *Server {
	// 作为Server的构造器
	s := &Server{Engine: gin.New()}
	s.StaticFS("/static", http.Dir("./static"))
	// 返回作为链式调用
	return s
}

func (s *Server) SetMiddlewares(middleware ...gin.HandlerFunc) *Server {
	s.Use(middleware...)
	return s
}

func (s *Server) SetStaticDir() *Server {
	s.StaticFS("/upload", http.Dir("./upload"))
	return s
}

func (s *Server) Listen() {
	err := s.Run(":9000")
	if err != nil {
		return
	}
}

func (s *Server) Route(controllers ...IController) *Server {
	// 遍历所有的控制层，这里使用接口，就是为了将Router实例化
	for _, controller := range controllers {
		controller.Router(s)
	}
	return s
}
