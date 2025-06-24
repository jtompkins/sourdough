package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONArray is a generic type for handling JSON arrays of any type
type JSONArray[T any] []T

// Value implements the driver.Valuer interface
// Uses a pointer receiver for consistency with Scan()
func (ja *JSONArray[T]) Value() (driver.Value, error) {
	if ja == nil || *ja == nil {
		return nil, nil
	}
	return json.Marshal(*ja)
}

// Scan implements the sql.Scanner interface
// Uses a pointer receiver since we need to modify the receiver
func (ja *JSONArray[T]) Scan(value any) error {
	if value == nil {
		*ja = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("cannot scan non-string value into JSONArray")
	}

	return json.Unmarshal(bytes, ja)
}
