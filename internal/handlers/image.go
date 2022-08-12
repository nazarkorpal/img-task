package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nazarkorpal/img-task/internal/services"
)

type imageHandler struct {
	ImageService services.Image
}

type Image interface {
	routes(rg *gin.RouterGroup)
	AddImage() gin.HandlerFunc
}

func NewImageHandler(image services.Image) Image {
	return &imageHandler{
		ImageService: image,
	}
}

func (a *imageHandler) routes(rg *gin.RouterGroup) {
	r := rg.Group("images")
	r.POST("/", a.AddImage())
}

func (a *imageHandler) AddImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println(ctx.Request.FormFile("file"))
		file, fileHeader, err := ctx.Request.FormFile("image")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		defer file.Close()

		if err := a.ImageService.AddImage(file, fileHeader); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	}
}
