package share

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const defaultBaseUrl = "https://avoronkov.github.io/waver"

func MakeLink(data string) (string, error) {
	code, err := Encode(data)
	if err != nil {
		return "", err
	}
	link := fmt.Sprintf("%v?code=%v", defaultBaseUrl, code)
	return link, nil
}

func MakeLinkFromFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	str := strings.TrimSpace(string(data))
	return MakeLink(str)
}
