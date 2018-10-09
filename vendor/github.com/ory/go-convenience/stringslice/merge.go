package stringslice

func Merge(parts ...[]string) []string {
	var result []string
	for _, part := range parts {
		result = append(result, part...)
	}

	return result
}
