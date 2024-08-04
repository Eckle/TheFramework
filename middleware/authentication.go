package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Eckle/TheFramework/utils/token"
	"github.com/golang-jwt/jwt/v5"
)

type Key string

func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "token is required")
			return
		}

		tok := strings.Split(authorization, " ")
		if tok[0] != "Bearer" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "token type should be bearer")
			return
		}

		claims := jwt.MapClaims{
			"nbf": time.Now().Unix(),
		}

		tokenClaims, valid := token.Verify(tok[1], &claims)
		if !valid {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), Key("claims"), tokenClaims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
