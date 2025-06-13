package utils

func ArrContains[T comparable](slice []T, element T) bool {
	m := make(map[T]struct{}, len(slice))
	for _, v := range slice {
		m[v] = struct{}{}
	}
	_, exists := m[element]
	return exists
}
