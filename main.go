package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)
var DB *gorm.DB

type Book struct {
	ID int `json:"id""`
	Name string `json:"name"`
	Price float32 `json:"price"`
}

func initMySQl() (err error)  {
	dsn := "root:123456@(127.0.0.1:3306)/tesdb?charset=utf8&mb4"
	DB, err = gorm.Open("mysql",dsn)
	if err != nil{
		return
	}
	return DB.DB().Ping()
}


func main()  {
	err := initMySQl()
	if err != nil{
		panic(err)
	}
	defer DB.Close()
	//模型绑定
	DB.AutoMigrate(&Book{})
	r := gin.Default()
	r.GET("/index", func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"msg":"This is my homepage!",
		})
	})
	r.GET("/add", func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"msg":"Add book success!",
		})
	})
	r.POST("/add", func(c *gin.Context) {
		//前端页面的请求会提交到这儿
		//1.从请求中吧数据拿出来
		//2.存入数据库
		//3.返回相应
		var book Book
		c.BindJSON(&book)
		if err = DB.Create(&book).Error;err != nil{
			c.JSON(http.StatusOK,gin.H{"error":err.Error()})
		}else {
			c.JSON(http.StatusOK,book)
		}
	})
	r.GET("/bookall", func(c *gin.Context) {
		var bookList []Book
		if err = DB.Find(&bookList).Error;err!=nil{
			c.JSON(http.StatusOK,gin.H{"error":err.Error()})
		}else {
			c.JSON(http.StatusOK,bookList)
		}
	})
	r.PUT("/book/:id", func(c *gin.Context) {
		id,ok := c.Params.Get("id")
		if !ok{
			c.JSON(http.StatusOK,gin.H{"error":"id not exist"})
			return
		}
		var book Book
		if err = DB.Where("id=?",id).First(&book).Error; err != nil{
			c.JSON(http.StatusOK,gin.H{"error":err.Error()})
		}
		c.BindJSON(&book)
		if err = DB.Save(&book).Error;err != nil{
			c.JSON(http.StatusOK,gin.H{"error":err.Error()})
		}else {
			c.JSON(http.StatusOK,book)
		}
	})
	r.DELETE("/book/:id", func(c *gin.Context) {
		id,ok := c.Params.Get("id")
		if !ok{
			c.JSON(http.StatusOK,gin.H{"error":"id not exist"})
			return
		}
		if err = DB.Where("id=?",id).Delete(Book{}).Error;err != nil{
			c.JSON(http.StatusOK,gin.H{"error":err.Error()})
		}else{
			c.JSON(http.StatusOK,gin.H{"id":"Delete success"})
		}
	})
	r.Run(":8888")
}


