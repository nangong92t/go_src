package main

import (
    "os"
    "git.masontest.com/branches/goserver/config"
    r "git.masontest.com/branches/goserver/core/running"
    "git.masontest.com/branches/goserver/plugins"
)

// link worker code.
func linkCode() bool {
    configFile      := os.Getenv("GOPATH") + "/src/git.masontest.com/branches/goserver/config.yml"     // 默认配置文件在当前运行命令目录下的config.yml
    env             := "dev"

    // 获取配置信息.
    mainConfig, serviceConfig   := config.GetConfig(configFile, env)

    // log
    log         := plugins.NewServerLog("")

    appNames    := []string{}
    for sn, config := range serviceConfig {
        if config["handle_path"] != "" {
            appNames = append(appNames, sn)
        }
    }

    r.BuildAppLink(appNames)

    if len(appNames) == 0 { return false }

    for sn, config := range serviceConfig {
        // to link apps code to runtime
        r.Build(mainConfig["env"], sn, config["handle_path"], log)
    }

    return true
}

func main() {
    linkCode()
}
