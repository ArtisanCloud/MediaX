package utils

// Contains 函数通过泛型检查切片中是否包含目标元素
func Contains[T comparable](arr []T, target T) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}
