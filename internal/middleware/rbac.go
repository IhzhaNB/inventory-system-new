package middleware

import (
	"net/http"
	"slices"

	"inventory-system/pkg/utils"
)

// RequireRole restricts access to endpoints based on the user's role.
// It accepts a variadic list of allowed roles.
// NOTE: This middleware MUST be placed AFTER the Authenticate middleware.
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// 1. Extract the user role from the request context.
			userRole, ok := r.Context().Value(UserRoleKey).(string)
			if !ok || userRole == "" {
				utils.Error(w, r, http.StatusForbidden, "Access denied: Role identifier is missing", nil)
				return
			}

			// 2 & 3. Check if the role is allowed AND block if it's not (IN ONE LINE!)
			// slices.Contains bakal otomatis ngecek apakah userRole ada di dalam allowedRoles
			if !slices.Contains(allowedRoles, userRole) {
				utils.Error(w, r, http.StatusForbidden, "Access denied: Insufficient permissions", nil)
				return
			}

			// 4. Role is valid, proceed to the actual handler.
			next.ServeHTTP(w, r)
		})
	}
}
