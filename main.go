package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	token := os.Getenv("SWAT_BOT_SLACK_TOKEN")
	if token == "" {
		log.Fatal("Need a non-empty slack token - check the SWAT_BOT_SLACK_TOKEN environment variable.")
	}

	ticker := time.NewTicker(time.Minute * 20)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(24 * time.Hour)
		done <- true
	}()
	postMessageURL := "https://slack.com/api/chat.postMessage"
	headers := []string{
		"This just in",
		"Just in from App Support",
		"Fresh from the release",
		"[Escalated]",
		"Heard from Client Care",
	}
	swatAlerts := []string{
		"Reports are up sitewide - an unexpected success has occurred.",
		"None of lfeldmandemo's 17,000 users can log in!",
		"Vampire users spotted in the LMS database.",
		"Wifi blip has been escalated to Jim T.",
		"SAML is not working in China!!",
		"Germany has just reported that Assessments are not working!!",
		"Referrer URL E2E tests are failing!!",
		"SUPERHEROADMIN user BROKEN!!!",
		"Academy users are maybe there, but maybe they aren't!",
		"Users are unable to fit their certificates into frames when printed!",
		"Data export utility has been running non-stop for 1 week!",
		"Enable HTTPS organization setting has been accidentally enabled for every site!",
		"Hierarchy tree with depth = 1,000,000 causing stack overflow!",
		"All web services are returning a -900. What could it mean?",
		"Request blip as azure web app was swapped - learner logged in one second slower.",
		"Rapid user save completely deletes the user and all associated CE credits.",
		"Authentication ticket did not account for Untraditional SSO - \"historic\" customers broken.",
		"User field was blanked out when blank value was provided!",
		"UpdateUser13 is returning a -600.",
		"Competency Tracker Reporter sees the \"Internal Reports\" report.",
		"LMS Account creation is happening too fast, causing token redemption successes on Academy.",
		"Site is down due to F5 key on the My Learning page.",
		"Site2 users can't proxy login into other sites!",
		"The newsfeed section has magically reappeared!",
	}

	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case t := <-ticker.C:
			if len(swatAlerts) == 0 {
				return
			}
			source := rand.NewSource(time.Now().UnixNano())
			random := rand.New(source)

			headerIndex := random.Intn(len(headers))
			index := random.Intn(len(swatAlerts))

			header := headers[headerIndex]
			swatAlert := swatAlerts[index]
			swatAlerts = remove(swatAlerts, index)
			resp, err := http.PostForm(postMessageURL, url.Values{"token": {token}, "channel": {"fake-swat"}, "text": {fmt.Sprintf("<!channel> %s - %s", header, swatAlert)}})
			defer resp.Body.Close()
			if err != nil {
				fmt.Printf("An error occurred sending the message - %s", err.Error())
			}
			if resp.StatusCode != http.StatusOK {
				bodyBytes, _ := ioutil.ReadAll(resp.Body)
				body := string(bodyBytes)
				fmt.Printf("Bad response - not a 200 - %s", body)
			}
			fmt.Println("Current time: ", t)
		}
	}
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}
