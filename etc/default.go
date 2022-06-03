package etc

import (
	_ "embed"
)

//go:embed config.yml
var DefaultConfig []byte

//go:embed example.pelia
var DefaultCodeExample []byte
