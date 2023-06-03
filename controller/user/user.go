package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (c *UserController) login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		if username == "admin" && password == "fuckdingding" {
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
		databasePath := ctx.PostForm("databasePath")
		recordDb, err := gorm.Open(sqlite.Open(databasePath+"attendance_record.db"), &gorm.Config{})
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": err.Error(),
			})
			return
		}
		userDb, err := gorm.Open(sqlite.Open(databasePath+"user.db"), &gorm.Config{})
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": err.Error(),
			})
			return
		}
		_ = userDb.AutoMigrate(&User{})
		_ = recordDb.AutoMigrate(&Record{})
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

func (c *UserController) getUserList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		databasePath := ctx.PostForm("databasePath")
		userDb, err := gorm.Open(sqlite.Open(databasePath+"user.db"), &gorm.Config{})
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": err.Error(),
			})
			return
		}
		_ = userDb.AutoMigrate(&Record{})
		var userList []User
		userDb.Find(&userList)
		ctx.JSON(http.StatusOK, gin.H{
			"code":     200,
			"message":  "success",
			"userList": userList,
		})
	}
}

func (c *UserController) addRecord() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		databasePath := ctx.PostForm("databasePath")
		recordDb, err := gorm.Open(sqlite.Open(databasePath+"attendance_record.db"), &gorm.Config{})
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": err.Error(),
			})
			return
		}
		_ = recordDb.AutoMigrate(&Record{})
		userId := ctx.PostForm("user_id")
		timestamp := ctx.PostForm("timestamp")
		status := ctx.PostForm("status")
		result := recordDb.Create(RecordDto{
			UserId:    userId,
			Timestamp: StrToInt(timestamp),
			Status:    StrToInt(status),
		})
		if result.Error == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "添加记录成功",
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": result.Error,
			})
		}

	}
}

func (c *UserController) deleteRecord() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		databasePath := ctx.PostForm("databasePath")
		recordDb, err := gorm.Open(sqlite.Open(databasePath+"attendance_record.db"), &gorm.Config{})
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": err.Error(),
			})
			return
		}
		recordDb.AutoMigrate(&Record{})
		id := ctx.PostForm("id")
		recordDb.Delete(Record{
			Id: StrToInt(id),
		})
		ctx.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "删除记录成功",
		})
	}
}

type UserInfo struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Sex  string `json:"sex"`
}

type Record struct {
	Id        int    `gorm:"id"`
	UserId    string `gorm:"user_id"`
	Timestamp int    `gorm:"timestamp"`
	Status    int    `gorm:"status"`
	UserName  string `gorm:"-"`
}

type RecordDto struct {
	UserId    string `gorm:"user_id"`
	Timestamp int    `gorm:"timestamp"`
	Status    int    `gorm:"status"`
}

type User struct {
	Id     uint
	UserId string
	Name   string
}

func (Record) TableName() string {
	return "record"
}

func (RecordDto) TableName() string {
	return "record"
}

func (User) TableName() string {
	return "user"
}

func StrToInt(str string) int {
	i, e := strconv.Atoi(str)
	if e != nil {
		return 0
	}
	return i
}
