package stringslice

// IndexOf return index of matching string in slice
//
// Returns index or -1 if string not found
func IndexOf(matcher string, slice []string) int {
	for index, str := range slice {
		if matcher == str {
			return index
		}
	}
	return -1
}
