package util

import "log/slog"

func Err(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func KV(key string, value any) slog.Attr {
	return slog.Any(key, value)
}
