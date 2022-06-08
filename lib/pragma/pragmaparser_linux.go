//go:build !js
// +build !js

package pragma

import (
	"log"

	"github.com/avoronkov/waver/lib/watch"
)

func (p *PragmaParser) Start(wtch bool) error {
	if err := p.Parse(); err != nil {
		return err
	}
	if !wtch {
		return nil
	}
	return watch.OnFileUpdate(p.file, func() {
		if err := p.Parse(); err != nil {
			log.Printf("Parsing pragmas of %v failed: %v", p.file, err)
		}
	})
}
