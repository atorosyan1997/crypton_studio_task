// Package model contains domain data structures.
package model

// SafeMapStats holds access and insert counters.
type SafeMapStats struct {
	AccessCount int64
	InsertCount int64
}
