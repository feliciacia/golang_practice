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
	"golang.org/x/crypto/bcrypt"
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
	user.Username = ctx.PostForm("username")
	password := ctx.PostForm("password")
	err := user.getUserByUsername()
	if err != nil {
		fmt.Println("error selecting password_hash in db by username, err:", err)
		ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{"message": "check username and password"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.pswdHash), []byte(password)) //using bycrypt for hashing password
	fmt.Println("err from bycrypt:", err)
	//if hash and password match
	if err == nil {
		session, _ := store.Get(ctx.Request, "session")
		session.Values["user"] = user
		session.Save(ctx.Request, ctx.Writer)
		ctx.HTML(http.StatusOK, "loggedin.html", gin.H{"username": user.Username})
		return
	}
	ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{"message": "check username and password"}) //handling unauthorized login
}

func profileHandler(ctx *gin.Context) {
	session, _ := store.Get(ctx.Request, "session")
	var user = &User{}            //user informations assigned in user value
	val := session.Values["user"] //store user informations in session
	var ok bool
	if user, ok = val.(*User); !ok {
		fmt.Println("was not of type *User")
		ctx.HTML(http.StatusForbidden, "login.html", nil) //so that other user cant render in the page of user's profile
		return
	}
	ctx.HTML(http.StatusOK, "profile.html", gin.H{"user": user})
}

func (u *User) getUserByUsername() error {
	stmt := "SELECT * FROM users WHERE username = ?"
	row := db.QueryRow(stmt, u.Username) //? will be replaced by user's username
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.pswdHash, &u.CreatedAt, &u.Active, &u.verHash, &u.timeout)
	if err != nil {
		fmt.Println("getUser() error selecting User, err:", err)
		return err
	}
	return nil //successful retrieval will return nil
}
