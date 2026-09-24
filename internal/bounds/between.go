package bounds

import "cmp"

// Between reports whether value is in the half-open interval [min, max).
func Between[T cmp.Ordered](value, min, max T) bool {
	return min <= value && value < max
}
