package main

import (
	"fmt"
	"time"

	"github.com/tbouvin/p123-screen-get/config"
	p123 "github.com/tbouvin/p123-screen-get/internal/selenium"
)

func main() {
	c, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", c)
	d := p123.Init(c)
	err = d.Login()
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)

	formattedDay := time.Now().Weekday().String()
	formattedMonth := time.Now().Month().String()
	formattedYear := string(time.Now().Year())
	formattedDate := formattedMonth + formattedDay + formattedYear

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

	count := 0
	for _, screenSet := range screenDay {
		for _, screen := range screenSet.Names {
			fileName := fmt.Sprintf("%s/%d%s.xls", c.DownloadPath, count, formattedDate)
			err = d.GetScreen(screen, fileName)
			if err != nil {
				panic(err)
			}
			count++
		}
	}
}
