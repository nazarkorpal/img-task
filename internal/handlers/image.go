package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nazarkorpal/img-task/internal/services"
)

type imageHandler struct {
	ImageService services.Image
}

type Image interface {
	routes(rg *gin.RouterGroup)
	AddImage() gin.HandlerFunc
	GetImage() gin.HandlerFunc
}

func NewImageHandler(image services.Image) Image {
	return &imageHandler{
		ImageService: image,
	}
}

func (a *imageHandler) routes(rg *gin.RouterGroup) {
	r := rg.Group("images")
	r.POST("/", a.AddImage())
	r.GET("/:id", a.GetImage())

	go a.GeneratingLessQualityImages()
}

func (a *imageHandler) AddImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, fileHeader, err := ctx.Request.FormFile("image")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		defer file.Close()

		id, err := a.ImageService.AddImage(file, fileHeader)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"imageID": id,
		})
	}
}

func (a *imageHandler) GetImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param, ok := ctx.Params.Get("id")
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "there is not id of image",
			})
			return
		}

		quality, _ := ctx.GetQuery("quality")

		id, err := strconv.Atoi(param)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		fileName, filePath, err := a.ImageService.GetImage(uint(id), quality)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Content-Disposition", "inline;filename="+fileName)
		ctx.Header("Content-Transfer-Encoding", "binary")

		fmt.Println(filePath)

		ctx.File(filePath)
	}
}

func (a *imageHandler) GeneratingLessQualityImages() {
	if err := a.ImageService.GenerateLessQualityImages(); err != nil {
		panic(err)
	}
}
