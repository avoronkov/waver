package etc

import (
	_ "embed"
)

//go:embed config.yml
var DefaultConfig []byte
