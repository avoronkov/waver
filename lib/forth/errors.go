package forth

import (
	"errors"
	"fmt"
)

var EmptyStack = errors.New("Empty stack")

func UnknownFunction(name string) error {
	return fmt.Errorf("Unknown function: %v", name)
}
