package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSercet = []byte("secret_key")

var req struct {
	Password string `json:"password"`
}

var todoPassword string

func passwordHash(pass string) string {
	sum := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(sum[:])
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	if req.Password != todoPassword {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Неверный пароль"})
		return
	}

	claims := jwt.MapClaims{"hash": passwordHash(todoPassword)}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(jwtSercet)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "не удалось сформировать токен"})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{"token": signed})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(todoPassword) > 0 {
			var tokenStr string
			if cookie, err := r.Cookie("token"); err != nil {
				tokenStr = cookie.Value
			}

			valid := false

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				return jwtSercet, nil
			})
			if err == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if claims["hash"] == passwordHash(todoPassword) {
						valid = true
					}
				}
			}

			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}
