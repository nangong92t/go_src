package core

import (
    "fmt"
    "syscall"
    "os"
    "os/exec"
    "os/signal"
    "log"
    "net"
    "time"
    "strconv"
    "runtime"
//	"github.com/sevlyar/go-daemon"
    "../workers"
    "../plugins"
)

type GoServer struct {
    MainCfg map[string]string
    ServerCfg map[string]map[string]string
    Log *plugins.ServerLog
    Ch chan os.Signal
}

var (
    // server统计信息
    serverStatusInfo = map[string]interface{}{
        "start_time": int64(0),
        "err_info": []interface{}{},
    }

    // worker config cache.
    workerConfigs   = map[string]string{}
)

func NewGoServer(mainConfig map[string]string, serviceConfig map[string]map[string]string) *GoServer {

    return &GoServer{
        MainCfg: mainConfig,
        ServerCfg: serviceConfig,
        Log: plugins.NewServerLog(""),
        Ch: make(chan os.Signal),
    }
}

func (s *GoServer) Run() {
    fmt.Printf("%#v\n", s.MainCfg)

    // 检查Server环境
    s.checkEnv()

    // 使之成为daemon进程
    s.daemonize()

    // 安装相关信号
    s.installSignal()

    // 创建监听进程
    if !s.createSocketsAndListen() {
        return
    }
}

// 检查必要环境
func (s *GoServer) checkEnv() {
    if syscall.Getppid() == 1 {
        log.Fatal("server already started")
        os.Exit(1)
    }

    // 检查指定的worker用户是否合法
    s.checkWorkerUserName()

    // 检查log目录是否可读
    s.checkLogWriteAble()

    // 检查配置和语法错误等
    s.checkWorkersConfig();

    // 检查文件限制
    s.checkLimit();
}

// 检查log目录是否可读
func (s *GoServer) checkLogWriteAble() {
    // TODO
}

// 检查启动worker进程的的用户是否合法 
func (s *GoServer) checkWorkerUserName() {
    // TODO
}

// 检查配置和语法错误等
func (s *GoServer) checkWorkersConfig() {
    // TODO
}

// 检查文件限制
func (s *GoServer) checkLimit() {
    // TODO
}

// 使之脱离终端，变为守护进程
func (s *GoServer) daemonize() {
    if plugins.IsInArg("-child") { return }

    args := os.Args[1:]
    args    = append(args, "-child=true")

    for sn, _ := range s.ServerCfg {
        newArgs    := append(args, "-server="+sn)
        cmd := exec.Command(os.Args[0], newArgs...)
        cmd.Start()

        msg := "Starting the server "+sn+" on [PID]"
        fmt.Println(msg, cmd.Process.Pid)
        s.Log.Add(msg, cmd.Process.Pid)
    }

    // 等待子进程启动.
    time.Sleep(1*time.Second)

    fmt.Println("GoServer starting successfully.")
}

// 安装相关信号
func (s *GoServer) installSignal() {
    // Handle SIGINT and SIGTERM.
    signal.Notify(s.Ch, syscall.SIGINT, syscall.SIGTERM)
}

// 创建监听进程
func (s *GoServer) createSocketsAndListen() bool {
    if !plugins.IsInArg("-child") { return false }
    childProcessName    := plugins.GetArg("-server")
    config              := s.ServerCfg[childProcessName]

    if config["name"] == "" {
        s.Log.Add("No found out the worker config for:" + childProcessName)
        return false
    }

    port    := ""
    if config["port"] == "" {
        return false
    } else {
        port = config["port"]
    }

    protocol:= config["protocol"]
    ip      := config["ip"]

    // 开启多核并行执行.
    NCPU := runtime.NumCPU()
    runtime.GOMAXPROCS(NCPU)

    switch config["protocol"] {
        case "tcp":
            return s.createTcpSocketsAndListen(protocol, ip, port, &config)
        case "udp":
            return s.createUdpSocketsAndListen(protocol,ip, port, &config)
        default:
            s.Log.Add("sorry can not support this protocol: %s", config["protocol"])
            return false
    }
}

func (s *GoServer) createTcpSocketsAndListen(protocol, ip, port string, config *map[string]string) bool {
    // 创建监听socket
    laddr, err := net.ResolveTCPAddr(protocol, ip+":"+port)
    if nil != err {
        s.Log.Add("add port error: ", err)
        return false
    }
    listener, err := net.ListenTCP(protocol, laddr)
    if nil != err {
        s.Log.Add("add listener error: ", err)
        return false
    }
    s.Log.Add((*config)["name"] + " listening on", listener.Addr())

    // Make a new service and send it into the background.
    chSync          := make(chan int)
    childCount, _   := strconv.Atoi((*config)["child_count"])
    for i:=1; i<=childCount; i++ {
        go func() {
            workers.NewRpcWorker(i, config, listener, chSync).Run()
        }()
    }
    for i:=1; i<=childCount; i++ {
        chSync <- i
    }

    log.Println(<-s.Ch)
    return true
}

func (s *GoServer) createUdpSocketsAndListen(protocol, ip, port string, config *map[string]string) bool {
    // TODO
    return false
}
