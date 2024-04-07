package util

// Option represents an optional value.
// If OK is true, then V contains the value that can be used.
// If OK is false, then value V should not be used.
type Option[T any] struct {
	V  T
	OK bool
}
