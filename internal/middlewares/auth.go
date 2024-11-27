package middlewares


import (
	"context"
	"net/http"
	"strings"
	"time"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

//todo need to use the .env for this
var jwtKey = []byte("your_secret_key")

// Claims structure
type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	// This adds standard fields like issuer, expirationtime and audience 
	jwt.RegisteredClaims
}

// Secure endpoint
func AuthMiddleware(next http.Handler)(http.Handler){
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// Extract token from authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer "){
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify token
		claims, err := VerifyJWT(tokenString)
		if err != nil{
			http.Error(w, "Unauthorized: " + err.Error(), http.StatusUnauthorized)
			return
		}

		// Attach claims to request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "role", claims.Role)
		r = r.WithContext(ctx)

		// Proceed with the next handler
		next.ServeHTTP(w, r)

	})
}

func GenerateJWT(userID int, role string)(string, error){
	// Sets the token’s validity to 24 hours.
	expirationTime := time.Now().Add(24* time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// Creates a token using the HS256 signing method and the claims
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Signs the token with the secret key, converting it into a string for transmission.
	return token.SignedString(jwtKey)
}

func VerifyJWT(tokenString string)(*Claims, error){
	claims := &Claims{}
	//Decode the JWT and checks the signature against jwtKey
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error){
		// Supplies the secret key for signature verification.
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
