package posts

import (
	"github.com/gin-gonic/gin"
)

type PostModelValidator struct {
	Post struct {
		Title       string   `form:"title" json:"title" binding:"exists,min=4"`
		Description string   `form:"description" json:"description" binding:"max=2048"`
		Body string          `form:"body" json:"body" binding:"max=2048"`
	} `json:"post"`
	postModel PostModel `json:"-"`
}

func NewPostModelValidator() PostModelValidator {
	return PostModelValidator{}
}

func NewArticleModelValidatorFillWith(postModel PostModel) PostModelValidator {
	postModelValidator := NewPostModelValidator()
	postModelValidator.Post.Title = postModel.Title
	postModelValidator.Post.Description = postModel.Description
	postModelValidator.Post.Body = postModel.Body
	return postModelValidator
}

func (s *PostModelValidator) Bind(c *gin.Context) error {
	err := Bind(c, s)
	if err != nil {
		return err
	}
	s.postModel.Title = s.Post.Title
	s.postModel.Description = s.Post.Description
	s.postModel.Body = s.Post.Body
	s.postModel.UserEmail = c.MustGet("user_email").(string)
	return nil
}