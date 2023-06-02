package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

func (c *UserController) login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		if username == "admin" && password == "123456" {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "登录成功",
				"token":   "hdu19doidj01ue",
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": "用户名或密码错误:" + username + password,
			})
		}
	}
}

func (c *UserController) getUserInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data": UserInfo{
				Id:   1,
				Name: "管理员",
				Sex:  "男",
			},
		})
	}
}

func (c *UserController) getRecordList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		recordDb, err := gorm.Open(sqlite.Open("attendance_record.db"), &gorm.Config{})
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": err.Error(),
			})
			return
		}
		userDb, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": err.Error(),
			})
			return
		}
		userDb.AutoMigrate(&User{})
		recordDb.AutoMigrate(&Record{})
		var recordList []Record
		recordDb.Find(&recordList)
		for i, record := range recordList {
			var user User
			userDb.Where("user_id = ?", record.UserId).First(&user)
			recordList[i].UserName = user.Name
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code":       200,
			"message":    "success",
			"recordList": recordList,
		})
	}
}

type UserInfo struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Sex  string `json:"sex"`
}

type Record struct {
	Id        uint
	UserId    string
	Timestamp int
	Status    int
	UserName  string
}

type User struct {
	Id     uint
	UserId string
	Name   string
}

func (Record) TableName() string {
	return "record"
}

func (User) TableName() string {
	return "user"
}
