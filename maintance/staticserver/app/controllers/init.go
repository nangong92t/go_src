// Copyright 2014, SuccessfumMatch Teq Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package controllers

import (
    "github.com/tonycbcd/cron"
    "github.com/revel/revel"
    "maintance/staticserver/app/components/uploader"
)

func init() {
    // to maintance the session cache
    revel.OnAppStart(func() {
        staticTotal := uploader.GetParamString("static.host.total", "0")
        if staticTotal != "0" {
            uploader.GetAllStaticHosts()
            cron := cron.New()
            cron.AddFunc("@every 1s", func() {
                uploader.SyncNewToAllStaticHosts()
            })
            cron.AddFunc("@every 10s", func() {
                uploader.TrySyncFailedStatic()
            })
            cron.Start()
        }
    })
}
