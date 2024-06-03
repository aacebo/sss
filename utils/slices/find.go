package slices

func Find[S any](arr []S, cb func(S) bool) (S, bool) {
	var res S

	for _, v := range arr {
		if cb(v) {
			return v, true
		}
	}

	return res, false
}
