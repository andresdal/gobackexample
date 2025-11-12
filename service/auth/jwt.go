package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/andresdal/gobackexample/config"
	"github.com/andresdal/gobackexample/types"
	"github.com/andresdal/gobackexample/utils"
	"github.com/dgrijalva/jwt-go"
)

type contextKey string

const UserKey contextKey = "userID"

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc { // necesitamos el UserStore para obtener el usuario
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header in the request
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		// Validate JWT
		validatedToken, err := validateToken(token)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !validatedToken.Valid {
			log.Printf("invalid token")
			permissionDenied(w)
			return
		}

		// fetch user from DB using userID from token claims
		claims := validatedToken.Claims.(jwt.MapClaims)
		userIDStr, ok := claims["userID"].(string)
		if !ok {
			log.Printf("invalid token claims")
			permissionDenied(w)
			return
		}

		userID, _ := strconv.Atoi(userIDStr)

		user, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user from DB: %v", err)
			permissionDenied(w)
			return
		}

		// set context "userID" to the user ID
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}
