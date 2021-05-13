package merge

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/tbouvin/p123-screen-get/config"
)

func findFiles(c config.Config, date string) ([]string, error) {
	matches, err := filepath.Glob(fmt.Sprintf("%s/*%s.csv", c.FilePaths.DownloadPath, date))
	if err != nil {
		return nil, err
	}

	return matches, nil
}

func ConvertToCSV(c config.Config, date string) error {
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

		scanner.Scan()
		nameLine := scanner.Text()

		_, err = f.Seek(0, 0)
		if err != nil {
			return err
		}

		scanner = bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)

		csvFile := fmt.Sprintf("%s/%s_%s.csv", c.FilePaths.CSVPath, date, nameLine)
		cf, err := os.OpenFile(csvFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}

		for scanner.Scan() {
			line := scanner.Text()
			var re = regexp.MustCompile(`"([a-zA-Z ]+),([a-zA-Z ]+)"`)
			substr := re.ReplaceAllString(line, `"$1$2"`)

			_, err = cf.WriteString(substr + "\n")
			if err != nil {
				return err
			}
		}

		err = f.Close()
		if err != nil {
			return err
		}

		err = cf.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func MergeFiles(c config.Config, date string, todaysScreens []config.ScreenPart) error {
	var fileArray map[int][]*bufio.Scanner
	fileArray = make(map[int][]*bufio.Scanner)
	for index, i := range todaysScreens {
		for _, j := range i.Names {
			f, err := os.Open(fmt.Sprintf("%s/%s_%s.csv", c.FilePaths.CSVPath, date, j))
			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(f)
			scanner.Split(bufio.ScanLines)
			fileArray[index] = append(fileArray[index], scanner)
		}
	}

	combinedFile, err := os.OpenFile(fmt.Sprintf("%s/%s_combined.csv", c.FilePaths.CombinedPath, date), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	for index, _ := range c.Screens.Monday {
		newData := true
		lineNum := 0
		curLine := 0

		caprange := "PLACEHOLDER"
		dateinfile := "PLACEHOLDER"
		for newData {
			writeNewLine := true
			curLine = curLine + 1
			for j, screen := range fileArray[index] {

				//var nameLine string
				screen.Scan()
				line := screen.Text()

				//Get the caprange
				if lineNum == 0 {
					delimitLine := strings.Split(line, ",")
					if len(delimitLine) != 1 {
						return err
					} else {
						//delimitLineUnderscoreStr := delimitLine[0][:len(delimitLine[0])-1]
						delimitLineUnderscore := strings.Split(delimitLine[0], "_")
						caprange = delimitLineUnderscore[len(delimitLineUnderscore)-1]
					}
				}

				//Get the download date
				var delimitLine []string
				if lineNum == 1 {
					delimitLine = strings.Split(line, ",")

					if len(delimitLine) != 1 {
						return err
					} else {
						if delimitLine != nil && delimitLine[0] != "" {
							dateinfile = delimitLine[0]
						}
					}
				}

				if lineNum < 3 {
					writeNewLine = false
					continue
				}

				if j > 0 && lineNum < 3 {
					writeNewLine = false
					continue
				}

				if line == "" {
					newData = false
					continue
				}

				if index > 0 && j == 1 && lineNum == 3 {
					writeNewLine = false
					continue
				}

				//If not the first file
				if j > 0 && lineNum > 2 {
					if strings.Contains(line, ",") == false {
						continue
					}

					delimitLine = strings.Split(line, ",")
					for k, line := range delimitLine {
						if k < 5 {
							continue
						}

						if k != 5 {
							_, err = combinedFile.Write([]byte(","))
							if err != nil {
								return err
							}
						}

						//Remove the newline char from the last cell in the row
						if k == len(delimitLine)-1 {
							delimitLine = delimitLine[:len(delimitLine)-1]
						} else {
							_, err = combinedFile.Write([]byte(line))
							if err != nil {
								return err
							}
						}
					}
				} else {
					if index == 1 && lineNum == 3 {
						continue
					}

					if index == 0 && j == 0 && lineNum == 3 {
						//Write line to merged python file
						_, err = combinedFile.Write([]byte("Download Date,Cap Range," + line))
						if err != nil {
							return err
						}
					} else if lineNum != 3 && newData {
						_, err = combinedFile.Write([]byte(dateinfile + "," + caprange + "," + line))
						if err != nil {
							return err
						}
					}
				}
			}

			if writeNewLine && newData {
				_, err = combinedFile.Write([]byte("\n"))
				if err != nil {
					return err
				}
			}

			lineNum++
		}

		if newData {
			_, err = combinedFile.Write([]byte("\n"))
			if err != nil {
				return err
			}
		}
	}

	time.Sleep(500 * time.Millisecond)

	err = combinedFile.Close()
	if err != nil {
		return err
	}

	return nil
}
