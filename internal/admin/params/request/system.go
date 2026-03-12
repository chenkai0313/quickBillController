package request

type SystemSummaryRequest struct {
	StartAt int64 `form:"start_at" binding:"required"`
	EndAt   int64 `form:"end_at" binding:"required"`
}
