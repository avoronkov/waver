//go:build !js

package parser

import (
	"log"

	"github.com/avoronkov/waver/lib/watch"
)

func (p *Parser) Start(wtch bool) error {
	if err := p.parse(); err != nil {
		return err
	}
	if !wtch {
		return nil
	}
	return watch.OnFileUpdate(p.file, func() {
		if err := p.parse(); err != nil {
			log.Printf("Parsing %v failed: %v", p.file, err)
		}
	})
}
