package util

import "log/slog"

// Err creates a slog.Attr for an error
// Returns an empty Attr if err is nil, otherwise creates an Attr with key "error"
// and the error message as value
func Err(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

// KV creates a slog.Attr with the given key and value
// This is a convenience wrapper around slog.Any
func KV(key string, value any) slog.Attr {
	return slog.Any(key, value)
}
