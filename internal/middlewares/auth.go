package middlewares


import (
	"context"
	"net/http"
	"strings"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

// Claims structure
type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler)(http.Handler){
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// Extract tpken from authorization header
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

func VerifyJWT(tokenString string)(*Claims, error){

}
