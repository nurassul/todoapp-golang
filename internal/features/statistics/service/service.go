package statistics_service

import (
	"context"
	"time"

	"github.com/nurassul/todoapp-golang/internal/core/domain"
)

type StatisticsService struct {
	statsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTasks(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}

func NewStatisticsService(
	statsRepository StatisticsRepository,
) *StatisticsService {
	return &StatisticsService{
		statsRepository: statsRepository,
	}
}
