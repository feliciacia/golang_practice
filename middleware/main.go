package main

import (
	"fmt"
	"net/http"

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
	router := gin.New()
	router.LoadHTMLGlob("middleware/templates/*.html")
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewareFunc1, middlewareFunc2, middlewareFunc3())
	router.GET("/movie", getAllMovies)
	authRouter := router.Group("/auth", gin.BasicAuth(gin.Accounts{
		"Joe":   "baseball",
		"Kelly": "1234",
	}))
	authRouter.GET("/movie", createMovieForm)
	authRouter.POST("/movie", createMovie)
	router.Run(":8000")
}

func middlewareFunc1(c *gin.Context) {
	fmt.Println("middlewareFunc1 running")

	c.Next()
}

func middlewareFunc2(c *gin.Context) {
	fmt.Println("middlewareFunc2 running")
	c.Abort()
	fmt.Println("middlewareFunc2 ending")
	c.Next()
}

func middlewareFunc3() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("middlewareFunc3 running")
		c.Next()
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
