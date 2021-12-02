package http

import (
	"github.com/gin-gonic/gin"
	"github.com/zlei1/goim/internal/logic"
)

type Server struct {
	engine *gin.Engine
	logic *logic.Logic
}

func InitHTTP(l *logic.Logic) {
	engine := gin.New()
	engine.Use(loggerHandler, recoverHandler)

	go func() {
		if err := engine.Run(l.C.Base.HttpAddr); err != nil {
			panic(err)
		}
	}()

	s := &Server{
		engine: engine,
		logic: l,
	}
	s.initRouter()
}

func (s *Server) initRouter() {
	group := s.engine.Group("/im")
	group.POST("/push", s.push)
	group.POST("/push/rooms", s.pushRoom)
}

func (s *Server) Close() {
}