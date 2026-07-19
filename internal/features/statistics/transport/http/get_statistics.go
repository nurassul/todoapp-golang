package statistics_trasport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nurassul/todoapp-golang/internal/core/domain"
	core_logger "github.com/nurassul/todoapp-golang/internal/core/logger"
	core_http_request "github.com/nurassul/todoapp-golang/internal/core/transport/http/request"
	core_http_response "github.com/nurassul/todoapp-golang/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"`
	TasksCompleted             int      `json:"tasks_completed"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time"`
}

func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, from, to, err := getUserIDFromQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/from/to query params",
		)

		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get stats",
		)
	}

	response := toDTOFromDomain(statistics)

	responseHandler.JSONResponse(response, http.StatusOK)

}

func toDTOFromDomain(stats domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if stats.TasksAverageCompletionTime != nil {
		duration := stats.TasksAverageCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               stats.TasksCreated,
		TasksCompleted:             stats.TasksCompleted,
		TasksCompletedRate:         stats.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

func getUserIDFromQueryParam(r *http.Request) (*int, *time.Time, *time.Time, error) {
	userID, err := core_http_request.GetIntQueryParam(r, "user_id")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	from, err := core_http_request.GetDateQueryParam(r, "from")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err := core_http_request.GetDateQueryParam(r, "to")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	return userID, from, to, nil
}
