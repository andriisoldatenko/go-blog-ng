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