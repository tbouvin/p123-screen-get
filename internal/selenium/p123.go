package selenium

import (
	"fmt"
	"os"
	"strings"

	"github.com/tebeka/selenium"

	"github.com/tbouvin/p123-screen-get/config"
)

type Driver struct {
	wd   selenium.WebDriver
	conf config.Config
}

func Init() Driver {
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
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}

	return Driver{wd: wd}
}

func (d Driver) GoHome() {
	if err := d.wd.Get("www.portfolio123.com"); err != nil {
		panic(err)
	}
}

func (d Driver) Login() error {
	err := d.wd.Get("https://www.portfolio123.com/app/auth")
	if err != nil {
		return err
	}

	err = d.enterText("user", "username")
	if err != nil {
		return err
	}

	err = d.enterText("passwd", "password")
	if err != nil {
		return err
	}

	err = d.clickButton("//*[text()='Submit']")
	if err != nil {
		return err
	}

	url, err := d.wd.CurrentURL()
	if err != nil {
		return err
	}

	if strings.Contains(url, "auth2fact") {

		err = d.enterText("passwd", "secondpassword")
		if err != nil {
			return err
		}

		err = d.clickButton("//*[@id=\"wrapper\"]/div[1]/div/div/div/div[1]/div/form/div/div/div[2]/div[2]/button")
		if err != nil {
			return err
		}
	}

	return nil
}

func (d Driver) clickButton(xpath string) error {
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
