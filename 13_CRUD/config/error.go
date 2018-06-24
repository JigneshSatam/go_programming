package config

// ParseError is used to parse error if present
func ParseError(err error) {
	if err != nil {
		panic(err)
	}
}
