package drivers

import (
    "log"

	"github.com/hackerspaceblumenau/capybara/models"
)

type MemoryStorage struct {
	reminders map[string]models.Reminder
}

func (m MemoryStorage) SaveReminder(r models.Reminder) error {
	if m.reminders == nil {
		m.reminders = map[string]models.Reminder{}
	}

	m.reminders[r.Title] = r
    log.Println("Saved new reminder in memory")
	return nil
}

func (m MemoryStorage) RemoveReminder(r models.Reminder) error {
    delete(m.reminders, r.Title)
    return nil
}

func (m MemoryStorage) GetReminders() ([]models.Reminder, error) {
	list := []models.Reminder{}
	for _, r := range m.reminders {
		list = append(list, r)
	}

	return list, nil
}
