---
layout: blog_post
title: "Build a Simple CRUD App with Go and Angular"
author: a_soldatenko
description: "This tutorial explains how to build a modern CRUD application with a Go as backend and an Angular frontend."
tags: [go, angular, crud]
---

# Build a Simple CRUD App with Go and Angular

## Instead of introduction

Go was invented in 2009 as a replacement of C by Robert Griesemer, Rob Pike, and Ken Thompson to solve system
programming real problems. Nowadays we can see huge amount of usage of Go in web, cloud native and even web for client
[WebAssembly](https://github.com/golang/go/wiki/WebAssembly#introduction).
But today I would like to share experience how to build modern blog application, where you
can create, read, update and delete blog posts.

## CRUD acronym

CRUD refer to SQL equivalents:

| Action        | SQL statement | HTTP Verb | URL Path    |
| ------------- |---------------| ----------|-------------|
| Create        | INSERT        | POST      | /posts      |
| Read          | SELECT        | GET       | /posts      |
| Update        | UPDATE        | PUT       | /post/<id>  |
| Delete        | DELETE        | DELETE    | /post/<id>  |


## Dependency management in Go

Today is't complicated to imagine modern web application without any dependency. By the way it's huge holy war in Go community
about package management and especially compare to Rust community. But I'm nit going to start discussing it here, I just
show you latest available in Go 1.11 [go modules](https://github.com/golang/go/wiki/Modules) from box which was proposed
in [Proposal: Versioned Go Modules](https://go.googlesource.com/proposal/+/master/design/24301-versioned-go.md)
by Russ Cox.


## Let's start from init project layout

I prefer to start from something real simple, and than step by step add / refactor code and structure.
But before we start let's test environment and installation. Assume you know how to install Go. If not please follow
this [getting started](https://golang.org/doc/install) instruction.

0. Go version

```bash
$ go version
go version go1.11.4 darwin/amd64
```

1. Create project structure:

```bash
cd $HOME/work/
tree -L 3 go-blog-ng
go-blog-ng
├── go.mod
└── main.go
```

2. Create you test `main.go` program:

```bash
cat main.go

package main

import "fmt"

func main() {
	fmt.Printf("Hello, Okta!\n")
}
```

3. Build and run

```bash
go run main.go
Hello, Okta!
```

As you can understand I need to install some web-framework, because I don't want to reinvent the wheel.
I prefer [Gin](https://github.com/gin-gonic/gin) http web framework written in Go for performance and productivity.


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

## Basic router structure

As mentioned earlier, we need to create basic operations for `Post`, I prefer from top to bottom approach,
where we can define all routes and than return to authorization middleware and database.

## Create frontend app using Angular CLI

Angular CLI is command line tool which can do all job for you: such as: ...

```
$ mkdir frontend && cd frontend
npm install -g @angular/cli
```

https://angular.io/tutorial/toh-pt1



## Migrate DB

```
go-blog-ng git:(master) ✗ go run main.go

(/Users/andrii/work/go-blog-ng-private/posts/models.go:17)
[2019-01-10 20:50:21]  [2.81ms]  CREATE TABLE "post_models" ("id" integer primary key autoincrement,"created_at" datetime,"updated_at" datetime,"deleted_at" datetime,"title" varchar(255),"description" varchar(2048),"body" varchar(2048),"author_id" integer )
[0 rows affected or returned ]

(/Users/andrii/work/go-blog-ng-private/posts/models.go:17)
[2019-01-10 20:50:21]  [0.86ms]  CREATE INDEX idx_post_models_deleted_at ON "post_models"(deleted_at)
[0 rows affected or returned ]
[GIN-debug] [WARNING] Now Gin requires Go 1.6 or later and Go 1.7 will be required soon.

[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)
```