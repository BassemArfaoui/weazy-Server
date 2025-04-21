package controllers

import (
	db "github.com/BassemArfaoui/Weazy-Server/config"
	"github.com/BassemArfaoui/Weazy-Server/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strconv"
)

func GetChatsByUserId(c *fiber.Ctx) error {

	userId, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid userId format",
		})
	}

	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil || limit <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid limit value",
		})
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid page value",
		})
	}

	offset := (page - 1) * limit

	var chats []models.Chat
	var totalCount int64

	result := db.DB.Where("user_id = ?", userId).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&chats)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch chats",
		})
	}

	err = db.DB.Model(&models.Chat{}).Where("user_id = ?", userId).Count(&totalCount).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch total chat count",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error":   false,
		"message": "Chats fetched successfully",
		"page":    page,
		"total":   totalCount,
		"data":    chats,
	})
}

func EditChat(c *fiber.Ctx) error {

	chatId, err := uuid.Parse(c.Params("chatId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid chatId format",
		})
	}


	var chat models.Chat
	if err := c.BodyParser(&chat); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid chat data",
		})
	}

	result := db.DB.Model(&models.Chat{}).Where("id = ?", chatId).Updates(&chat).First(&chat)
	if result.Error != nil {
		if(result.Error.Error() == "record not found") {
			return c.Status(404).JSON(fiber.Map{
				"error":   true,
				"message": "Chat not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update chat",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error":   false,
		"message": "Chat updated successfully",
		"data":    chat,
	})
}

func DeleteChat(c *fiber.Ctx) error {
	chatId, err := uuid.Parse(c.Params("chatId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid chatId format",
		})
	}
	var chat models.Chat

	result := db.DB.Where("id = ?", chatId).Delete(&chat)


	if (result.Error != nil) {

		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete chat",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error":   false,
		"message": "Chat deleted successfully",
	})
}