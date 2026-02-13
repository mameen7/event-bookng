package utils

import (
	"testing"
	"time"

	"event-booking/testutil"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken_CreatesValidToken(t *testing.T) {
	testutil.SetupTestEnv(t)

	email := "test@example.com"
	userId := int64(123)

	token, err := GenerateToken(email, userId)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Contains(t, token, ".", "JWT should contain dots")
}

func TestGenerateToken_ContainsCorrectClaims(t *testing.T) {
	testutil.SetupTestEnv(t)

	email := "test@example.com"
	userId := int64(456)

	token, err := GenerateToken(email, userId)
	require.NoError(t, err)

	extractedUserId, verifyErr := VerifyToken(&token)
	require.NoError(t, verifyErr)
	assert.Equal(t, userId, extractedUserId)
}

func TestGenerateToken_HasExpiration(t *testing.T) {
	testutil.SetupTestEnv(t)

	email := "test@example.com"
	userId := int64(789)

	token, err := GenerateToken(email, userId)
	require.NoError(t, err)

	// Parse token to check expiration
	parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret-key-for-testing-only-do-not-use-in-production"), nil
	})

	claims := parsedToken.Claims.(jwt.MapClaims)
	exp := int64(claims["exp"].(float64))
	expirationTime := time.Unix(exp, 0)

	expectedExpiration := time.Now().Add(2 * time.Hour)
	timeDiff := expirationTime.Sub(expectedExpiration).Abs()

	assert.Less(t, timeDiff, 10*time.Second, "Expiration should be approximately 2 hours from now")
}

func TestVerifyToken_ValidToken(t *testing.T) {
	testutil.SetupTestEnv(t)

	email := "test@example.com"
	userId := int64(123)
	token, _ := GenerateToken(email, userId)

	extractedUserId, err := VerifyToken(&token)

	require.NoError(t, err)
	assert.Equal(t, userId, extractedUserId)
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	testutil.SetupTestEnv(t)

	invalidToken := "invalid.token.here"

	userId, err := VerifyToken(&invalidToken)

	require.Error(t, err)
	assert.Equal(t, int64(0), userId)
}

func TestVerifyToken_MalformedToken(t *testing.T) {
	testutil.SetupTestEnv(t)

	malformedToken := "notajwt"

	userId, err := VerifyToken(&malformedToken)

	require.Error(t, err)
	assert.Equal(t, int64(0), userId)
}

func TestVerifyToken_WrongSigningMethod(t *testing.T) {
	testutil.SetupTestEnv(t)

	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"email":  "test@example.com",
		"userId": float64(123),
		"exp":    time.Now().Add(2 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	userId, err := VerifyToken(&tokenString)

	require.Error(t, err)
	assert.Equal(t, int64(0), userId)
}

func TestVerifyToken_ExpiredToken(t *testing.T) {
	testutil.SetupTestEnv(t)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  "test@example.com",
		"userId": float64(123),
		"exp":    time.Now().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
	})
	tokenString, _ := token.SignedString([]byte("test-secret-key-for-testing-only-do-not-use-in-production"))

	userId, err := VerifyToken(&tokenString)

	require.Error(t, err)
	assert.Equal(t, int64(0), userId)
}

func TestVerifyToken_EmptyToken(t *testing.T) {
	testutil.SetupTestEnv(t)

	emptyToken := ""

	userId, err := VerifyToken(&emptyToken)

	require.Error(t, err)
	assert.Equal(t, int64(0), userId)
}
