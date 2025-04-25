package controllers

import (
	"context"
	"github.com/BassemArfaoui/Weazy-Server/config"
	"github.com/BassemArfaoui/Weazy-Server/models"
	"github.com/BassemArfaoui/Weazy-Server/utils"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func BoolPtr(b bool) *bool {
	return &b
}

func Upload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to parse form",
		})
	}
	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "No file provided",
		})
	}
	file := files[0]

	if err := utils.ValidateImageFile(file); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	fileReader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to read file",
		})
	}
	defer fileReader.Close()

	timestamp := time.Now().Unix()
	uniquePublicID := strings.TrimSuffix(file.Filename, strings.ToLower(filepath.Ext(file.Filename))) + "_" + strconv.FormatInt(timestamp, 10)

	ctx := context.Background()
	uploadParams := uploader.UploadParams{
		Folder:       "user_uploads",
		PublicID:     uniquePublicID,
		Overwrite:    BoolPtr(false),
		ResourceType: "image",
	}

	resp, err := config.CloudinaryClient.Upload.Upload(ctx, fileReader, uploadParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to upload image to Cloudinary",
		})
	}

	response := models.UploadResponse{
		PublicID: resp.PublicID,
		URL:      resp.SecureURL,
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Picture uploaded successfully",
		"data":    response,
	})
}
