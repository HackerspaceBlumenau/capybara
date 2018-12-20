package ticker

import (
	"time"

	"github.com/hackerspaceblumenau/capybara/contracts"
)

var (
	tasks []func(time.Time)
)

type ticker struct {
	contracts.Server
	storage contracts.Storage
}

func (t ticker) registerTask(task func(time.Time)) {
	tasks = append(tasks, task)
}

func (t ticker) registerTasks() {
	t.registerTask(t.updateReminders)
	t.registerTask(t.checkWeeklyReminders)
	t.registerTask(t.checkDailyReminders)
	t.registerTask(t.checkNearReminders)
}

func (t ticker) Run() error {
	t.registerTasks()

	tickChan := time.Tick(time.Minute)
	for now := range tickChan {
		for _, task := range tasks {
			task(now)
		}
	}

	return nil
}

func NewServer(st contracts.Storage) contracts.Server {
	return ticker{storage: st}
}
