package auth

import "context"

// custom type biar gak bentrok antar package
type contextKey string

const contextKeyUserID contextKey = "userID"

// helper biar handler gak akses key langsung
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(contextKeyUserID).(string)
	return id, ok
}
