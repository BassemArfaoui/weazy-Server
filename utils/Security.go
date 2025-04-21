package utils



import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)



func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}


func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}


func GenerateJWT(claims map[string]interface{}) (string, error) {

	jwtSecret := os.Getenv("JWT_SECRET")

	if _, exists := claims["exp"]; !exists {
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix() 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	signedToken, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
