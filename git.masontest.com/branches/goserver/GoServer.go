package main

import (
    "os"
    "./config"
    "./core"
)

func main() {
    // 默认配置文件在当前运行命令目录下的config.yml
    configFile      := os.Getenv("GOPATH") + "/src/git.masontest.com/branches/goserver/config.yml"
    env             := "dev"            // 默认开发环境为测试开发环境

    // 获取配置信息.
    mainConfig, serviceConfig   := config.GetConfig(configFile, env)

    // 开始服务运行.
    core.NewGoServer(mainConfig, serviceConfig).Run()
}
