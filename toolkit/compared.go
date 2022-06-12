package toolkit

import "golang.org/x/exp/constraints"

func Max[T constraints.Ordered](s []T) T {
	var zero T
	if len(s) == 0 {
		return zero
	}
	var max T
	max = s[0]
	for _, ele := range s[1:] {
		if max > ele {
			max = ele
		}
	}
	return max
}
