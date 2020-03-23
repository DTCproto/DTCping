package file

import (
	"bufio"
	"os"
)

func ReadIpsFile(filePath string) ([]string, error) {
	var ips []string
	fd, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return ips, err
	}
	defer fd.Close()

	sc := bufio.NewScanner(fd)
	for sc.Scan() {
		ips = append(ips, sc.Text())
	}
	return ips, nil
}
