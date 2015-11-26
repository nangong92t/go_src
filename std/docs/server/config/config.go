package config

import (
    "fmt"
    "flag"
    "log"
    "github.com/kylelemons/go-gypsy/yaml"
)

var (
    serviceFieldNames    = map[string]string {
    //  字段名                  默认值
        "protocol":             "tcp",      // [必填]tcp udp
        "name":                 "",         // [必填]用户相关服务名称
        "port":                 "",         // [必填]监听的端口
        "handle_path":          "",         // [必填]调用实体服务路径
        "child_count":          "1",        // [必填]worker进程数, 默认为1个
        "ip":                   "0.0.0.0",  // [选填]绑定的ip，注意：生产环境要绑定内网ip 不配置默认是0.0.0.0
        "recv_timeout":         "1000",     // [选填]从客户端接收数据的超时时间, 默认为1000毫秒
        "process_timeout":      "30000",    // [选填]业务逻辑处理超时时间
        "send_timeout":         "1000",     // [选填]发送数据到客户端超时时间
        "persistent_connection":"true",     // [选填]是否是长连接, 不配置默认是短链接, 默认为长连接
        "max_requests":         "1000",     // [选填]进程接收多少请求后退出, 不配置默认是0，不退出
        "rpc_secret_key":       "",         // 数据签名用私匙
    }

    configFieldName     = map[string]string {
    //  字段名                  默认值
        "env":                  "dev",      // 开发环境: dev or production
        "worker_user":          "",         // 运行worker的用户,正式环境应该用低权限用户运行worker进程.
    }
)

// 获取配置信息
func GetConfig(configFile string, env string) (map[string]string, map[string]map[string]string) {

    // 装在配置文件对象
    file        := flag.String("file", configFile, "(Simple) YAML file to read")

    // 读取配置文件内容
    configHandle, err := yaml.ReadFile(*file)

    if err != nil {
        log.Fatalf("readfile %q: %s", *file, err)
    }

    serviceTotal, err := configHandle.Count("services")

    mainConfig      := make(map[string]string)
    serviceConfig   := map[string]map[string]string{}

    for i:=0; i<serviceTotal; i++ {
        serviceConf := make(map[string]string)

        // 循环获取Service配置字段值.
        for key, defaultVal := range(serviceFieldNames) {
            serviceVal, err := configHandle.Get("services[" + fmt.Sprintf("%d", i) + "]." + key);
            if err != nil {
                serviceVal = defaultVal
            }

            if serviceVal == "" && defaultVal != "" {
                serviceVal = defaultVal
            }

            serviceConf[ key ]  = serviceVal
        }

        serviceConfig[ serviceConf["name"] ] = serviceConf
    }

    // 获取Env环境变量.
    for key, defaultVal := range(configFieldName) {
        configVal, err  := configHandle.Get(key)
        if err != nil {
            if key == "env" {
                mainConfig[key] = env
            } else {
                mainConfig[key] = defaultVal
            }
        } else {
            if key == "env" {
                env = configVal
            }
            mainConfig[key] = configVal
        }
    }

    log.Printf("current ENV: %s", env)

    return mainConfig, serviceConfig
}

