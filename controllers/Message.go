package controllers

import (
	db "github.com/BassemArfaoui/Weazy-Server/config"
	"github.com/BassemArfaoui/Weazy-Server/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

func SavePayload(c *fiber.Ctx) error {
	var payload models.MessagePayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid message payload",
		})
	}

	if payload.Request.ChatId == uuid.Nil ||
		(payload.Request.Text == "" && len(payload.Request.ImageURLs) == 0) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Request must include chatId and text or image_urls",
		})
	}

	if payload.Response.ChatId == uuid.Nil ||
		(payload.Response.Text == "" && len(payload.Response.ImageURLs) == 0 && len(payload.Response.Products) == 0) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Response must include chatId and at least one of text, image_urls, or products",
		})
	}

	// Generate timestamps
	reqTime := time.Now()
	respTime := reqTime.Add(time.Second)

	payload.Request.CreatedAt = reqTime
	payload.Response.CreatedAt = respTime

	// Save request
	if err := db.DB.Create(&payload.Request).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to save request message: " + err.Error(),
		})
	}

	// Save response
	if err := db.DB.Create(&payload.Response).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to save response message: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Request and response saved successfully",
		"data": fiber.Map{
			"request":  payload.Request,
			"response": payload.Response,
		},
	})
}

func SaveResponse(c *fiber.Ctx) error {
	var resp models.Message

	if err := c.BodyParser(&resp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid response payload: " + err.Error(),
		})
	}

	if resp.ChatId == uuid.Nil ||
		(resp.Text == "" && len(resp.ImageURLs) == 0 && len(resp.Products) == 0) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Response must include chat_id and at least one of text, image_urls, or products",
		})
	}

	resp.CreatedAt = time.Now()

	if err := db.DB.Create(&resp).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to save response message: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Response saved successfully",
		"data":    resp,
	})
}
