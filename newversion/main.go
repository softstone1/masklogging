package main

import (
	"context"
	"os"

	"log/slog"
)

type SensitiveData string

func (s SensitiveData) LogValue() slog.Value {
	return slog.StringValue("REDACTED")
}

type Country struct {
	Name string        `json:"name"`
	Code SensitiveData `json:"code"`
}

type Address struct {
	Street   string        `json:"street"`
	City     string        `json:"city"`
	Postcode SensitiveData `json:"postcode"`
	Country  Country       `json:"country"`
}

type User struct {
	Name        string        `json:"name"`
	Email       SensitiveData `json:"email"`
	PhoneNumber SensitiveData `json:"phone_number"`
	Address     Address       `json:"address"`
}

// Convert the User struct to slog.Attr
func (u User) ToLogAttrs(attrs []slog.Attr) []slog.Attr {
	attrs = append(attrs,
		slog.String("Name", u.Name),
		slog.Any("Email", u.Email),
		slog.Any("PhoneNumber", u.PhoneNumber),
		slog.Group("Address",
			slog.String("Street", u.Address.Street),
			slog.String("City", u.Address.City),
			slog.Any("Postcode", u.Address.Postcode),
			slog.Group("Country",
				slog.String("Name", u.Address.Country.Name),
				slog.Any("Code", u.Address.Country.Code),
			),
		),
	)
	return attrs
}

func main() {
	handlerOptions := &slog.HandlerOptions{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))

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
	attrs := make([]slog.Attr, 0, 20) // Increase initial capacity to reduce resizes

	// Log the user struct
	logger.LogAttrs(context.Background(), slog.LevelInfo, "User data", user.ToLogAttrs(attrs[:0])...)

}
