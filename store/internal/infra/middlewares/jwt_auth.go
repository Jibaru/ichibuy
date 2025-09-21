package middlewares

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"ichibuy/store/internal/infra/client"
)

type JWTAuthMiddleware struct {
	authClient *client.AuthClient
	cache      map[string]*rsa.PublicKey
}

func NewJWTAuthMiddleware(authClient *client.AuthClient) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		authClient: authClient,
		cache:      make(map[string]*rsa.PublicKey),
	}
}

func (m *JWTAuthMiddleware) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("kid not found in token header")
			}

			publicKey, err := m.getPublicKey(kid)
			if err != nil {
				return nil, err
			}

			return publicKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
			c.Abort()
			return
		}

		slog.InfoContext(c, "token validated successfully", "user_id", userID)
		c.Set("user_id", userID)
		c.Next()
	}
}

func (m *JWTAuthMiddleware) getPublicKey(kid string) (*rsa.PublicKey, error) {
	if cachedKey, exists := m.cache[kid]; exists {
		return cachedKey, nil
	}

	jwks, err := m.authClient.GetJWKS()
	if err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			publicKey, err := m.jwkToPublicKey(key)
			if err != nil {
				return nil, err
			}

			m.cache[kid] = publicKey
			return publicKey, nil
		}
	}

	return nil, fmt.Errorf("public key not found for kid: %s", kid)
}

func (m *JWTAuthMiddleware) jwkToPublicKey(jwk client.JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.URLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode n: %w", err)
	}

	eBytes, err := base64.URLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode e: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	publicKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return publicKey, nil
}
