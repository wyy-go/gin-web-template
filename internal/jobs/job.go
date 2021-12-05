package jobs

import (
	"fmt"
	"github.com/robfig/cron"
)

var (
	testSpec = "*/5 * * * * *"
)

type JobInfo struct {
	Spec string
	Job  func()
}

var (
	c    = cron.New()
	Jobs = map[string]JobInfo{}
)

func addJob(name string, spec string, job func()) {
	Jobs[name] = JobInfo{Spec: spec, Job: job}
	if spec == "@manual" {
		return
	}
	if err := c.AddFunc(spec, job); err != nil {
		panic(err)
	}
}

func Setup() {
	addJob("test", testSpec, test)
}

func Start() {
	c.Start()
}

func test() {
	fmt.Println("test ...")
}