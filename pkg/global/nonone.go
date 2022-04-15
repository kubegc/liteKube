package global

func NoValueString(s ...string) bool {
	if len(s) < 1 {
		return true
	}

	for _, str := range s {
		if str != "" {
			return false
		}
	}

	return true
}
