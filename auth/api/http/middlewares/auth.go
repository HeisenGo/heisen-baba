package middlewares

import (
    "net/http"
    "strings"

    "github.com/your-repo/auth/pkg/jwt"
    "github.com/your-repo/auth/pkg/rbac"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing authorization header", http.StatusUnauthorized)
            return
        }

        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 {
            http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
            return
        }

        claims, err := jwt.ValidateToken(bearerToken[1])
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Add user information to the request context
        ctx := r.Context()
        ctx = context.WithValue(ctx, "user_id", claims.UserID)
        ctx = context.WithValue(ctx, "roles", claims.Roles)

        next.ServeHTTP(w, r.WithContext(ctx))
    }
}

func RBACMiddleware(permission string) func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            roles, ok := r.Context().Value("roles").([]string)
            if !ok {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            if !rbac.IsGranted(roles, permission) {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        }
    }
}
