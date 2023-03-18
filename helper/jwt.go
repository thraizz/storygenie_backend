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
		log.Fatalf("error getting Auth client: %v\n", err)
		return nil, err
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
		return nil, err
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
	return claims["uid"].(string), nil
}
