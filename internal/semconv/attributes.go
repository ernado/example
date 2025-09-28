package semconv

import (
	"go.opentelemetry.io/otel/attribute"
)

// TODO: Generate from a semconv specification

const (
	AttrTaskID    = "task.id"
	AttrTaskTitle = "task.title"

	AttrDBOperation = "db.operation"
	AttrDBResult    = "db.operation.result"
)

func DBOperation(operation string) attribute.KeyValue {
	return attribute.String(AttrDBOperation, operation)
}

const (
	DBResultOkValue       = "ok"
	DBResultErrorValue    = "error"
	DBResultCanceledValue = "canceled"
)

func DBResultOk() attribute.KeyValue {
	return attribute.String(AttrDBResult, DBResultOkValue)
}

func DBResultError() attribute.KeyValue {
	return attribute.String(AttrDBResult, DBResultErrorValue)
}

func DBResultCanceled() attribute.KeyValue {
	return attribute.String(AttrDBResult, DBResultCanceledValue)
}
