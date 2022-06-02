//go:build !js
// +build !js

package parser

import (
	"log"

	"gitlab.com/avoronkov/waver/lib/watch"
)

func (p *Parser) Start(wtch bool) error {
	if err := p.parse(); err != nil {
		return err
	}
	if wtch {
		err := watch.OnFileUpdate(p.file, func() {
			if err := p.parse(); err != nil {
				log.Printf("Parsing %v failed: %v", p.file, err)
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}
