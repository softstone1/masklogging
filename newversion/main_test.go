package main

import (
	"log/slog"
	"testing"
)

// BenchmarkToLogAttrs benchmarks the ToLogAttrs method
func BenchmarkToLogAttrs(b *testing.B) {
	user := User{
		Name:        "John Doe",
		Email:       SensitiveData("john.doe@example.com"),
		PhoneNumber: SensitiveData("+1234567890"),
		Address: Address{
			Street:   "123 Elm Street",
			City:     "Springfield",
			Postcode: SensitiveData("12345"),
			Country: Country{
				Name: "USA",
				Code: SensitiveData("US"),
			},
		},
	}
	// Pre-allocate slice for log attributes
	attrs := make([]slog.Attr, 0, 10)

	for i := 0; i < b.N; i++ {
		_ = user.ToLogAttrs(attrs)
	}
}
