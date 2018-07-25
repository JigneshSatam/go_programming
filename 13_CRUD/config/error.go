package config

import (
	"log"
)

// ParseError is used to parse error if present
func ParseError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
