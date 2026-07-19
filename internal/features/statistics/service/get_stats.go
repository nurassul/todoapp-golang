package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/nurassul/todoapp-golang/internal/core/domain"
	core_errors "github.com/nurassul/todoapp-golang/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {

	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf(
				"'to' must be after 'from': %w",
				core_errors.ErrInvalidArgument)
		}
	}

	tasks, err := s.statsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf(
			"get tasks from repository: %w", err)
	}

	stats := calcStats(tasks)

	return stats, nil
}

func calcStats(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.Statistics{
			TasksCreated:               0,
			TasksCompleted:             0,
			TasksCompletedRate:         nil,
			TasksAverageCompletionTime: nil,
		}
	}

	tasksCreated := len(tasks)

	tasksCompleted := 0
	var totalCompletedDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}

		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletedDuration += *completionDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAverageCompletionTime *time.Duration
	if tasksCompleted > 0 && totalCompletedDuration != 0 {
		avg := totalCompletedDuration / time.Duration(tasksCompleted)

		tasksAverageCompletionTime = &avg
	}

	return domain.Statistics{
		TasksCreated:               tasksCreated,
		TasksCompleted:             tasksCompleted,
		TasksCompletedRate:         &tasksCompletedRate,
		TasksAverageCompletionTime: tasksAverageCompletionTime,
	}

}
