package handlers

import (
	// import the image formats even if we dont call them directly
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "image/png"

	"image"
	"image/jpeg"

	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dchote/survey-engine/config"
	"github.com/dchote/survey-engine/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

var (
	imageManager *models.ImageManager
)

// ImageObject struct
type ImageObject struct {
	ImageName string `json:"imageName"`
	URL       string `json:"url"`
}

// ImageList struct
type ImageList struct {
	Images []*ImageObject `json:"images"`
}

// InitImageHandler initialize the logic manager
func InitImageHandler(db *gorm.DB) {
	imageManager = models.NewImageManager(db)
}

// UploadImage method for handling the upload
func UploadImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		imageName := c.Param("name")

		if imageName == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Please specify image name")
		}

		// append .jpg extension
		imageName = imageName + ".jpg"

		// create upload directory
		uploadDir := config.Config.Files.Path
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			log.Printf("unable to create uploadDir: %s", uploadDir)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		file, err := c.FormFile("file")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		defer src.Close()

		// Destination
		destFilePath := filepath.Join(uploadDir, imageName)
		log.Printf("saving to file: %s", destFilePath)

		dst, err := os.Create(destFilePath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		defer dst.Close()

		img, _, err := image.Decode(src)
		if err != nil {
			log.Printf("error decoding image: %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		err = jpeg.Encode(dst, img, nil)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		// check for existing record
		image := imageManager.FindImage(imageName)
		if image.ImageName == "" {
			// create new if not
			// TODO: probably dont need full file path...
			image, _ = imageManager.NewImage(imageName, destFilePath)
		}

		imageObject := &ImageObject{
			ImageName: image.ImageName,
			URL:       "/images/" + image.ImageName,
		}

		return c.JSON(http.StatusOK, imageObject)
	}
}

// ListImages return all functions
func ListImages() echo.HandlerFunc {
	return func(c echo.Context) error {
		images := imageManager.AllImages()

		imageList := new(ImageList)

		for _, image := range images {
			imageObject := &ImageObject{
				ImageName: image.ImageName,
				URL:       "/images/" + image.ImageName,
			}
			imageList.Images = append(imageList.Images, imageObject)
		}

		return c.JSON(http.StatusOK, imageList)
	}
}

// FetchImage retrieve single function
func FetchImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		imageName := c.Param("name")

		if imageName == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Please specify image name")
		}

		image := imageManager.FindImage(imageName)

		if image.ImageName == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Image not found")
		}

		imageObject := &ImageObject{
			ImageName: image.ImageName,
			URL:       "/images/" + image.ImageName,
		}

		return c.JSON(http.StatusOK, imageObject)
	}
}

// DeleteImage delete function
func DeleteImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ValidateAdmin(c.Get("user").(*jwt.Token))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Permission Denied")
		}

		imageName := c.Param("name")

		image := imageManager.FindImage(imageName)

		if image.ImageName == "" {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		imageManager.DeleteImage(image)

		err = os.Remove(image.ImagePath)
		if err != nil {
			log.Printf("unable to remove file: %s", image.ImagePath)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, true)
	}
}

// ViewImage pass the image through
func ViewImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		imageName := c.Param("name")

		if imageName == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Please specify image name")
		}

		image := imageManager.FindImage(imageName)

		if image.ImageName == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Image not found")
		}

		fileName := filepath.Join(config.Config.Files.Path, image.ImageName)
		return c.File(fileName)
	}
}
