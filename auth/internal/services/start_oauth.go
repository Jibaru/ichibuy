package services

import (
	"context"

	"golang.org/x/oauth2"
)

type StartOAuth struct {
	cfg *oauth2.Config
}

type StartOAuthResp struct {
	URL string
}

func NewStartOAuth(cfg *oauth2.Config) *StartOAuth {
	return &StartOAuth{cfg: cfg}
}

func (s *StartOAuth) Exec(ctx context.Context) (*StartOAuthResp, error) {
	url := s.cfg.AuthCodeURL("")
	return &StartOAuthResp{URL: url}, nil
}
