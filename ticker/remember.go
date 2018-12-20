package ticker

import (
	"log"
	"time"

	"github.com/hackerspaceblumenau/capybara/models"
	"github.com/hackerspaceblumenau/capybara/slack"
)

var (
	reminders           []models.Reminder
	remindersLastUpdate *time.Time
)

func filterReminders(fn func(models.Reminder) bool) []models.Reminder {
	filtered := []models.Reminder{}
	for _, reminder := range reminders {
		if fn(reminder) {
			filtered = append(filtered, reminder)
		}
	}

	return filtered
}

func (t ticker) checkWeeklyReminders(date time.Time) {
	if date.Weekday() != time.Monday {
		return
	}

	log.Println("Checking weekly reminders")
}

func (t ticker) checkDailyReminders(date time.Time) {
	if date.Hour() < 7 && date.Hour() > 9 {
		return
	}

	log.Println("Checking daily reminders")
}

func (t ticker) checkNearReminders(date time.Time) {
	inFifteen := filterReminders(func(r models.Reminder) bool {
		return r.When.After(date) && r.When.Before(date.Add(15*time.Minute))
	})

	if len(inFifteen) == 0 {
		return
	}

	next := inFifteen[0]
	slack.SendMessage(next.Title, next.Channel)
}

func (t ticker) updateReminders(date time.Time) {
	if remindersLastUpdate != nil {
		if remindersLastUpdate.Before(date.Add(time.Hour)) {
			return
		}
	}

	log.Println("Updating reminders list")

	current, err := t.storage.GetReminders()
	if err != nil {
		return
	}

	reminders = current
	now := time.Now()
	remindersLastUpdate = &now
}
