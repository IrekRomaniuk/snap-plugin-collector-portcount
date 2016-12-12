package targets

import (
	"os"
	"bufio"
	"fmt"
)

func ReadTargets(target string) ([]string, error) {
	/*if pathExists(target) {
		lines, err := readHosts(target)
		hosts = DeleteEmpty(lines)
		if err != nil {
			if err != nil { return nil, err }
		}
	}*/
	if !pathExists(target) {
		return nil, fmt.Errorf("File %s does not exist", target)
	}
	return readHosts(target)
}


func DeleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func readHosts(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//lines = append(lines, scanner.Text())
		line := scanner.Text()
		if line != "" {
			// omit line if empty
			lines = append(lines, scanner.Text())
			}
	}
	return lines, scanner.Err()
}

func pathExists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil { return true }
	if os.IsNotExist(err) { return false }
	return true
}
