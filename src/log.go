package douDianSdk

import "time"

type Logger interface {
	Error(...interface{})
	Errorf(string, ...interface{})
	Errorv(interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Infov(interface{})
	Slow(...interface{})
	Slowf(string, ...interface{})
	Slowv(interface{})
	WithDuration(time.Duration) Logger
}
