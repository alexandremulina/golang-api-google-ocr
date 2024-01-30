package api

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

// HealthCheck
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OK",
		"status":  200,
	})
}

type FileContent struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

// API REQUEST METHOD
func Ocr(c *gin.Context) {
	var fileContent FileContent

	if err := c.ShouldBind(&fileContent); err != nil {
		c.JSON(400, gin.H{
			"message": "Bad Request",
			"status":  400,
		})
		return
	}

	// Authenticate with Google Cloud Vision API
	ctx := context.Background()
	credsPath := os.Getenv("CREDENTIALS_FILE_PATH")

	creds, err := os.ReadFile(credsPath) // Replace with the actual path to your credentials file
	if err != nil {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("Error reading credentials file: %v", err),
			"status":  500,
		})
		return
	}

	visionService, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsJSON(creds))
	if err != nil {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("Error creating Vision API client: %v", err),
			"status":  500,
		})
		return
	}

	// Get the file content
	file, err := fileContent.File.Open()
	if err != nil {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("Error opening file: %v", err),
			"status":  500,
		})
		return
	}
	defer file.Close()

	fileContentBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("Error reading file content: %v", err),
			"status":  500,
		})
		return
	}

	// Perform OCR on the image
	annotations, err := visionService.DetectTexts(ctx, &visionpb.Image{
		Content: fileContentBytes,
	}, &visionpb.ImageContext{}, 1)

	if err != nil {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("Error in DetectTexts: %v", err),
			"status":  500,
		})
		return
	}

	var ocrResults []string
	for _, annotation := range annotations {
		ocrResults = append(ocrResults, annotation.Description)
	}

	c.JSON(200, gin.H{
		"message": "OCR Successful",
		"status":  200,
		"results": ocrResults,
	})
}
