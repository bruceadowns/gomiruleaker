package lib

const (
	maxLength = 80
)

// Trunc ...
func Trunc(s string) string {
	res := s
	if len(res) > maxLength {
		res = res[:maxLength] + "..."
	}

	return res
}
