package middleware

import (
	"context"
	"inventory-system/pkg/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// ContextKey adalah custom type untuk key di context biar gak tabrakan sama package lain
type ContextKey string

const (
	UserIDKey   ContextKey = "user_id"
	UserRoleKey ContextKey = "user_role"
)

// ==========================================
// 1. AUTHENTICATION MIDDLEWARE
// ==========================================

// Authenticate mengecek keaslian JWT token dari header Authorization
func Authenticate(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Ambil token dari header "Authorization: Bearer <token>"
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.Error(w, r, http.StatusUnauthorized, "Missing authorization header", nil)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.Error(w, r, http.StatusUnauthorized, "Invalid authorization format", nil)
				return
			}

			tokenString := parts[1]

			// 2. Validasi dan Parse Token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Pastikan algoritma signing-nya sesuai
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, http.ErrAbortHandler
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				utils.Error(w, r, http.StatusUnauthorized, "Invalid or expired token", nil)
				return
			}

			// 3. Ambil data payload (claims) dari dalam token
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				utils.Error(w, r, http.StatusUnauthorized, "Invalid token claims", nil)
				return
			}

			userID := claims["user_id"].(string)
			role := claims["role"].(string)

			// 4. Simpan UserID dan Role ke dalam Context Request
			// Biar Handler yang dipanggil setelah ini bisa tau siapa yang lagi login
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserRoleKey, role)

			// Lanjut ke rute berikutnya dengan context yang udah diselipin data user
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ==========================================
// 2. AUTHORIZATION MIDDLEWARE (RBAC)
// ==========================================

// RequireRoles mengecek apakah user yang login punya role yang diizinkan
func RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Ambil role dari context (yang tadi dimasukin sama middleware Authenticate)
			userRole, ok := r.Context().Value(UserRoleKey).(string)
			if !ok || userRole == "" {
				utils.Error(w, r, http.StatusUnauthorized, "Unauthorized access", nil)
				return
			}

			// 2. Cek apakah role user ada di dalam daftar allowedRoles
			isAllowed := false
			for _, role := range allowedRoles {
				if userRole == role {
					isAllowed = true
					break
				}
			}

			// 3. Kalau gak ada yang cocok, tolak aksesnya (403 Forbidden)
			if !isAllowed {
				utils.Error(w, r, http.StatusForbidden, "You don't have permission to access this resource", nil)
				return
			}

			// Kalau aman, lanjut ke handler!
			next.ServeHTTP(w, r)
		})
	}
}
