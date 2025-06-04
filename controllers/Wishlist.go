package controllers

import (
	"strconv"
	"time"

	db "github.com/BassemArfaoui/Weazy-Server/config"
	"github.com/BassemArfaoui/Weazy-Server/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetWishlistByUserId(c *fiber.Ctx) error {
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

	var products []models.Product
	var totalCount int64

	//todo make it ordered
	result := db.DB.Raw(`
	SELECT DISTINCT p.*, wl.created_at 
	FROM products p
	JOIN wishlists wl ON p.id = wl.product_id
	WHERE wl.user_id = ?
	ORDER BY wl.created_at DESC
	LIMIT ? OFFSET ?
`, userId, limit, offset).Scan(&products)

	for i := range products {
		products[i].IsLiked = true
	}

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch wishlist",
		})
	}

	err = db.DB.Raw(`
	SELECT COUNT(*) FROM (
		SELECT DISTINCT product_id
		FROM wishlists
		WHERE user_id = ?
	) AS distinct_products
`, userId).Scan(&totalCount).Error





	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch total wishlist products count",
		})
	}

	hasNextPage := int64(offset+limit) < totalCount

	return c.Status(200).JSON(fiber.Map{
		"error":       false,
		"message":     "wishlist fetched successfully",
		"page":        page,
		"total":       totalCount,
		"hasNextPage": hasNextPage,
		"data":        products,
	})
}

func SaveToWishlist(c *fiber.Ctx) error {
	var wishlist models.Wishlist

	if err := c.BodyParser(&wishlist); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid wishlist data",
		})
	}

	if wishlist.UserId == uuid.Nil || wishlist.ProductId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "User_id and product_id are required",
		})
	}

	wishlist.CreatedAt = time.Now()

	if err := db.DB.Create(&wishlist).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to save wishlist item: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "wishlist item added successfully",
		"data":    wishlist,
	})
}

func DeleteWishlistItem(c *fiber.Ctx) error {

	productId := c.Params("productId")

	userId, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid userId format",
		})
	}

	var wishlist models.Wishlist

	result := db.DB.Where("user_id = ? and  product_id= ?", userId, productId).Delete(&wishlist)

	if result.Error != nil {

		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete wishlist item",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error":   false,
		"message": "wishlist item deleted successfully",
	})
}
