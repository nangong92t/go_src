package config

const (
    // 支持的协议
    PROTOCOL_TCP      = "tcp"
    PROTOCOL_UDP      = "udp"

    // 服务的各种状态
    STATUS_STARTING   = 1
    STATUS_RUNNING    = 2
    STATUS_SHUTDOWN   = 4
    STATUS_RESTARTING_WORKERS = 8

    // 配置相关
    SERVER_MAX_WORKER_COUNT = 1000
    // 某个进程内存达到这个值时安全退出该进程 单位K
    MAX_MEM_LIMIT = 83886
    // 单个进程打开文件数限制
    MIN_SOFT_OPEN_FILES = 10000
    MIN_HARD_OPEN_FILES = 10000
    // worker从客户端接收数据超时默认时间 毫秒
    WORKER_DEFAULT_RECV_TIMEOUT = 1000
    // worker业务逻辑处理默认超时时间 毫秒
    WORKER_DEFAULT_PROCESS_TIMEOUT = 30000
    // worker发送数据到客户端默认超时时间 毫秒
    WORKER_DEFAULT_SEND_TIMEOUT = 1000
    // 心跳包ping没回复多少次则停止对应进程
    PONG_MISS_LIMIT = 12

    // Basic Server Path
    SERVER_BASE     = "."

    // LOG FOLDER NAME
    LOG_FOLDER      = "logs"
)


