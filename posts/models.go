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