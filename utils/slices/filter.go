package slices

func Filter[S any](arr []S, cb func(S) bool) []S {
	res := []S{}

	for _, v := range arr {
		if cb(v) {
			res = append(res, v)
		}
	}

	return res
}
