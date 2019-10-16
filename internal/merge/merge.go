package merge

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/tbouvin/p123-screen-get/config"
)

func findFiles(c config.Config, date string) ([]string, error) {
	matches, err := filepath.Glob(fmt.Sprintf("%s/*%s*", c.DownloadPath, date))
	if err != nil {
		return nil, err
	}

	return matches, nil
}

func MergeFiles(c config.Config, date string) error {
	files, err := findFiles(c, date)
	if err != nil {
		return err
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)

		//var nameLine string
		nameLine := scanner.Text()
		nameLine = nameLine[:len(nameLine)-1]

		_, err = f.Seek(0, 0)
		if err != nil {
			return err
		}

		csvFile := c.CSVPath + date + "_" + nameLine + ".csv"
		cf, err := os.OpenFile(csvFile, os.O_RDWR|os.O_CREATE, 0)
		if err != nil {
			return err
		}

		for scanner.Scan() {
			line := scanner.Text()
			//TODO This line is not right and needs work
			r, err := regexp.Compile(`r'\"([a-zA-Z ]+),([a-zA-Z ]+)\"', r'"\1\2"`)
			if err != nil {
				return err
			}

			substr := r.FindString(line)
			_, err = cf.WriteString(substr)
			if err != nil {
				return err
			}
		}

		f.Close()
		cf.Close()
	}

	return nil
}
