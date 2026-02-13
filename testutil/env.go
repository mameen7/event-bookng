package testutil

import "testing"

func SetupTestEnv(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-for-testing-only-do-not-use-in-production")
}
