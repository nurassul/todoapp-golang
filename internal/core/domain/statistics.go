package domain

import "time"

type Statistics struct {
	TasksCreated               int
	TasksCompleted             int
	TasksCompletedRate         *float64
	TasksAverageCompletionTime *time.Duration
}

func NewStatistics(
	tasksCreated int,
	tasksCompleted int,
	taskCompletedRate *float64,
	tasksAverageCompletionTime *time.Duration,
) Statistics {
	return Statistics{
		TasksCreated:               tasksCreated,
		TasksCompleted:             tasksCompleted,
		TasksCompletedRate:         taskCompletedRate,
		TasksAverageCompletionTime: tasksAverageCompletionTime,
	}
}
