package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/tbouvin/p123-screen-get/config"
	"github.com/tbouvin/p123-screen-get/internal/merge"
	p123 "github.com/tbouvin/p123-screen-get/internal/selenium"
)

func didScreenRun(screenName string, path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	found := false
	for scanner.Scan() {
		t := scanner.Text()
		if screenName == t {
			found = true
		}
	}

	err = scanner.Err()
	if err != nil || !found {
		return false
	}

	return true
}

func main() {
	mergeFiles := flag.Bool("mergefiles", true, "False will not merge CSVs")
	getScreens := flag.Bool("getscreens", true, "False will not get screens from p123")
	convertFiles := flag.Bool("convertfiles", true, "False will not convert screen files to CSV")
	configFile := flag.String("config", "resources/local/config.yml", "Location of yaml configuration")
	flag.Parse()

	c, err := config.GetConfig(configFile)
	if err != nil {
		panic(err)
	}

	formattedDay := time.Now().Day()
	formattedMonth := time.Now().Month()
	formattedYear := time.Now().Year() - 2000
	formattedDate := fmt.Sprintf("%02d%02d%02d", formattedMonth, formattedDay, formattedYear)

	alreadyRanFileName := c.FilePaths.DownloadPath + "/." + formattedDate

	var screenDay []config.ScreenPart

	weekday := time.Now().Weekday()
	switch weekday {
	case time.Monday:
		screenDay = c.Screens.Monday
		break
	case time.Tuesday:
		screenDay = c.Screens.Tuesday
		break
	case time.Wednesday:
		screenDay = c.Screens.Wednesday
		break
	case time.Thursday:
		screenDay = c.Screens.Thursday
		break
	case time.Friday:
		screenDay = c.Screens.Friday
		break
	case time.Saturday:
		screenDay = c.Screens.Saturday
		break
	case time.Sunday:
		screenDay = c.Screens.Sunday
		break
	default:
		panic(nil)
	}

	file, err := os.OpenFile(alreadyRanFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	if *getScreens {
		d := p123.Init(c)
		err = d.Login()
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)

		count := 1
		for _, screenSet := range screenDay {
			for _, screen := range screenSet.Names {
				if !didScreenRun(screen, alreadyRanFileName) {
					fileName := fmt.Sprintf("%s/%d_%s.csv", c.FilePaths.DownloadPath, count, formattedDate)
					for retry := 0; retry < 5; retry++ {
						err = d.GetScreen(screen, fileName)
						if err == nil {
							_, _ = file.Write([]byte(screen))
							_, _ = file.Write([]byte("\n"))
							break
						}
					}

					if err != nil {
						panic(err)
					}
				}
				count++
			}
		}

		err = d.Logout()
		if err != nil {
			panic(err)
		}
	}

	if *convertFiles {
		err = merge.ConvertToCSV(c, formattedDate)
		if err != nil {
			panic(err)
		}
	}

	time.Sleep(500 * time.Millisecond)

	if *mergeFiles {
		err = merge.MergeFiles(c, formattedDate, screenDay)
		if err != nil {
			panic(err)
		}
	}

	time.Sleep(time.Duration(c.SleepTime) * time.Second)
}
