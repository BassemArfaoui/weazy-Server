package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	db "github.com/BassemArfaoui/Weazy-Server/config"
	"github.com/BassemArfaoui/Weazy-Server/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetChatsByUserId(c *fiber.Ctx) error {
	userId , err := uuid.Parse(c.Params("userId"))
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
		Order("created_at DESC , id DESC").
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

	hasNextPage := int64(offset+limit) < totalCount

	return c.Status(200).JSON(fiber.Map{
		"error":       false,
		"message":     "Chats fetched successfully",
		"page":        page,
		"total":       totalCount,
		"hasNextPage": hasNextPage,
		"data":        chats,
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

func CreateChat(c *fiber.Ctx) error {
	var chat models.Chat

	if err := c.BodyParser(&chat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid chat data",
		})
	}


	if chat.UserId == uuid.Nil || chat.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "UserId and Message are required",
		})
	}

	if chat.Title == "" {
		chat.Title = "Untitled"
	}
	fmt.Println((chat.ImageURLs))



	fmt.Println((chat.ImageURLs))


	chat.Id = uuid.New()
	chat.CreatedAt = time.Now()

	if err := db.DB.Create(&chat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create chat: " + err.Error(),
		})
	}

	message := models.Message{
		Id:         uuid.New(),
		ChatId:     chat.Id,
		SenderRole: "user",
		Text:       chat.Message,
		ImageURLs : chat.ImageURLs,
		CreatedAt:  time.Now(),
	}

	if err := db.DB.Create(&message).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create initial message: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Chat and initial message created successfully",
		"data":    chat,
	})
}



func GetChatById(c *fiber.Ctx) error {
	chatId := c.Params("chatId")
	if chatId == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "chatId is required",
		})
	}

	query := `
	SELECT
		messages.id AS id,
		messages.sender_role as sender,
		messages.text as message,
		messages.image_urls,
		messages.created_at AS message_created_at
	FROM messages
	WHERE messages.chat_id = ?
	ORDER BY messages.created_at
	`

	rows, err := db.DB.Raw(query, chatId).Rows()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to execute query",
		})
	}
	defer rows.Close()

	var messages []map[string]interface{}

	for rows.Next() {
		var (
			id          string
			sender      string
			message       string
			imageUrls   *string
			createdAt   time.Time
		)

		err := rows.Scan(&id, &sender, &message, &imageUrls, &createdAt)
		if err != nil {
			continue
		}

		var urls []string
		if imageUrls != nil {
			cleaned := strings.Trim(*imageUrls, "{}")
			if cleaned != "" {
				urls = strings.Split(cleaned, ",")
			}
		}

		messages = append(messages, map[string]interface{}{
			"id":          id,
			"sender":      sender,
			"message":     message,
			"image_urls":  urls,
			"created_at":  createdAt,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error":   false,
		"message": "Messages fetched successfully",
		"data":    messages,
	})
}
