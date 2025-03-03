package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims struct
type JWTClaims struct {
	UserID        string `json:"user_id"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	EntityID      string `json:"entity_id"`
	WalletAddress string `json:"wallet_address"` // ✅ เพิ่ม WalletAddress
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token
func GenerateToken(userID, email, role, entityID, walletAddress string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // ✅ อายุ 1 วัน

	claims := &JWTClaims{
		UserID:        userID,
		Email:         email,
		Role:          role,
		EntityID:      entityID,
		WalletAddress: walletAddress, // ✅ เพิ่ม WalletAddress เข้าไปใน Token
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "supplychain_backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// ✅ Debug ตรวจสอบค่า UserID, Email, Role และ WalletAddress ก่อนสร้าง Token
	fmt.Println("🛠 [GenerateToken] Creating token for User ID:", userID, "Email:", email, "Role:", role, "Wallet:", walletAddress)

	return signedToken, nil
}

// ValidateToken from Cookie
func ValidateToken(tokenString string) (*JWTClaims, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// ✅ Debug ค่า UserID, Email, Role, EntityID, และ WalletAddress
	fmt.Println("🔍 [ValidateToken] Extracted - User ID:", claims.UserID,
		"Email:", claims.Email,
		"Role:", claims.Role,
		"EntityID:", claims.EntityID,
		"Wallet:", claims.WalletAddress)

	return claims, nil
}

// AuthMiddleware
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ✅ ดึง Token จาก Cookie
		tokenString := c.Cookies("auth_token")
		fmt.Println("🔍 [AuthMiddleware] Received Token from Cookie:", tokenString)

		if tokenString == "" {
			fmt.Println("❌ [AuthMiddleware] No token found in Cookie")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization token required"})
		}

		// ✅ Validate Token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			fmt.Println("❌ [AuthMiddleware] Invalid token:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// ✅ Debug ตรวจสอบค่า UserID และ WalletAddress
		fmt.Println("✅ [AuthMiddleware] Authenticated User ID:", claims.UserID, "Wallet:", claims.WalletAddress)

		// ✅ Store claims in context
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		c.Locals("entityID", claims.EntityID)
		c.Locals("walletAddress", claims.WalletAddress) // ✅ เพิ่ม WalletAddress เข้าไปใน Context

		return c.Next()
	}
}
