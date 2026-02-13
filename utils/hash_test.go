package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword_GeneratesDifferentHashes(t *testing.T) {
	password := "mypassword123"

	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2, "Hashes should be different due to bcrypt salting")
}

func TestHashPassword_Success(t *testing.T) {
	password := "validpassword123"

	hash, err := HashPassword(password)

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash, "Hash should be different from plaintext password")
	assert.Greater(t, len(hash), len(password), "Hash should be longer than password")
}

func TestHashPassword_TooShort(t *testing.T) {
	password := "short"

	hash, err := HashPassword(password)

	require.Error(t, err)
	assert.Equal(t, "", hash)
	assert.Equal(t, "Password must be at least 8 characters", err.Error())
}

func TestCheckPasswordHash_ValidPassword(t *testing.T) {
	password := "mypassword123"
	hash, _ := HashPassword(password)

	result := CheckPasswordHash(password, hash)

	assert.True(t, result)
}

func TestCheckPasswordHash_InvalidPassword(t *testing.T) {
	password := "mypassword123"
	wrongPassword := "wrongpassword"
	hash, _ := HashPassword(password)

	result := CheckPasswordHash(wrongPassword, hash)

	assert.False(t, result)
}

func TestCheckPasswordHash_EmptyPassword(t *testing.T) {
	password := "mypassword123"
	hash, _ := HashPassword(password)

	result := CheckPasswordHash("", hash)

	assert.False(t, result)
}

func TestCheckPasswordHash_EmptyHash(t *testing.T) {
	password := "mypassword123"

	result := CheckPasswordHash(password, "")

	assert.False(t, result)
}

func TestCheckPasswordHash_InvalidHash(t *testing.T) {
	password := "mypassword123"

	result := CheckPasswordHash(password, "not-a-valid-bcrypt-hash")

	assert.False(t, result)
}
