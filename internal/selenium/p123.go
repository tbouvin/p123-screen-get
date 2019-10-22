package p123

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/tebeka/selenium"

	"github.com/tbouvin/p123-screen-get/config"
)

type Driver struct {
	wd   selenium.WebDriver
	conf config.Config
	stop func() error
	cmd  *exec.Cmd
}

//conf.Selenium.Command
func Init(conf config.Config) Driver {
	args := append(conf.Selenium.Arguments, "-port")
	args = append(args, conf.Selenium.Port)
	cmd := exec.Command(conf.Selenium.Command, args...)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%s/wd/hub", conf.Selenium.Port))
	if err != nil {
		panic(err)
	}

	return Driver{wd: wd, conf: conf, stop: nil, cmd: cmd}
}

func (d Driver) Logout() error {
	err := d.wd.Close()
	if err != nil {
		return err
	}

	err = d.wd.Quit()
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	err = d.cmd.Process.Kill()
	if err != nil {
		return err
	}

	return nil
}

func (d Driver) Login() error {
	err := d.wd.Get(d.conf.URLs.Login)
	if err != nil {
		return err
	}

	err = d.enterText(d.conf.IDs.Username, d.conf.Credentials.Username)
	if err != nil {
		return err
	}

	err = d.enterText(d.conf.IDs.Password, d.conf.Credentials.Password)
	if err != nil {
		return err
	}

	err = d.clickXpath(d.conf.Xpaths.LoginButton)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	url, err := d.wd.CurrentURL()
	if err != nil {
		return err
	}

	if strings.Contains(url, "auth2fact") {
		err = d.enterText(d.conf.IDs.SecondaryPassword, d.conf.Credentials.SecondaryPassword)
		if err != nil {
			return err
		}

		err = d.clickXpath(d.conf.Xpaths.SecondaryLoginButton)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d Driver) GetScreen(screenName string, fileName string) error {
	err := d.wd.Get(d.conf.URLs.Screen)
	if err != nil {
		return err
	}

	err = d.clickXpath(d.conf.Xpaths.ShowAllScreenButton)
	if err != nil {
		return err
	}

	err = d.clickLink(screenName)
	if err != nil {
		return err
	}

	err = d.clickID(d.conf.IDs.RunScreenButton)
	if err != nil {
		return err
	}

	err = d.clickLink(d.conf.IDs.TickerLink)
	if err != nil {
		return err
	}

	elem, err := d.wd.FindElement(selenium.ByXPATH, d.conf.Xpaths.ScreenDownload)
	if err != nil {
		return err
	}

	link, err := elem.GetAttribute("href")
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return err
	}

	cookies, err := d.wd.GetCookies()
	for _, cookie := range cookies {
		req.AddCookie(&http.Cookie{Name: cookie.Name, Value: cookie.Value})
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//write response content to file
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (d Driver) clickXpath(xpath string) error {
	btn, err := d.wd.FindElement(selenium.ByXPATH, xpath)
	if err != nil {
		return err
	}

	err = btn.Click()
	if err != nil {
		return err
	}

	return nil
}

func (d Driver) clickID(id string) error {
	btn, err := d.wd.FindElement(selenium.ByID, id)
	if err != nil {
		return err
	}

	err = btn.Click()
	if err != nil {
		return err
	}

	return nil
}

func (d Driver) clickLink(link string) error {
	elem, err := d.wd.FindElement(selenium.ByLinkText, link)
	if err != nil {
		return err
	}

	err = elem.Click()
	if err != nil {
		return err
	}

	return nil
}

func (d Driver) enterText(id string, text string) error {
	elem, err := d.wd.FindElement(selenium.ByID, id)
	if err != nil {
		return err
	}

	err = elem.SendKeys(text)
	if err != nil {
		return err
	}

	return nil
}
