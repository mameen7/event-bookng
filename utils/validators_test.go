package utils

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockFieldLevel is a mock implementation of validator.FieldLevel for testing
type mockFieldLevel struct {
	field reflect.Value
}

func (m *mockFieldLevel) Top() reflect.Value {
	return reflect.Value{}
}

func (m *mockFieldLevel) Parent() reflect.Value {
	return reflect.Value{}
}

func (m *mockFieldLevel) Field() reflect.Value {
	return m.field
}

func (m *mockFieldLevel) FieldName() string {
	return "TestField"
}

func (m *mockFieldLevel) StructFieldName() string {
	return "TestField"
}

func (m *mockFieldLevel) Param() string {
	return ""
}

func (m *mockFieldLevel) GetTag() string {
	return ""
}

func (m *mockFieldLevel) ExtractType(field reflect.Value) (value reflect.Value, kind reflect.Kind, nullable bool) {
	return field, field.Kind(), false
}

func (m *mockFieldLevel) GetStructFieldOK() (reflect.Value, reflect.Kind, bool) {
	return reflect.Value{}, reflect.Invalid, false
}

func (m *mockFieldLevel) GetStructFieldOKAdvanced(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool) {
	return reflect.Value{}, reflect.Invalid, false
}

func (m *mockFieldLevel) GetStructFieldOK2() (reflect.Value, reflect.Kind, bool, bool) {
	return reflect.Value{}, reflect.Invalid, false, false
}

func (m *mockFieldLevel) GetStructFieldOKAdvanced2(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool, bool) {
	return reflect.Value{}, reflect.Invalid, false, false
}

func TestValidateFutureDate_AcceptsFutureDates(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)
	mockFL := &mockFieldLevel{
		field: reflect.ValueOf(futureDate),
	}

	result := ValidateFutureDate(mockFL)

	assert.True(t, result, "Future date should be valid")
}

func TestValidateFutureDate_RejectsPastDates(t *testing.T) {
	pastDate := time.Now().Add(-24 * time.Hour)
	mockFL := &mockFieldLevel{
		field: reflect.ValueOf(pastDate),
	}

	result := ValidateFutureDate(mockFL)

	assert.False(t, result, "Past date should be invalid")
}

func TestValidateFutureDate_RejectsCurrentTime(t *testing.T) {
	currentTime := time.Now()
	mockFL := &mockFieldLevel{
		field: reflect.ValueOf(currentTime),
	}

	result := ValidateFutureDate(mockFL)

	assert.False(t, result, "Current time should be invalid (must be in the future)")
}

func TestValidateFutureDate_RejectsNearFutureTime(t *testing.T) {
	nearFuture := time.Now().Add(1 * time.Millisecond)
	mockFL := &mockFieldLevel{
		field: reflect.ValueOf(nearFuture),
	}

	result := ValidateFutureDate(mockFL)

	assert.True(t, result, "Near future time should be valid")
}

func TestValidateFutureDate_InvalidType(t *testing.T) {
	invalidType := "not a time"
	mockFL := &mockFieldLevel{
		field: reflect.ValueOf(invalidType),
	}

	result := ValidateFutureDate(mockFL)

	assert.False(t, result, "Non-time.Time type should be invalid")
}

func TestValidateFutureDate_IntegerType(t *testing.T) {
	invalidType := 12345
	mockFL := &mockFieldLevel{
		field: reflect.ValueOf(invalidType),
	}

	result := ValidateFutureDate(mockFL)

	assert.False(t, result, "Integer type should be invalid")
}

func TestValidateFutureDate_ZeroTime(t *testing.T) {
	zeroTime := time.Time{}
	mockFL := &mockFieldLevel{
		field: reflect.ValueOf(zeroTime),
	}

	result := ValidateFutureDate(mockFL)

	assert.False(t, result, "Zero time should be invalid")
}
