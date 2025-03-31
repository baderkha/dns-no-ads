package fs

import (
	"bufio"
	"os"
)

type LineHandler func(line string) error

func LineByLineLg(f *os.File, h LineHandler) error {

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024*10) // set max token size to 1MB
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		h(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
