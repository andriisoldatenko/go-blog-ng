package posts

import (
	"github.com/gin-gonic/gin"
)

type PostModelValidator struct {
	Post struct {
		Title       string   `form:"title" json:"title" binding:"exists,min=4"`
		Description string   `form:"description" json:"description" binding:"max=2048"`
	} `json:"post"`
	postModel PostModel `json:"-"`
}

func NewPostModelValidator() PostModelValidator {
	return PostModelValidator{}
}

func NewPostModelValidatorFillWith(postModel PostModel) PostModelValidator {
	postModelValidator := NewPostModelValidator()
	postModelValidator.Post.Title = postModel.Title
	postModelValidator.Post.Description = postModel.Description
	return postModelValidator
}

func (s *PostModelValidator) Bind(c *gin.Context) error {
	s.postModel.Title = s.Post.Title
	s.postModel.Description = s.Post.Description
	s.postModel.UserEmail = c.MustGet("user_email").(string)
	return nil
}