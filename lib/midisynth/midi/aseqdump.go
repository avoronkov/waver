package midi

import (
	"io"
	"os/exec"
	"strconv"
)

func aseqdump(p int) (*exec.Cmd, io.Reader, error) {
	cmd := exec.Command("aseqdump", "-p", strconv.Itoa(p))
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}
	return cmd, reader, nil
}
