package scheduler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gproc/pkg/types"
)

type CronScheduler struct {
	tasks   map[string]*types.ScheduledTask
	running map[string]bool
	mutex   sync.RWMutex
	ticker  *time.Ticker
	stop    chan bool
}

func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		tasks:   make(map[string]*types.ScheduledTask),
		running: make(map[string]bool),
		stop:    make(chan bool),
	}
}

func (cs *CronScheduler) Start() {
	cs.ticker = time.NewTicker(1 * time.Minute)
	go cs.run()
}

func (cs *CronScheduler) Stop() {
	if cs.ticker != nil {
		cs.ticker.Stop()
	}
	cs.stop <- true
}

func (cs *CronScheduler) AddTask(task *types.ScheduledTask) error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	
	nextRun, err := cs.parseNextRun(task.Cron)
	if err != nil {
		return fmt.Errorf("invalid cron expression: %v", err)
	}
	
	task.NextRun = nextRun
	cs.tasks[task.Name] = task
	return nil
}

func (cs *CronScheduler) RemoveTask(name string) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	
	delete(cs.tasks, name)
	delete(cs.running, name)
}

func (cs *CronScheduler) ListTasks() []*types.ScheduledTask {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	
	var tasks []*types.ScheduledTask
	for _, task := range cs.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (cs *CronScheduler) run() {
	for {
		select {
		case <-cs.ticker.C:
			cs.checkAndRunTasks()
		case <-cs.stop:
			return
		}
	}
}

func (cs *CronScheduler) checkAndRunTasks() {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	
	now := time.Now()
	for name, task := range cs.tasks {
		if now.After(task.NextRun) && !cs.running[name] {
			go cs.executeTask(task)
			
			// Schedule next run
			nextRun, err := cs.parseNextRun(task.Cron)
			if err == nil {
				task.NextRun = nextRun
			}
		}
	}
}

func (cs *CronScheduler) executeTask(task *types.ScheduledTask) {
	cs.mutex.Lock()
	cs.running[task.Name] = true
	cs.mutex.Unlock()
	
	defer func() {
		cs.mutex.Lock()
		cs.running[task.Name] = false
		cs.mutex.Unlock()
	}()
	
	log.Printf("Executing scheduled task: %s", task.Name)
	// Execute the task command
	// This would integrate with the process manager
}

func (cs *CronScheduler) parseNextRun(cronExpr string) (time.Time, error) {
	// Simple cron parser - supports basic patterns
	// Format: "minute hour day month weekday"
	// Example: "0 2 * * *" = daily at 2 AM
	
	now := time.Now()
	
	// Basic patterns
	switch cronExpr {
	case "@hourly", "0 * * * *":
		return now.Add(time.Hour).Truncate(time.Hour), nil
	case "@daily", "0 0 * * *":
		return now.AddDate(0, 0, 1).Truncate(24 * time.Hour), nil
	case "@weekly", "0 0 * * 0":
		return now.AddDate(0, 0, 7).Truncate(24 * time.Hour), nil
	case "@monthly", "0 0 1 * *":
		return now.AddDate(0, 1, 0).Truncate(24 * time.Hour), nil
	default:
		// For now, default to hourly for unknown patterns
		return now.Add(time.Hour), nil
	}
}