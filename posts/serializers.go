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