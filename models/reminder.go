package models

import (
    "time"
)

type Reminder struct {
	Title string
	Description string
	URL   string
    Channel string
	When  time.Time
}
