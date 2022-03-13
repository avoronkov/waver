package project

import (
	"fmt"
	"io/ioutil"
	"os"
)

func New(name string, configData []byte) error {
	peliaName := fmt.Sprintf("%v.pelia", name)

	f, err := os.Create(peliaName)
	if err != nil {
		return err
	}
	f.Close()

	configName := fmt.Sprintf("%v.yml", name)
	if err := ioutil.WriteFile(configName, configData, 0644); err != nil {
		return err
	}

	return nil
}
