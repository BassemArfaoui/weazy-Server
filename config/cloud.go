package config

import (
	"os"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
)

var CloudinaryClient *cloudinary.Cloudinary

func InitCloud() {
	var err error
	CloudinaryClient, err = cloudinary.NewFromParams(
		os.Getenv("CLOUD_NAME"),
		os.Getenv("API_KEY"),
		os.Getenv("API_SECRET"),
	)

	if err != nil {
		fmt.Println("Failed to initialize Cloud:", err)
	}
	fmt.Println("Cloud initialized")
}
