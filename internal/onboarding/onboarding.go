package onboarding

import "fmt"

type Step struct {
	ID       string
	Title    string
	Complete bool
}

type Onboarding struct {
	steps []Step
}

func NewOnboarding(steps []Step) *Onboarding {
	return &Onboarding{steps: steps}
}

func (o *Onboarding) Next() (Step, error) {
	for _, step := range o.steps {
		if !step.Complete {
			return step, nil
		}
	}
	return Step{}, fmt.Errorf("onboarding complete")
}

func (o *Onboarding) Complete(id string) error {
	for i := range o.steps {
		if o.steps[i].ID == id {
			o.steps[i].Complete = true
			return nil
		}
	}
	return fmt.Errorf("step not found: %s", id)
}

func (o *Onboarding) Status() []Step {
	return append([]Step(nil), o.steps...)
}
