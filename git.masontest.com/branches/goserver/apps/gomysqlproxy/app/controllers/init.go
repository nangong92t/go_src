package controllers

import (
    "sync"
    "github.com/tonycbcd/cron"

    "git.masontest.com/branches/goserver/workers/revel"
    "git.masontest.com/branches/goserver/plugins"
    "git.masontest.com/branches/goserver/apps/gomysqlproxy/app/modules"
)

var (
    isRunning = false
)

func init() {
    revel.OnAppStart(func() {
        // 为了避免在多个app 同时
        if isRunning { return }
        childProcessName    := plugins.GetArg("-server")
        if childProcessName != "HostTracker" { return }

        wg  := &sync.WaitGroup{}
        wg.Add(1)

        cron := cron.New()
        cron.AddFunc("@every 30s", func() { modules.Run() })
        cron.Start()
        isRunning   = true
    })
}

func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}
