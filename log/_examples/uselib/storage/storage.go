package storage

var drivers map[string]DriverFactory

type Driver interface {
	Get(k string) (string, error)
	Set(k string, v string)
}

type DriverFactory func() Driver

func Register(name string, factory DriverFactory) {
	drivers[name] = factory
}

func Get(name string) Driver {
	return drivers[name]()
}

func init() {
	drivers = make(map[string]DriverFactory, 3)
}
