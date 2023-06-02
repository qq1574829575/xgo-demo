package user

import "dingding/server"

type UserController struct {
}

func Controllers() *UserController {
	return &UserController{}
}

func (c *UserController) Router(server *server.Server) {
	server.Handle("POST", "/user/login", c.login())
	server.Handle("GET", "/user/getUserInfo", c.getUserInfo())
	server.Handle("POST", "/user/getRecordList", c.getRecordList())
}
