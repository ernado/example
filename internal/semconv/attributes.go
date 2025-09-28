package semconv

import (
	"go.opentelemetry.io/otel/attribute"
)

// TODO: Generate from a semconv specification

const (
	AttrTaskID    = "task.id"
	AttrTaskTitle = "task.title"

	AttrOperation = "operation"
	AttrResult    = "operation.result"
	AttrSystem    = "operation.system"
)

func Operation(operation string) attribute.KeyValue {
	return attribute.String(AttrOperation, operation)
}

func System(system string) attribute.KeyValue {
	return attribute.String(AttrSystem, system)
}

const (
	DBResultOkValue       = "ok"
	DBResultErrorValue    = "error"
	DBResultCanceledValue = "canceled"
)

func ResultOk() attribute.KeyValue {
	return attribute.String(AttrResult, DBResultOkValue)
}

func ResultError() attribute.KeyValue {
	return attribute.String(AttrResult, DBResultErrorValue)
}

func ResultCanceled() attribute.KeyValue {
	return attribute.String(AttrResult, DBResultCanceledValue)
}
