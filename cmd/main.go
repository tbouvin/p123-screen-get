package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/tbouvin/p123-screen-get/config"
	"github.com/tbouvin/p123-screen-get/internal/merge"
	p123 "github.com/tbouvin/p123-screen-get/internal/selenium"
)

const alreadyRanFileName = "/screensRan"

func checkIfDownloadsExist(d string, path string) bool {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}

	date := string(b)
	if date == d {
		return true
	}

	return false
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

	exists := checkIfDownloadsExist(formattedDate, c.FilePaths.DownloadPath+alreadyRanFileName)
	if exists {
		return
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
				fileName := fmt.Sprintf("%s/%d_%s.csv", c.FilePaths.DownloadPath, count, formattedDate)
				for retry := 0; retry < 5; retry++ {
					err = d.GetScreen(screen, fileName)
					if err == nil {
						break
					}
				}

				if err != nil {
					panic(err)
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

	err = ioutil.WriteFile(c.FilePaths.DownloadPath+alreadyRanFileName, []byte(formattedDate), 0644)
	if err != nil {
		panic(err)
	}
}
