package stringslice

// Exists return bool if string exists in string slice
func Exists(matcher string, slice []string) bool {
	return IndexOf(matcher, slice) != -1
}
