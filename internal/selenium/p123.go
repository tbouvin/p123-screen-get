package selenium

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/tebeka/selenium"

	"github.com/tbouvin/p123-screen-get/config"
)

type Driver struct {
	wd   selenium.WebDriver
	conf config.Config
	stop func() error
}

func Init(conf config.Config) Driver {
	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	const (
		// These paths will be different on your system.
		seleniumPath    = "vendor/selenium-server.jar"
		geckoDriverPath = "vendor/geckodriver"
		port            = 8080
	)
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),            // Output debug information to STDERR.
	}
	selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}

	return Driver{wd: wd, conf: conf, stop: service.Stop}
}

func (d Driver) GoHome() {
	if err := d.wd.Get("www.portfolio123.com"); err != nil {
		panic(err)
	}
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

func (d Driver) GetScreen(screenName string) error {
	err := d.wd.Get(d.conf.URLs.Login)
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

	//write response content to file

	err = resp.Body.Close()
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
