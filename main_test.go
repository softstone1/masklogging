package main

import "testing"

// BenchmarkConvertStructToAttrs benchmarks the ConvertStructToAttrs function
func BenchmarkConvertStructToAttrs(b *testing.B) {
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

	for i := 0; i < b.N; i++ {
		ConvertStructToAttrs(user)
	}
}
