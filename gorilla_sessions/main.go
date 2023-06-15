package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

type User struct {
	ID        string
	Username  string
	Email     string
	pswdHash  string
	CreatedAt string
	Active    string
	verHash   string
	timeout   string
}

var db *sql.DB
var store = sessions.NewCookieStore([]byte("super-secret"))

func init() {
	store.Options.HttpOnly = true
	store.Options.Secure = true
	gob.Register(&User{})
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	var err error
	db, err = sql.Open("mysql", "root:super-secret-password@tcp(localhost:3306)/gin-db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	authRouter := router.Group("/user", auth)

	router.GET("/", indexHandler)
	router.GET("/login", loginGEThandler)
	router.POST("/login", loginPOSThandler)

	authRouter.GET("/profile", profileHandler)
	err = router.Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
}

func auth(ctx *gin.Context) {
	fmt.Println("auth middleware running")
	session, _ := store.Get(ctx.Request, "session")
	fmt.Println("session:", session)
	_, ok := session.Values["user"]
	if !ok {
		ctx.HTML(http.StatusForbidden, "login.html", nil)
		ctx.Abort()
		return
	}
	fmt.Println("middleware done")
	ctx.Next()
}

func indexHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}

func loginGEThandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", nil)
}

func loginPOSThandler(ctx *gin.Context) {
	var user User
	user.Username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	err := user.getUserbyUsername()
	if err != nil {
		fmt.Println("error selecting password_hash in db by username, err:", err)
		ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{"message": "check username and password"})
		return
	}
}
