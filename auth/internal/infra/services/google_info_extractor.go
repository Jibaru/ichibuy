package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GoogleInfoExtractor(token string) (string, string, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("error gtting user info, status: %s", resp.Status)
	}

	type GoogleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
	}

	var user GoogleUser
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", "", err
	}

	return user.Name, user.Email, nil
}
