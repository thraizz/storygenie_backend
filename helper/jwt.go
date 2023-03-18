package helper

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func VerifyJWT(ctx context.Context, idToken string) (jwt.MapClaims, error) {
	// If the token contains bearer, remove it
	if len(idToken) > 6 && idToken[:7] == "Bearer " {
		idToken = idToken[7:]
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %v", err)
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil || token == nil {
		return nil, fmt.Errorf("error verifying ID token: %v", err)
	}
	return token.Claims, err
}

func GetUserFromRequest(ctx *gin.Context) (string, error) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		return "", fmt.Errorf("no token provided")
	}
	claims, err := VerifyJWT(ctx, token)
	if err != nil {
		return "", err
	}
	log.Println(claims)
	uid := claims["user_id"]
	if uid == nil {
		return "", fmt.Errorf("no uid in token")
	}
	return uid.(string), nil
}
