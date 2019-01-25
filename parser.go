package cli

var _ MultiParser = Const(0)

// Const parser returns constant value if flag is present.
func Const(v interface{}) MultiParser {
	return func([]string) (interface{}, int, error) {
		return v, 0, nil
	}
}
