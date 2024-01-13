package utils

func In[T string | int](i T, list []T) bool {
	for _, el := range list {
		if el == i {
			return true
		}
	}
	return false
}

func InOps[T string | int](i T, list []T, cp func(el, i T) bool) bool {
	for _, el := range list {
		if cp(el, i) {
			return true
		}
	}
	return false
}
