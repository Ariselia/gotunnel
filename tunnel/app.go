//
//   date  : 2014-06-11
//   author: xjdrew
//
package tunnel

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type Options struct {
	Listen     string
	Server     string // tunnel server or client
	Count      int    // tunnel count underlayer
	Capacity   uint16
	ConfigFile string
	LogLevel   int
	Rc4Key     []byte
}

var options *Options

func init() {
	rand.Seed(time.Now().Unix())
}

type Service interface {
	Start() error
	Reload() error
	Stop()
	Wait()
}

type App struct {
	services []Service
	wg       sync.WaitGroup
}

func (self *App) Add(service Service) {
	self.services = append(self.services, service)
}

func (self *App) Start() error {
	for _, service := range self.services {
		err := service.Start()
		if err != nil {
			return err
		}
	}

	for _, service := range self.services {
		self.wg.Add(1)
		go func(s Service) {
			defer self.wg.Done()
			s.Wait()
			Info("service finish: %v", s)
		}(service)
	}
	return nil
}

func (self *App) Reload() {
	for _, service := range self.services {
		service.Reload()
	}
}

func (self *App) Stop() {
	for _, service := range self.services {
		service.Stop()
	}
}

func (self *App) Wait() {
	self.wg.Wait()
}

func (self *App) Status() {
	Log("num goroutine: %d", runtime.NumGoroutine())
}

func NewApp(o *Options) *App {
	options = o
	return new(App)
}