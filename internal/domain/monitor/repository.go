package monitor

type Repository interface {
	Save(monitor *URLMonitor) error
	FindByID(id string) (*URLMonitor, error)
	FindAll() ([]*URLMonitor, error)
	Delete(id string) error
	Update(monitor *URLMonitor) error
}
