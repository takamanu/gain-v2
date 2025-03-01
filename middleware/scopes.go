package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gain-v2/configs"
	"gain-v2/features/users"
	"gain-v2/helper"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ScopesMiddlewareInterface interface {
	ValidateScope(scopes []string) echo.MiddlewareFunc
	IsLoggedIn() echo.MiddlewareFunc
	IsSocketLoggedIn() echo.MiddlewareFunc
	GetTokenData() echo.MiddlewareFunc
	GainSpecialMiddleware() echo.MiddlewareFunc
}

type ScopesMiddleWare struct {
	jwt helper.JWTInterface
	us  users.UserServiceInterface
	cfg configs.ProgrammingConfig
}

func NewMiddleware(jwt helper.JWTInterface, us users.UserServiceInterface, cfg configs.ProgrammingConfig) ScopesMiddlewareInterface {
	return &ScopesMiddleWare{
		jwt: jwt,
		us:  us,
		cfg: cfg,
	}
}

const (
	SecretKey = "secrettest" // Replace with your actual secret key
)

// var IsLoggedInUser = echojwt.WithConfig(echojwt.Config{SigningKey: sm.getSecretKey()})

func hasPermission(scopes, permissions []string) bool {
	permSet := make(map[string]struct{}, len(permissions))
	for _, perm := range permissions {
		permSet[perm] = struct{}{}
	}

	for _, scope := range scopes {
		if _, exists := permSet[scope]; exists {
			return true
		}
	}
	return false
}

func (sm *ScopesMiddleWare) ValidateScope(scopes []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// permissions, err := sm.us.GetAllPermissionsByUserID(uint(c.Get(USER_ID).(int)))

			// if err != nil {
			// 	return c.JSON(http.StatusUnauthorized, helper.FormatResponse(false, "unauthorized: unknown auth token", nil))
			// }

			// scopes = []string{"pasti.lu.ga.punya.ini"}
			// fmt.Println("Check permissions: ", permissions)
			// fmt.Println("Check Scopes: ", scopes)

			// if !hasPermission(scopes, permissions) {
			// 	return c.JSON(http.StatusUnauthorized, helper.FormatResponse(false, "unauthorized: you dont have permit to access this feature", nil))
			// }
			// c.Set(USER_ID, getID)

			return next(c)
		}
	}
}

func (sm *ScopesMiddleWare) IsLoggedIn() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			getID, err := sm.jwt.GetID(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, helper.FormatResponse(false, "unauthorized: unknown auth token", nil))
			}

			authHeader := c.Request().Header.Get("Authorization")

			c.Set(USER_ID, getID)
			c.Set(BEARER_TOKEN, authHeader[7:])

			_, err = sm.us.UpdateProfile(int(getID), users.UpdateProfile{
				LastLoginDate: time.Now(),
				IP:            c.RealIP(),
				UserAgent:     c.Request().UserAgent(),
			})

			if err != nil {
				fmt.Println("err: unable to update user last active, ", err.Error())
			}

			return next(c)
		}
	}
}

func (sm *ScopesMiddleWare) IsSocketLoggedIn() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			// var (
			// 	token  string
			// 	userID uint
			// )

			// authHeader := c.Request().Header.Get("Authorization")

			// // Jika Authorization tidak ada, gunakan Sec-WebSocket-Protocol
			// if authHeader == "" {
			// 	token = c.Request().Header.Get("Sec-WebSocket-Protocol")

			// 	userID, err = sm.jwt.GetIDFromStringToken("Bearer " + token)
			// 	if err != nil {
			// 		return c.JSON(http.StatusUnauthorized, helper.FormatResponse(false, "unauthorized: unknown auth token", nil))
			// 	}

			// 	authHeader = "Bearer " + token
			// } else {
			// 	userID, err = sm.jwt.GetIDFromStringToken(authHeader)
			// 	if err != nil {
			// 		return c.JSON(http.StatusUnauthorized, helper.FormatResponse(false, "unauthorized: unknown auth token", nil))
			// 	}
			// }

			// c.Set(USER_ID, int(userID))
			// c.Set(BEARER_TOKEN, authHeader[7:])

			return next(c)
		}
	}
}

func (sm *ScopesMiddleWare) GetTokenData() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authHeader := c.Request().Header.Get("Authorization")

			c.Set(BEARER_TOKEN, authHeader[7:])

			return next(c)
		}
	}
}

func generateSignature(clientID, requestID, requestTimestamp, requestTarget, digest string) string {
	// Create the signature string
	signatureString := fmt.Sprintf(
		"Client-Id:%s\nRequest-Id:%s\nRequest-Timestamp:%s\nRequest-Target:%s\nDigest:%s",
		clientID, requestID, requestTimestamp, requestTarget, digest,
	)

	// Generate HMAC-SHA256
	h := hmac.New(sha256.New, []byte(SecretKey))
	h.Write([]byte(signatureString))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return "HMACSHA256=" + signature
}

func (sm *ScopesMiddleWare) GainSpecialMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			clientID := c.Request().Header.Get("Client-Id")
			requestID := c.Request().Header.Get("Request-Id")
			requestTimestamp := c.Request().Header.Get("Request-Timestamp")
			signature := c.Request().Header.Get("Signature")

			if clientID == "" || requestID == "" || requestTimestamp == "" || signature == "" {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "err: header is not complete", nil))
			}

			reqTime, err := time.Parse(time.RFC3339, requestTimestamp)
			if err != nil || time.Since(reqTime).Minutes() > 5 {
				return c.JSON(http.StatusUnauthorized, helper.FormatResponse(false, "err: max timestamp of request is 5 minutes before now", nil))
			}

			reqTarget := c.Request().URL.Path

			body, _ := io.ReadAll(c.Request().Body)

			// Restore the request body for future reads
			c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

			fmt.Println("Raw Body:", string(body))

			// If body is empty or null, set it to "{}"

			var jsonBody map[string]interface{}
			if err := json.Unmarshal(body, &jsonBody); err != nil {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "err: invalid req body", nil))
			}

			// Re-marshal to ensure consistent JSON formatting
			normalizedBody, _ := json.Marshal(jsonBody)

			// Generate SHA-256 hash and Base64 encode it
			hash := sha256.Sum256(normalizedBody)
			digest := base64.StdEncoding.EncodeToString(hash[:])

			fmt.Println("Digest:", digest)

			// Generate expected signature
			expectedSignature := generateSignature(clientID, requestID, requestTimestamp, reqTarget, digest)

			if !hmac.Equal([]byte(expectedSignature), []byte(signature)) {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "err: inavlid signature", nil))
			}

			return next(c)
			// return c.JSON(http.StatusOK, map[string]string{"message": "Authorized"})

		}
	}
}
