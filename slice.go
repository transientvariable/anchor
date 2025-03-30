package anchor

// Unique returns the slice containing the unique values from the provided slice.
func Unique[S comparable](values ...S) []S {
	if len(values) <= 1 {
		return values
	}

	k := make(map[S]bool)
	var d []S
	for _, v := range values {
		if _, value := k[v]; !value {
			k[v] = true
			d = append(d, v)
		}
	}
	return d
}
