package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/tbouvin/p123-screen-get/internal/merge"

	"github.com/tbouvin/p123-screen-get/config"
	p123 "github.com/tbouvin/p123-screen-get/internal/selenium"
)

func main() {
	mergeFiles := flag.Bool("mergefiles", true, "False will not merge CSVs")
	getScreens := flag.Bool("getscreens", true, "False will not get screens from p123")
	convertFiles := flag.Bool("convertfiles", true, "False will not convert screen files to CSV")
	flag.Parse()
	*convertFiles = true
	*mergeFiles = true
	*getScreens = false

	c, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	formattedDay := time.Now().Day()
	formattedMonth := time.Now().Month()
	formattedYear := time.Now().Year() - 2000
	formattedDate := fmt.Sprintf("%02d%02d%02d", formattedMonth, formattedDay, formattedYear)

	var screenDay []config.ScreenPart

	switch time.Now().Weekday() {
	case time.Monday:
		screenDay = c.Screens.Monday
	case time.Tuesday:
		screenDay = c.Screens.Tuesday
	case time.Wednesday:
		screenDay = c.Screens.Wednesday
	case time.Thursday:
		screenDay = c.Screens.Thursday
	case time.Friday:
		screenDay = c.Screens.Friday
	case time.Saturday:
		screenDay = c.Screens.Saturday
	case time.Sunday:
		screenDay = c.Screens.Sunday
	default:
		panic(nil)
	}

	if *getScreens {
		d := p123.Init(c)
		err = d.Login()
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)

		count := 0
		for _, screenSet := range screenDay {
			for _, screen := range screenSet.Names {
				fileName := fmt.Sprintf("%s/%d%s.xls", c.FilePaths.DownloadPath, count, formattedDate)
				err = d.GetScreen(screen, fileName)
				if err != nil {
					panic(err)
				}
				count++
			}
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
}
