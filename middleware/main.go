package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type movie struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Director string `json:"director"`
	Price    string `json:"price"`
}

var movies = []movie{
	{ID: "1", Title: "The Dark Knight", Director: "Christopher Nolan", Price: "5.99"},
	{ID: "2", Title: "Tommy Boy", Director: "Peter Segal", Price: "2.99"},
	{ID: "3", Title: "The Shawshank Redemption", Director: "Frank Darabont", Price: "7.99"},
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("middleware/templates/*.html")
	router.Use(middlewareFunc1(), middlewareFunc2, middlewareFunc3())
	//router.Use(StartReqTimeAndFinishReqTimeMiddleware())
	//router.Use(CheckAuthMiddleware())
	/*router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "succeed",
		})
	})*/
	router.GET("/movie", getAllMovies)
	authRouter := router.Group("/auth", gin.BasicAuth(gin.Accounts{
		"Joe":   "baseball",
		"Kelly": "1234",
	}))
	authRouter.GET("/movie", createMovieForm)
	authRouter.POST("/movie", createMovie)
	router.Run(":8080")
}

func StartReqTimeAndFinishReqTimeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("Start at:", time.Now())
		ctx.Next()
		log.Println("End at:", time.Now())
	}
}

func CheckAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token != "secrettoken" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "You have no authorited token",
			})
		}
		ctx.Next()
	}
}

func middlewareFunc1() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("middlewareFunc1 running")
		ctx.Next()
	}
}

func middlewareFunc2(ctx *gin.Context) {
	fmt.Println("middlewareFunc2 running")
	//ctx.Abort()
	fmt.Println("middlewareFunc2 ending")
	ctx.Next()
}

func middlewareFunc3() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("middlewareFunc3 running")
		ctx.Next()
	}
}

func getAllMovies(c *gin.Context) {
	c.HTML(http.StatusOK, "allmovies.html", movies)
}

func createMovieForm(c *gin.Context) {
	c.HTML(http.StatusOK, "createmovieform.html", nil)
}

func createMovie(c *gin.Context) {
	var newMovie movie
	newMovie.ID = c.PostForm("id")
	newMovie.Title = c.PostForm("title")
	newMovie.Director = c.PostForm("director")
	newMovie.Price = c.PostForm("price")
	movies = append(movies, newMovie)
	c.HTML(http.StatusOK, "allmovies.html", movies)
}
