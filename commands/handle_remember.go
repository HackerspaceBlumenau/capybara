package commands

import (
    "fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/hackerspaceblumenau/capybara/models"
	"github.com/hackerspaceblumenau/capybara/slack"
)

func (s server) getOpenGraphFromURL(url string) (map[string]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var content []byte
	for {
		buf := make([]byte, 256)
		resp.Body.Read(buf)
		content = append(content, buf...)
		if strings.Index(string(content), "</head>") >= 0 {
			break
		}
	}

	reg := regexp.MustCompile(`(?m)<meta [^>]*property=[\"']og:(\w+)[\"'] [^>]*content=[\"']([^'^\"]+?)[\"'][^>]*>`)
	og := map[string]string{}
	for _, match := range reg.FindAllSubmatch(content, -1) {
		og[string(match[1])] = string(match[2])
	}

	return og, nil
}

func (s server) Remember(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	byteContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	content := string(byteContent)
	values, err := url.ParseQuery(content)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	split := strings.Split(values["text"][0], " ")
	when, err := time.Parse("02/01/2006 15:04", strings.Join(split[1:], " "))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

    if when.Before(time.Now()) {
        log.Printf("%s %s", when.String(), time.Now().String())
        slack.SendMessage("Ainda não tenho uma máquina do tempo...", values["channel_id"][0])
        return
    }

	uri := split[0]
	uri = strings.Trim(uri, " ")
	og, err := s.getOpenGraphFromURL(uri)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reminder := models.Reminder{
		Title:       og["title"],
		Description: og["description"],
		URL:         og["url"],
		Channel:     values["channel_id"][0],
		When:        when,
	}

	err = s.storage.SaveReminder(reminder)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

    msg := fmt.Sprintf("Opa! Eu vou lembrar todos neste canal no futuro sobre '%s'", reminder.Title)
	slack.SendMessage(msg, reminder.Channel)
}
