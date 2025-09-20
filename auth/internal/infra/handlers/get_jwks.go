package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"

	"ichibuy/auth/config"
	"ichibuy/auth/internal/domain"
)

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

// GetJWKS godoc
// @Summary      GetJWKS
// @Description  Returns JSON Web Key Set for JWT verification
// @Accept       json
// @Produce      json
// @Success      200    {object}    JWKS
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/auth/.well-known/jwks.json [get]
func GetJWKS(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		privateKeyPEM := cfg.JWTPrivateKey
		if privateKeyPEM == "" {
			c.JSON(http.StatusInternalServerError, ErrorResp{Error: "JWT private key not configured"})
			return
		}

		privateKey, err := domain.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Error: fmt.Sprintf("Failed to parse RSA private key: %v", err)})
			return
		}

		publicKey := &privateKey.PublicKey
		publicKeyPEM, err := domain.PublicKeyToPEM(publicKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Error: fmt.Sprintf("Failed to convert public key to PEM: %v", err)})
			return
		}

		hash := sha256.Sum256([]byte(publicKeyPEM))
		kid := base64.URLEncoding.EncodeToString(hash[:8])

		nBytes := publicKey.N.Bytes()
		eBytes := big.NewInt(int64(publicKey.E)).Bytes()

		jwk := JWK{
			Kty: "RSA",
			Use: "sig",
			Kid: kid,
			Alg: "RS256",
			N:   base64.URLEncoding.EncodeToString(nBytes),
			E:   base64.URLEncoding.EncodeToString(eBytes),
		}

		jwks := JWKS{
			Keys: []JWK{jwk},
		}

		c.JSON(http.StatusOK, jwks)
	}
}
