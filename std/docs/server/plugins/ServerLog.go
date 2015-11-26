package plugins

import (
    "time"
    "os"
    "log"
    "fmt"
    "../config"
)

type ServerLog struct {
    // root folder path
    Rp string

    // log file name
    Fn string

    // start time.
    Stime time.Time

    // log type
    Lt string
}

func NewServerLog(lt string) *ServerLog {
    rootFolder  := config.SERVER_BASE + "/" + config.LOG_FOLDER
    // build log root path
    if isExists, _ := IsFolder(rootFolder); !isExists {
        err := MakeFolder(config.SERVER_BASE, config.LOG_FOLDER)
        if err != nil {
            panic("Sorry can not create the log folder:" + rootFolder)
        }
    }

    // build log type path
    if lt != "" {
        ltFolder    := rootFolder + "/" + lt
        if isExists, _ := IsFolder(ltFolder); !isExists {
            err := MakeFolder(rootFolder, lt)
            if err != nil {
                panic("Sorry can not create the log type folder:" + ltFolder)
            }
        }

        rootFolder  = ltFolder
    }

    return &ServerLog {
        Rp: rootFolder,
        Fn: "server.log",
        Stime: time.Now(),
        Lt: lt,
    }
}

func (s *ServerLog) Add(mesg string, v ...interface{}) {
    // build log date path
    curDate := time.Now().Format("2006-01-02")
    dateFolder  := s.Rp + "/" + curDate
    if isExists, _ := IsFolder(dateFolder); !isExists {
        err := MakeFolder(s.Rp, curDate)
        if err != nil {
            log.Fatalln("Sorry can not create the log folder:" + dateFolder)
            return
        }
    }

    // open the log file
    var file *os.File
    var err error

    logFile     := dateFolder + "/" + s.Fn

    if isExists := IsExistsFile(logFile); isExists {
        file, err   = os.OpenFile(logFile, os.O_WRONLY | os.O_APPEND, 0777)
    } else {
        file, err   = os.Create(logFile)
    }
    if err != nil {
        log.Fatalln("Sorry can not open the log file:" + logFile)
        return
    }
    defer file.Close()

    mesg    = time.Now().Format("2006-01-02 15:04:05") + ": " + mesg
    _, err  = file.WriteString(fmt.Sprintf(mesg, v...) + "\n")
    if err != nil {
        log.Fatalln("Sorry can not write log to file:" + logFile)
        return
    }
}

func (s *ServerLog) SetType(lt string) {
    s.Lt    = lt
}
