package workers

import (
)

const (
    // 最大buffer长度
    MAX_BUFFER_SIZE     = 524288
)

var (
    // 上次写数据到磁盘的时间
    logLastWriteTime    = int64(0)
    stLastWriteTime     = int64(0)
    lastClearTime       = int64(0)

    // log数据
    logBuffer           = make([]byte, MAX_BUFFER_SIZE)
    statisticData       = []string{}

    // 与统计中心通信所用的协议
    protocolToCenter    = "udp"

    // 多长时间写一次log数据
    stSendTimeLong      = 300

    // 多长时间清除一次统计数据
    clearTimeLong       = 86400

    // 日志过期时间 14days
    logExpTimeLong      = 1296000

    // 统计结果过期时间 14days
    stExpTimeLong       = 1296000

)

type StatisticWorker struct {

}

func NewStatisticWorker() *StatisticWorker {
    return &StatisticWorker {

    }
}
