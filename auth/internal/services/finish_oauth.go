package services

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"

	"ichibuy/auth/config"
	"ichibuy/auth/internal/domain"
	"ichibuy/auth/internal/domain/dao"
)

type FinishOAuthReq struct {
	Code string `json:"code"`
}

type FinishOAuthResp struct {
	URL string
}

type FinishOAuth struct {
	userDAO            dao.UserDAO
	oauthCfg           *oauth2.Config
	oauthInfoExtractor domain.InfoExtractor
	nextID             domain.NextID
	cfg                config.Config
}

func NewFinishOAuth(
	userDAO dao.UserDAO,
	oauthCfg *oauth2.Config,
	oauthInfoExtractor domain.InfoExtractor,
	nextID domain.NextID,
	cfg config.Config,
) *FinishOAuth {
	return &FinishOAuth{
		userDAO:            userDAO,
		oauthCfg:           oauthCfg,
		oauthInfoExtractor: oauthInfoExtractor,
		nextID:             nextID,
		cfg:                cfg,
	}
}

func (s *FinishOAuth) Exec(ctx context.Context, req FinishOAuthReq) (*FinishOAuthResp, error) {
	token, err := s.oauthCfg.Exchange(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	username, email, err := s.oauthInfoExtractor(token.AccessToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userDAO.FindOne(ctx, "email = $1", "", email)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		newUser, err := domain.NewUser(s.nextID(), email, username)
		if err != nil {
			return nil, err
		}

		err = s.userDAO.Create(ctx, newUser)
		if err != nil {
			return nil, err
		}

		user = newUser
	} else if err != nil {
		return nil, err
	}

	privateKey, err := domain.ParseRSAPrivateKeyFromPEM(s.cfg.JWTPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
	}

	tokenString, err := generateToken(user.ID, user.Email, privateKey)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/auth/callback/google?token=%s&id=%s&email=%s&username=%s",
		s.cfg.WebBaseURI,
		tokenString,
		user.ID,
		user.Email,
		user.Username,
	)

	return &FinishOAuthResp{
		URL: url,
	}, nil
}

func generateToken(userID string, userEmail string, privateKey *rsa.PrivateKey) (string, error) {
	publicKeyBytes, err := domain.PublicKeyToPEM(&privateKey.PublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to convert public key to PEM: %w", err)
	}

	hash := sha256.Sum256([]byte(publicKeyBytes))
	kid := base64.URLEncoding.EncodeToString(hash[:8])

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_id": userID,
		"email":   userEmail,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     "ichibuy-auth",
	})

	token.Header["kid"] = kid

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
