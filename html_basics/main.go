package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("html_basics/templates/*.html")
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
	c.HTML(http.StatusOK, "greeting.html", nil)
}

func getGreetName(c *gin.Context) {
	name := c.Param("name")
	c.HTML(http.StatusOK, "customGreeting.html", name)
}

func getManyData(c *gin.Context) {
	foods := []string{"chicken sandwich", "fries", "soda", "cookie"}
	c.HTML(http.StatusOK, "manyData.html", gin.H{
		"name":  "Carl",
		"foods": foods,
	})
}

func getForm(c *gin.Context) {
	c.HTML(http.StatusOK, "form.html", nil)
}

func postForm(c *gin.Context) {
	name := c.PostForm("name")
	food := c.PostForm("food")
	c.HTML(http.StatusOK, "formresult.html", gin.H{
		"name": name,
		"food": food,
	})
}
