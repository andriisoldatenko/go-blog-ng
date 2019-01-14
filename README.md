---
layout: blog_post
title: "Build a Simple CRUD App with Go and Angular"
author: a_soldatenko
description: "This tutorial explains how to build a modern CRUD application with a Go as backend and an Angular frontend."
tags: [go, angular, crud, okta, jwt, gin]
---

# Build a Simple CRUD App with Go and Angular

## Instead of introduction

Go was invented in 2009 as a replacement of C by Robert Griesemer, Rob Pike, and Ken Thompson to solve system
programming real problems. Nowadays we can see huge amount of usage of Go in web, cloud native and even web for client
[WebAssembly](https://github.com/golang/go/wiki/WebAssembly#introduction).
But today I would like to share experience how to build modern blog application, where you
can create, read, update and delete blog posts.

## Few words about requirements

For go through this tutorial we will use for backed some Go requirements:

### Backend:

- `gin` for building REST Api
- `gorm` for manipulating with database
- `okta-jwt-verifier-golang` for manipulates with Json web tokens
- `delve` to demonstrate how to debug code

### Database:

- `SQLite` - opensource relational database written in C.

### Frontend:

- `Angular 7` - modern opensource frontend framework
- npm

### Authontication && autorization:

- Okta

## CRUD acronym

CRUD refer to SQL equivalents:

| Action        | SQL statement | HTTP Verb | URL Path    |
| ------------- |---------------| ----------|-------------|
| Create        | INSERT        | POST      | /posts      |
| Read          | SELECT        | GET       | /posts      |
| Update        | UPDATE        | PUT       | /post/<id>  |
| Delete        | DELETE        | DELETE    | /post/<id>  |


## Let's start from init project layout

I prefer to start from something real simple, and than step by step add / refactor code and structure.
But before we start let's test environment and installation. Assume you know how to install Go. If not please follow
this [getting started](https://golang.org/doc/install) instruction. After succesfully installation you can see something like it:

```bash
$ go version
go version go1.11.4 darwin/amd64
```

## Create project structure:

First of all let's create project folder `go-blog-ng` and put `main.go` inside:

```bash
ls -l go-blog-ng
main.go
```

And now we can add simple hello okta code to test `main.go` program:

```go

package main

import "fmt"

func main() {
	fmt.Printf("Hello, Okta!\n")
}
```

## Build and run

```bash
go run main.go
Hello, Okta!
```

## Add first dependency `gin`:

As you can understand I need to install some web-framework, because I don't want to reinvent the wheel.
I prefer [Gin](https://github.com/gin-gonic/gin) http web framework written in Go for performance and productivity.

## Dependency management in Go

Today is't complicated to imagine modern web application without any dependency. By the way it's huge holy war in Go community
about package management and especially compare to Rust community. But I'm nit going to start discussing it here, I just
show you latest available in Go 1.11 [go modules](https://github.com/golang/go/wiki/Modules) from box which was proposed
in [Proposal: Versioned Go Modules](https://go.googlesource.com/proposal/+/master/design/24301-versioned-go.md)
by Russ Cox.

To start using some dependency first create `go.mod` file inside `go-blog-ng` folder:

```go
module github.com/andriisoldatenko/go-blog-ng

require (
	github.com/gin-gonic/gin v1.3.0
)
```

and run `go run main.go` again, and you can see how go modules starts to download `gin` module.


## ping/pong GIN app

Now you can create simple ping/pong app to check installation.

```bash
cat main.go
package main

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return r
}


func main() {
	r := setupRouter()
	r.Run(":8080")
}
```

TIP:

> If you add a new import to your source code that is not yet covered by a require in go.mod, any go command run
> (e.g., 'go build') will automatically look up the proper module and add the highest version of that new direct
> dependency to your module's go.mod as a require directive.

```bash
go run main.go
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /ping                     --> main.main.func1 (3 handlers)
[GIN-debug] Environment variable PORT is undefined. Using port :8080 by default
[GIN-debug] Listening and serving HTTP on :8080
[GIN] 2018/10/29 - 17:57:39 | 404 |         766ns |             ::1 | GET      /
[GIN] 2018/10/29 - 17:57:39 | 404 |       1.108µs |             ::1 | GET      /favicon.ico
[GIN] 2018/10/29 - 17:57:43 | 200 |      75.323µs |             ::1 | GET      /ping
^Csignal: interrupt
```

## Few notes regarding testing

World divided into 2 parts TDD fans and engineers who understands that we have to write tests. Both sides
understands that test fundament for refactoring and future quality, let's follow this rules.
In golang there are great set of tools already included:
```bash
Go is a tool for managing Go source code.
	doc         show documentation for package or symbol
	test        test packages
	tool        run specified go tool
	version     print Go version
	... some parts missing
```
Note:
`go test` expexted go files with `_test.go` prefix.

Let's create first test:

```bash
cat main_test.go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}
```

```bash
go test
PASS
ok  	github.com/andriisoldatenko/go-blog	0.019s
# Awesome!
```

## Basic posts structure

As mentioned earlier, we need to create basic operations for `Post`, I prefer from top to bottom approach,
where we can define all routes and than return to authorization middleware and database.

Let's create `posts` folder inside `go-blog-ng` and add few files inside:

- `routers.go` - routers is main file for all routes regarding CRUD operation for `Posts`.
- `models.go` - basic models for `Post` and operation on it.
- `serializers.go` - helper for serialize/deserialize models to database types and vice versa.
- `validators.go` - just for validate post fields.
- `utils.go` - general helpers.

Now we can create and discuss each file in depth.


## Routers

A clean and elegant URL scheme is an important detail in a high-quality web application. Now we can start designing URLs.
Usually router is kind of URL dispatcher, who is responsible for mapping between URL path and expressions to Go functions.

Since we already created empty file `routers.go` inside `go-blog-ng/posts`, now you can review source code of it:

```go
package posts

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func PostsRegister(router *gin.RouterGroup) {
	router.GET("/", PostList)
	router.POST("/", PostCreate)
	router.GET("/:id", PostRetrieve)
	router.PUT("/:id", PostUpdate)
	router.DELETE("/:id", PostDelete)
}

func PostList(c *gin.Context) {
	authorEmail, _ := c.Get("user_email")
	articleModels, err := FindManyPost(authorEmail.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, NewError("posts", errors.New("invalid param")))
		return
	}
	serializer := PostsSerializer{c, articleModels}
	c.JSON(http.StatusOK, gin.H{"posts": serializer.Response()})
}

func PostRetrieve(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	postModel, err := FindOnePost(&PostModel{ID: uint(id)})
	if err != nil {
		c.JSON(http.StatusNotFound, NewError("posts", errors.New("invalid id")))
		return
	}
	serializer := PostSerializer{c, postModel}
	c.JSON(http.StatusOK, gin.H{"post": serializer.Response()})
}

func PostCreate(c *gin.Context) {
	articleModelValidator := NewPostModelValidator()
	if err := articleModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, NewValidatorError(err))
		return
	}
	if err := SaveOne(&articleModelValidator.postModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, NewError("database", err))
		return
	}
	serializer := PostSerializer{c, articleModelValidator.postModel}
	c.JSON(http.StatusCreated, gin.H{"article": serializer.Response()})
}

func PostUpdate(c *gin.Context) {
	id, _:= strconv.ParseUint(c.Param("id"), 10, 32)
	postModel, err := FindOnePost(&PostModel{ID: uint(id)})

	if err != nil {
		c.JSON(http.StatusNotFound, NewError("posts", errors.New("invalid slug")))
		return
	}

	postModelValidator := NewArticleModelValidatorFillWith(postModel)
	if err := postModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, NewValidatorError(err))
		return
	}

	postModelValidator.postModel.ID = postModel.ID
	if err := postModel.Update(postModelValidator.postModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, NewError("database", err))
		return
	}
	serializer := PostSerializer{c, postModel}
	c.JSON(http.StatusOK, gin.H{"post": serializer.Response()})
}

func PostDelete(c *gin.Context) {
	id, _:= strconv.ParseUint(c.Param("id"), 10, 32)
	err := DeletePostModel(&PostModel{ID: uint(id)})
	if err != nil {
		c.JSON(http.StatusNotFound, NewError("articles", errors.New("invalid id")))
		return
	}
	c.JSON(http.StatusOK, gin.H{"article": "Delete success"})
}
```

Main function here is `PostsRegister` where we can see all declared CRUD operations based on table we discussed earlier. Also if you need to get implementation details of each route, you can dig into each function.


## Models

A model is the one and complete source of information about all your data. It contains the all fields and behaviors of the data you’re storing. Generally, each model maps to a single database table.
In this tutorial I decided to demonstrate the fantastic ORM library for Go is [GORM](http://gorm.io/)

The basics are:

- Each model is a Go [struct](https://tour.golang.org/moretypes/2) that includes `gorm.Model`.
- Each field include [struct tags](https://flaviocopes.com/go-tags/) to define some metadata ``gorm:"type:varchar(100);unique_index"``

Now let's define all models and behaviour for mnipulating blog posts inside `models.go`:

```go
package posts

import (
	"github.com/andriisoldatenko/go-blog-ng/storage"
	"github.com/jinzhu/gorm"
)

type PostModel struct {
	ID          uint    `gorm:"primary_key"`
	Title       string
	Description string  `gorm:"size:2048"`
	Body        string  `gorm:"size:2048"`
	UserEmail   string  `gorm:"column:user_email"`
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(PostModel{})
}

func FindOnePost(condition interface{}) (PostModel, error) {
	db := storage.GetDB()
	var model PostModel
	tx := db.Begin()
	tx.Where(condition).First(&model)
	err := tx.Commit().Error
	return model, err
}

func FindManyPost(userEmail string) ([]PostModel, error) {
	db := storage.GetDB()
	var models []PostModel
	var err error

	tx := db.Begin()
	tx.Where(PostModel{UserEmail: userEmail}).Find(&models)
	err = tx.Commit().Error
	return models, err
}

func SaveOne(data interface{}) error {
	db := storage.GetDB()
	err := db.Save(data).Error
	return err
}

func DeletePostModel(condition interface{}) error {
	db := storage.GetDB()
	err := db.Where(condition).Delete(PostModel{}).Error
	return err
}

func (model *PostModel) Update(data interface{}) error {
	db := storage.GetDB()
	err := db.Model(model).Update(data).Error
	return err
}
```


## Migrations

```
go-blog-ng git:(master) ✗ go run main.go

(/Users/andrii/work/go-blog-ng/posts/models.go:17)
[2019-01-10 20:50:21]  [2.81ms]  CREATE TABLE "post_models" ("id" integer primary key autoincrement,"created_at" datetime,"updated_at" datetime,"deleted_at" datetime,"title" varchar(255),"description" varchar(2048),"body" varchar(2048),"author_id" integer )
[0 rows affected or returned ]

(/Users/andrii/work/go-blog-ng/posts/models.go:17)
[2019-01-10 20:50:21]  [0.86ms]  CREATE INDEX idx_post_models_deleted_at ON "post_models"(deleted_at)
[0 rows affected or returned ]
[GIN-debug] [WARNING] Now Gin requires Go 1.6 or later and Go 1.7 will be required soon.

[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)
```

## Serializers

Serializers helps to convert any custom type or content to []byte or string, for Go Programming Language and vice versa.

```go
package posts

import "github.com/gin-gonic/gin"

type PostSerializer struct {
	C *gin.Context
	PostModel
}

type PostResponse struct {
	ID             uint                  `json:"id"`
	Title          string                `json:"title"`
	Description    string                `json:"description"`
	Body           string                `json:"body"`
}

type PostsSerializer struct {
	C        *gin.Context
	Posts []PostModel
}

func (s *PostSerializer) Response() PostResponse {
	response := PostResponse{
		ID:          s.ID,
		Title:       s.Title,
		Description: s.Description,
		Body:        s.Body,
	}
	return response
}

func (s *PostsSerializer) Response() []PostResponse {
	var response []PostResponse
	for _, post := range s.Posts {
		serializer := PostSerializer{s.C, post}
		response = append(response, serializer.Response())
	}
	return response
}
```


## Middleware

```go
package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"os"
	"strings"
)

func isAuthenticated(c *gin.Context) bool {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		return false
	}
	tokenParts := strings.Split(authHeader, "Bearer ")
	bearerToken := tokenParts[1]

	tv := map[string]string{}
	tv["aud"] = "api://default"
	tv["cid"] = os.Getenv("SPA_CLIENT_ID")
	jv := verifier.JwtVerifier{
		Issuer:           os.Getenv("ISSUER"),
		ClaimsToValidate: tv,
	}

	myJwt, err := jv.New().VerifyAccessToken(bearerToken)
	if err != nil {
		return false
	}
	c.Set("user_email", myJwt.Claims["sub"])
	return true
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isAuthenticated(c) {
			err := errors.New("auth error")
			c.AbortWithError(401, err)
		}
	}
}
```

## Storage

```
package storage

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Opening a database and save the reference to `Database` struct.
func Init() *gorm.DB {
	db, err := gorm.Open("sqlite3", "/tmp/blog.db")
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
```

## Create frontend app using Angular CLI

Angular CLI is command line tool which can do all job for you: such as: ...

```
$ mkdir frontend && cd frontend
npm install -g @angular/cli
```

If you are not familar with Angular I higly recommend this tutorial https://angular.io/tutorial/toh-pt1.

