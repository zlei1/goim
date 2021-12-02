package http

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func (s *Server) push(c *gin.Context) {
	var arg struct {
		UserId int64 `form:"user_id"`
	}
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	if arg.UserId == 0 {
		errors(c, RequestErr, "user_id empty")
		return
	}

	msg, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	if err = s.logic.Push(arg.UserId, msg); err != nil {
		errors(c, ServerErr, err.Error())
		return
	}
	result(c, nil, OK)
}

func (s *Server) pushRoom(c *gin.Context) {
	var arg struct {
		RoomId int64 `form:"room_id"`
	}
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	if arg.RoomId == 0 {
		errors(c, RequestErr, "room_id empty")
		return
	}

	msg, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	if err = s.logic.PushRoom(arg.RoomId, msg); err != nil {
		errors(c, ServerErr, err.Error())
		return
	}
	result(c, nil, OK)
}