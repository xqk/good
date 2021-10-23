package worker

// Worker could scheduled by good or customized scheduler
type Worker interface {
	Run() error
	Stop() error
}
