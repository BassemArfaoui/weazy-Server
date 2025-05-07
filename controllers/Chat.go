package controllers

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	db "github.com/BassemArfaoui/Weazy-Server/config"
	"github.com/BassemArfaoui/Weazy-Server/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
		if result.Error.Error() == "record not found" {
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

	if result.Error != nil {

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

	if chat.UserId == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "UserId is required",
		})
	}

	if chat.Message == "" && len(chat.ImageURLs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Either a Message or at least one ImageURL is required",
		})
	}

	if chat.Title == "" {
		chat.Title = "Untitled"
	}

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
		ImageURLs:  chat.ImageURLs,
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

	userId := c.Params("userId")
	if userId == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "userId is required",
		})
	}

	query := `
		SELECT 
			m.id AS id,
			m.sender_role AS sender,
			m.text AS message,
			m.image_urls,
			m.created_at AS message_created_at,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', p.id,
						'gender', p.gender,
						'mastercategory', p.mastercategory,
						'subcategory', p.subcategory,
						'articletype', p.articletype,
						'basecolour', p.basecolour,
						'season', p.season,
						'year', p.year,
						'usage', p.usage,
						'productdisplayname', p.productdisplayname,
						'link', p.link,
						'is_liked', CASE WHEN w.product_id IS NOT NULL THEN true ELSE false END
					)
					ORDER BY pid.ordinality
				) FILTER (WHERE p.id IS NOT NULL),
				'[]'
			) AS products
		FROM messages m
		LEFT JOIN LATERAL UNNEST(m.products) WITH ORDINALITY AS pid(product_id, ordinality) ON TRUE
		LEFT JOIN products p ON p.id = pid.product_id
		LEFT JOIN wishlists w ON w.product_id = p.id AND w.user_id = ?
		WHERE m.chat_id = ?
		GROUP BY m.id
		ORDER BY m.created_at;
		`

	rows, err := db.DB.Raw(query, userId, chatId).Rows()
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
			message     string
			imageUrls   *string
			createdAt   time.Time
			productsRaw []byte
		)

		err := rows.Scan(&id, &sender, &message, &imageUrls, &createdAt, &productsRaw)
		if err != nil {
			continue
		}

		// Handle image URLs
		var urls []string
		if imageUrls != nil {
			cleaned := strings.Trim(*imageUrls, "{}")
			if cleaned != "" {
				urls = strings.Split(cleaned, ",")
			}
		}

		// Decode products JSON
		var products []models.Product
		if err := json.Unmarshal(productsRaw, &products); err != nil {
			products = []models.Product{}
		}

		messages = append(messages, map[string]interface{}{
			"id":         id,
			"sender":     sender,
			"message":    message,
			"image_urls": urls,
			"created_at": createdAt,
			"products":   products,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error":   false,
		"message": "Messages fetched successfully",
		"data":    messages,
	})
}
