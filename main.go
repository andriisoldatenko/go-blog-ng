package main

import (
	"github.com/andriisoldatenko/go-blog-ng/auth"
	"github.com/andriisoldatenko/go-blog-ng/posts"
	"github.com/andriisoldatenko/go-blog-ng/storage"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	db := storage.Init()
	posts.Migrate(db)
	defer db.Close()

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("./frontend/dist/frontend/", true)))
	r.Use(static.Serve("/implicit/callback", static.LocalFile("./frontend/dist/frontend/", true)))
	r.Use(static.Serve("/posts", static.LocalFile("./frontend/dist/frontend/", true)))
	r.Use(static.Serve("/posts/new", static.LocalFile("./frontend/dist/frontend/", true)))

	v1 := r.Group("/api")
	v1.Use(auth.AuthMiddleware())
	posts.PostsRegister(v1.Group("/posts"))

	testAuth := r.Group("/api/ping")
	testAuth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run()
}