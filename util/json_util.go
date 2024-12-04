package util

import (
	"encoding/json"
	"fmt"
)

// UnmarshalSafe TypeSafe JSON Unmarshal
func UnmarshalSafe[T any](data []byte) (T, error) {
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return result, nil
}

// UnmarshalUnsafe Unsafe JSON unmarshal to map
func UnmarshalUnsafe(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return result, nil
}

// ToJSON struct to json
func ToJSON(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal: %w", err)
	}
	return string(data), nil
}

// GetField get fields details from json
func GetField(data []byte, field string) (interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	if val, ok := m[field]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("field %s not found", field)
}

func MustUnmarshalSafe[T any](data []byte) T {
	result, err := UnmarshalSafe[T](data)
	if err != nil {
		panic(err)
	}
	return result
}

// UnmarshalWithDefault unmarshal with default value
func UnmarshalWithDefault[T any](data []byte, defaultValue T) T {
	result, err := UnmarshalSafe[T](data)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetNestedField safely get embedded
func GetNestedField(data []byte, fields ...string) (interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	current := m
	for i, field := range fields[:len(fields)-1] {
		if next, ok := current[field].(map[string]interface{}); ok {
			current = next
		} else {
			return nil, fmt.Errorf("field %s at depth %d is not an object", field, i)
		}
	}

	lastField := fields[len(fields)-1]
	if val, ok := current[lastField]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("field %s not found", lastField)
}

// PrettyPrint Pretty print JSON
func PrettyPrint(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
