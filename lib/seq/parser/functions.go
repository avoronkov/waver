package parser

import (
	"fmt"
	"strconv"

	"gitlab.com/avoronkov/waver/lib/seq/common"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

type ModParser = func(fields []string) (types.Modifier, int, error)

// : 4
func parseEvery(fields []string) (types.Modifier, int, error) {
	if len(fields) < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for Every (':')")
	}

	n, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil, 0, err
	}

	return common.Every(int64(n)), 2, nil
}
