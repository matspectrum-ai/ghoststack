package workflow

import "context"

type Step struct {
	ID    string
	Run   func(context.Context) error
	After []string
}

type Workflow struct {
	ID    string
	Steps []Step
}

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Run(ctx context.Context, wf Workflow) error {
	return nil
}
