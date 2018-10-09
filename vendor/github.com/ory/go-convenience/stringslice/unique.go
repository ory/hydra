package stringslice

func Unique(i []string) []string {
	u := make([]string, 0, len(i))
	m := make(map[string]bool)

	for _, val := range i {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}
