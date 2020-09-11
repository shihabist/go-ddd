package aggregate

type AggregateRoot struct {
	eventRecorder EventRecorder
	ID            string
}

func (root *AggregateRoot) AddEvent(event interface{}) { root.eventRecorder.Record(event) }

func (root *AggregateRoot) Clear() { root.eventRecorder.Clear() }

func (root *AggregateRoot) Events() []interface{} { return root.eventRecorder.Events() }
