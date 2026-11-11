package middlewares

import (
	"aaa_bekend/internal/auth"
	"aaa_bekend/internal/lmodel"
	"net/http"
	"strings"
)

func AuthN(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization") // ambil header authorization

		// TODO
		// menangani header kosong
		if strings.TrimSpace(authHeader) == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// mengecek header yang formatnya salah
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// TODO
		// mengambil string JWT token
		tokenStr := strings.TrimSpace(parts[1])

		// validasi JWT token
		claims, err := auth.ValidateToken(tokenStr)

		// kasus error
		if err != nil {
			http.Error(w, "Invalid or expired JWT token", http.StatusUnauthorized)
			return
		}

		// kasus berhasil & pembuatan UserContext
		userCtx := lmodel.UserContext{
			ID:       claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		}

		ctx := lmodel.ContextWithUser(r.Context(), userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
