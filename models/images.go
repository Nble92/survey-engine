package models

import (
	//"fmt"

	"github.com/jinzhu/gorm"
	// the blank import is the way GORM wants you to do it
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Image struct
type Image struct {
	gorm.Model `json:"-"`
	ImageName  string `gorm:"not null;unique" json:"imageName"`
	ImagePath  string `gorm:"not null" json:"imagePath"`
}

// ImageManager manager
type ImageManager struct {
	db *gorm.DB
}

// NewImageManager for handling the DB operations
func NewImageManager(db *gorm.DB) *ImageManager {
	imageManager := ImageManager{}
	imageManager.db = db

	imageManager.db.AutoMigrate(&Image{})

	return &imageManager
}

// Image

// FindImage simple search
func (imageManager *ImageManager) FindImage(name string) *Image {
	image := Image{}

	imageManager.db.Where("image_name=?", name).Find(&image)

	return &image
}

// AllImages return all stored functions
func (imageManager *ImageManager) AllImages() (images []Image) {
	err := imageManager.db.Find(&images).Error

	if err != nil {
		return images
	}

	return images
}

// NewImage create new function
func (imageManager *ImageManager) NewImage(imageName string, imagePath string) (image *Image, err error) {
	image = &Image{
		ImageName: imageName,
		ImagePath: imagePath,
	}

	err = imageManager.db.Create(&image).Error
	if err != nil {
		return nil, err
	}

	return image, nil
}

// DeleteImage delete a function
func (imageManager *ImageManager) DeleteImage(updateImage *Image) (err error) {
	err = imageManager.db.Delete(&updateImage).Error
	if err != nil {
		return err
	}

	return nil
}
