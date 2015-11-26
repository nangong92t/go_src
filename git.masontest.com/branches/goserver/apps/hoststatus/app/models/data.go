package models

type CpuStatus struct {
    User    int64       // 从系统启动开始累计到当前时刻，用户态的CPU时间（单位：jiffies） ，不包含 nice值为负进程。1jiffies=0.01秒
    Nice    int64       // 从系统启动开始累计到当前时刻，nice值为负的进程所占用的CPU时间（单位：jiffies） 
    System  int64       // 从系统启动开始累计到当前时刻，核心时间（单位：jiffies）
    Idel    int64       // 从系统启动开始累计到当前时刻，除硬盘IO等待时间以外其它等待时间（单位：jiffies）
    IoWait  int64       // 从系统启动开始累计到当前时刻，硬盘IO等待时间（单位：jiffies）
    Irq     int64       // 从系统启动开始累计到当前时刻，硬中断时间（单位：jiffies）
    SoftIrq int64       // 从系统启动开始累计到当前时刻，软中断时间（单位：jiffies）
}

type StatStatus struct {
    Cpu         CpuStatus
    SubCpu      []CpuStatus
    CTxt        int64   // 累计上下文交换次数
    BTime       int64   // 系统总运行时间，单位为秒
    Processes   int64   // 所有创建的任务数
    ProcsRunning    int64   // 当前运行队列的任务的数目
    ProcsBlocked    int64   // 当前被阻塞的任务的数目
}

type MemeryStatus struct {
    MemTotal    int64   // kb
    MemFree     int64   // kb
    Buffers     int64
    Cached      int64
    SwapTotal   int64   // kb
    SwapCached  int64
    SwapFree    int64   // kb
}

type HDStatus struct {
    Filesystem  string
    Blocks      int64
    Used        int64
    Available   int64
    Use         string
    Mounted     string
}

type NetSub struct {
    Bytes       int64
    Packets     int64
    Errs        int64
    Drop        int64
}

type NetStatus struct {
    Interface   string
    Receive     *NetSub
    Transmit    *NetSub
}

type SysStatus struct {
    R           int64
    B           int64
    Swpd        int64
    Free        int64
    Buff        int64
    Cache       int64
    Si          int64
    So          int64
    Bi          int64
    Bo          int64
    In          int64
    Cs          int64
    Us          int64
    Sy          int64
    Id          int64
    Wa          int64
    St          int64
}

type OsStatus struct {
    Cpu         *StatStatus
    Mem         *MemeryStatus
    Hd          []HDStatus
    Net         []NetStatus
    Sys         *SysStatus
}

