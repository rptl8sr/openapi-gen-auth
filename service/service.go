package service

import (
	"context"
)

func YourOwnGetOrCreateUser(_ context.Context, username, _ string) (string, error) {
	// Hash password
	// Call DB with ctx
	// Check hashed password with stored hash
	// return id/uid/something or custom error
	return username, nil
}
