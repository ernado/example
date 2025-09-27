package semconv

import (
	"go.uber.org/zap"
)

// TODO: Generate from a semconv specification for log attributes.

// TaskTitle returns zap field with task title.
func TaskTitle(title string) zap.Field {
	return zap.String(AttrTaskTitle, title)
}

// TaskID returns zap field with task ID.
func TaskID(id int64) zap.Field {
	return zap.Int64(AttrTaskID, id)
}
