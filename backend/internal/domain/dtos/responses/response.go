package responses

import (
	"encoding/json"
	"reflect"
)

type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

type Response[T any] struct {
	Status  Status `json:"status" example:"error|success"`
	Data    T      `json:"data"`
	Message string `json:"message,omitempty"`
}

// Custom JSON marshal: success + nil'ler güzelce normalize edilir
func (r Response[T]) MarshalJSON() ([]byte, error) {
	type Alias Response[T]

	// success değilse direkt bas
	if r.Status != StatusSuccess {
		return json.Marshal((Alias)(r))
	}

	v := any(r.Data)
	if v == nil {
		// Tür tamamen belirsizse [] tercih edelim (isteğin doğrultusunda)
		return json.Marshal(struct {
			Status  Status `json:"status"`
			Data    []any  `json:"data"`
			Message string `json:"message,omitempty"`
		}{r.Status, []any{}, r.Message})
	}

	rv := reflect.ValueOf(v)
	// Pointer veya interface nil'iyse -> {} (tekil nesne beklentisi)
	if (rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface) && rv.IsNil() {
		return json.Marshal(struct {
			Status  Status         `json:"status"`
			Data    map[string]any `json:"data"`
			Message string         `json:"message,omitempty"`
		}{r.Status, map[string]any{}, r.Message})
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		// nil slice -> []
		if rv.IsNil() {
			return json.Marshal(struct {
				Status  Status `json:"status"`
				Data    []any  `json:"data"`
				Message string `json:"message,omitempty"`
			}{r.Status, []any{}, r.Message})
		}
	case reflect.Map:
		// nil map -> {}
		if rv.IsNil() {
			return json.Marshal(struct {
				Status  Status         `json:"status"`
				Data    map[string]any `json:"data"`
				Message string         `json:"message,omitempty"`
			}{r.Status, map[string]any{}, r.Message})
		}
	}

	// Normal durum: olduğu gibi bas
	return json.Marshal((Alias)(r))
}

func Success[T any](data T, message string) Response[T] {
	return Response[T]{Status: StatusSuccess, Data: data, Message: message}
}

func Error[T any](data T, message string) Response[T] {
	return Response[T]{Status: StatusError, Data: data, Message: message}
}
