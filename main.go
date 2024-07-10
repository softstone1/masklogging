package main

import (
	"context"
	"log"
	"os"
	"reflect"
	"runtime/pprof"
	"time"

	"log/slog"
)

// SensitiveData is a custom type for sensitive fields
type SensitiveData string

// Implement the LogValuer interface for SensitiveData
func (s SensitiveData) LogValue() slog.Value {
	return slog.StringValue("REDACTED")
}

// Nested structs with sensitive fields
type Country struct {
	Name      string        `json:"name"`
	Code      SensitiveData `json:"code"`
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

// ConvertStructToAttrs converts a struct to a list of slog.Attr, handling nested structs
func ConvertStructToAttrs(v interface{}) []slog.Attr {
	return convertStructToAttrs(reflect.ValueOf(v), "")
}

func convertStructToAttrs(val reflect.Value, prefix string) []slog.Attr {
	attrs := []slog.Attr{}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		if prefix != "" {
			fieldName = prefix + "." + fieldName
		}

		if field.Kind() == reflect.Struct && !field.Type().Implements(reflect.TypeOf((*slog.LogValuer)(nil)).Elem()) {
			// Recursively process nested struct
			attrs = append(attrs, convertStructToAttrs(field, fieldName)...)
		} else {
			attrs = append(attrs, slog.Any(fieldName, field.Interface()))
		}
	}

	return attrs
}

func main() {
    // Create a new logger with JSON handler
    handlerOptions := &slog.HandlerOptions{}
    logger := slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))

    // Create a user instance
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

    // Start CPU profiling
    f, err := os.Create("cpu.prof")
    if err != nil {
        log.Fatal("could not create CPU profile: ", err)
    }
    if err := pprof.StartCPUProfile(f); err != nil {
        log.Fatal("could not start CPU profile: ", err)
    }
    defer pprof.StopCPUProfile()

    // Run the conversion function multiple times
    for i := 0; i < 1000000; i++ {
        _ = ConvertStructToAttrs(user)
    }

    // Convert the user struct to slog.Attrs and log it
    attrs := ConvertStructToAttrs(user)
    logger.LogAttrs(context.Background(), slog.LevelInfo, "User data", attrs...)

    // Allow some time for profiling data to be written
    time.Sleep(2 * time.Second)
}


