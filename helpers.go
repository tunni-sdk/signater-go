package signater

// String returns a pointer to s, for optional params fields.
func String(s string) *string { return &s }

// Bool returns a pointer to b, for optional params fields.
func Bool(b bool) *bool { return &b }

// Int returns a pointer to i, for optional params fields.
func Int(i int) *int { return &i }

// Float64 returns a pointer to f, for optional params fields.
func Float64(f float64) *float64 { return &f }
