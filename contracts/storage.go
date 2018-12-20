package contracts

import (
    "github.com/hackerspaceblumenau/capybara/models"
)

type Storage interface {
	SaveReminder(models.Reminder) error
    RemoveReminder(models.Reminder) error
    GetReminders() ([]models.Reminder, error)
}
