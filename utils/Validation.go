package utils

import (
	"errors"
	"mime/multipart" 
	"path/filepath"
	"regexp"
	"strings"
	// db "github.com/BassemArfaoui/Quinsat-Server-Side/config"
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("weak Password : password must be at least 8 characters long")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("weak Password : password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.New("weak Password : password must contain at least one lowercase letter")
	}
	if !regexp.MustCompile(`\d`).MatchString(password) {
		return errors.New("weak Password : password must contain at least one number")
	}
	return nil
}

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters long")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}
	return nil
}

func ValidateEmail(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if !regexp.MustCompile(emailRegex).MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidateImageFile checks if the file is a valid image and within size limits
func ValidateImageFile(file *multipart.FileHeader) error {
	// Check for empty filename
	if file.Filename == "" {
		return errors.New("filename is empty")
	}

	// Check file extension
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedTypes[ext] {
		return errors.New("only JPG, JPEG, PNG, or GIF files are allowed")
	}

	// Check file size (max 10MB)
	const MaxFileSize = 10 * 1024 * 1024
	if file.Size > MaxFileSize {
		return errors.New("file size exceeds 10MB limit")
	}

	return nil
}


// func EmailExists(email string) (bool, error) {

// 	var student models.Student

// 	result := db.DB.Where("email = ?", email).First(&student)

// 	if result.Error != nil {
// 		if result.Error.Error() == "record not found" {
// 		return false, nil
// 		}
// 	return false, result.Error
// 	}

// 	return true, nil
// }
