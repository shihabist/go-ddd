package behaviour

import (
	"context"
	"log"

	"orderContext/core/mediator"
)

type Logger struct{}

func NewLogger() *Logger { return &Logger{} }

func (l *Logger) Process(ctx context.Context, cmd interface{}, next mediator.Next) error {

	log.Println("Pre Process!")

	result := next(ctx)

	log.Println("Post Process")

	return result
}
