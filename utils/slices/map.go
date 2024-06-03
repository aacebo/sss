package slices

func Map[S, R any](arr []S, cb func(S) R) []R {
	res := make([]R, len(arr))

	for i := range arr {
		res[i] = cb(arr[i])
	}

	return res
}
