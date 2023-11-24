package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Trunk: if true, return only the username part of the email address
func extractAndValidateToken(r *http.Request, trunk bool) (string, error) {
	token := extractToken(r)
	if token == "" {
		return "", fmt.Errorf("no token found")
	}
	email, err := validateToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	if trunk {
		// since the username is an email address, we should cut off the domain part
		emailParts := strings.Split(email, "@")
		if len(emailParts) != 2 {
			return "", fmt.Errorf("invalid email address: %s", email)
		}
		return emailParts[0], nil
	}
	return email, nil
}

// extractToken extracts the token from the Authorization header
func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	if bearToken == "" {
		return ""
	}
	strArr := strings.Split(bearToken, " ")
	if len(strArr) != 2 {
		return ""
	}
	return strArr[1]
}

// validateToken validates the token with the auth service
func validateToken(token string) (string, error) {
	type TokenRequest struct {
		Token string `json:"token"`
	}
	tokenReq := TokenRequest{Token: token}
	tokenReqBytes, err := json.Marshal(tokenReq)
	if err != nil {
		return "", err
	}

	// Send the token to the auth service
	resp, err := http.Post("http://auth-svc:8080/auth/validate", "application/json", bytes.NewBuffer(tokenReqBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read the response body for an error message
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("invalid token; also failed to read error message: %v", err)
		}
		return "", fmt.Errorf("invalid token; response from auth service: %s, status: %s", string(body), resp.Status)

	}

	type authResponse struct {
		Email string  `json:"email"`
		Exp   float64 `json:"exp"`
	}

	var respData authResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}
	return respData.Email, nil

}
