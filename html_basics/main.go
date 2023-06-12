package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("template/*.html")
	router.GET("/hello", getHello)
	router.GET("/greet", getGreet)
	router.GET("/greet/:name", getGreetName)
	router.GET("/many", getManyData)
	router.GET("/form", getForm)
	router.POST("/form", postForm)

	err := router.Run("localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
}

func getHello(c *gin.Context) {
	c.String(http.StatusOK, "Hello world!")
}

func getGreet(c *gin.Context) {
	c.String
}
