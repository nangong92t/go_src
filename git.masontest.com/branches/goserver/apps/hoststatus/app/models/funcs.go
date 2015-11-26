package models

import (
    "io/ioutil"
    "bytes"
    "strconv"
    "strings"
    "errors"
    "reflect"
    "os/exec"
    "fmt"
)

// get the network io status
func GetNet() ([]NetStatus, error) {
    info, err   := ioutil.ReadFile("/proc/net/dev")
    if err != nil { return nil, err }

    infoS       := string(info)
    lines       := strings.Split(infoS, "\n")

    if len(lines) <= 3 { return nil, errors.New("net has some error!") }

    net     := []NetStatus{}

    for i, line := range lines {
        if i < 2 { continue }
        rece        := NetSub{}
        tran        := NetSub{}

        sectionName := strings.Split(line, ":")
        if len(sectionName) < 2 { continue }
        interName   := strings.TrimLeft(sectionName[0], " ")

        sections    := strings.Split(sectionName[1], " ")
        foundCnt    := 1
        for _, section  := range sections {
            if section == "" { continue }
            curInt, err   := strconv.ParseInt(section, 10, 0)
            if err == nil {
                switch foundCnt {
                    case 1:
                        rece.Bytes      = curInt
                    case 2:
                        rece.Packets    = curInt
                    case 3:
                        rece.Errs       = curInt
                    case 4:
                        rece.Drop       = curInt

                    case 9:
                        tran.Bytes      = curInt
                    case 10:
                        tran.Packets    = curInt
                    case 11:
                        tran.Errs       = curInt
                    case 12:
                        tran.Drop       = curInt
                }
            }
            foundCnt++
        }

        net = append(net, NetStatus{
            Interface:  interName,
            Receive:    &rece,
            Transmit:   &tran,
        })
    }

    return net, nil
}

// get the hardware status.
func GetHd() ([]HDStatus, error) {
    result      := ToRun("df", "-l")
    lines       := strings.Split(result, "\n")
    hds         := []HDStatus{}
    for i, line := range lines {
        if i == 0 { continue }
        sections    := strings.Split(line, " ")
        foundCnt    := 1
        hd          := HDStatus{}
        for _, section  := range sections {
            if section == "" { continue }
            switch foundCnt {
                case 1:
                    hd.Filesystem   = section
                case 2:
                    curInt, _   := strconv.ParseInt(section, 10, 0)
                    hd.Blocks   = curInt
                case 3:
                    curInt, _   := strconv.ParseInt(section, 10, 0)
                    hd.Used     = curInt
                case 4:
                    curInt, _   := strconv.ParseInt(section, 10, 0)
                    hd.Available= curInt
                case 5:
                    hd.Use      = section
                case 6:
                    hd.Mounted      = section
            }

            foundCnt++
        }

        if hd.Use != "" {
            hds     = append(hds, hd)
        }
    }
    return hds, nil
}

func GetSys() (*SysStatus, error) {
    result      := ToRun("vmstat")
    lines       := strings.Split(result, "\n")
    if len(lines) != 4 { return nil, errors.New("result data error!") }

    // check the sysstatus struct.
    sysStat     := SysStatus{}
    sysStatS    := reflect.ValueOf(&sysStat).Elem()

    // get struct's column name
    typeOfT := sysStatS.Type()
    fields  := map[string]int{}
    fieldCnt:= 0
    for ; fieldCnt<sysStatS.NumField(); fieldCnt++ {
        fields[typeOfT.Field(fieldCnt).Name]    = 1
    }

    lineMap     := map[int]string{}
    titleSections   := strings.Split(lines[1], " ")
    foundCnt    :=  1
    for _, section  := range titleSections {
        if section == "" { continue }
        section     = strings.ToUpper(section[:1]) + section[1:]
        if _, ok := fields[section]; ok {
            lineMap[foundCnt]  = section
        }

        foundCnt++
    }

    varSections := strings.Split(lines[2], " ")
    foundCnt    = 1
    for _, section := range varSections {
        if section == "" { continue }
        if key, ok := lineMap[foundCnt]; ok {
            curInt, err := strconv.ParseInt(section, 10, 0)
            if err == nil {
                sysStatS.FieldByName(key).SetInt( curInt )
            }
        }

        foundCnt++
    }
    return &sysStat, nil

}

func ToRun(cmdS string, arg ...string) string {
    cmd := exec.Command(cmdS, arg...)
    result, err := cmd.Output()
    if err != nil {
        return fmt.Sprintf("%s", err)
    }

    return string(result)
}

// Get the memery status.
func GetMem() (*MemeryStatus, error) {
    info, err    := ioutil.ReadFile("/proc/meminfo")
    if err != nil { return nil, err }

    mem     := MemeryStatus{}
    // try to user reflect.
    memS    := reflect.ValueOf(&mem).Elem()

    // get struct's column name
    typeOfT := memS.Type()
    fields  := []string{}
    fieldCnt:= 0
    for ; fieldCnt<memS.NumField(); fieldCnt++ {
        fields  = append(fields, typeOfT.Field(fieldCnt).Name)
    }
    fieldCnt++

    putFieldCnt := 0
    lines   := strings.Split(string(info), "\n")
    for _, line := range lines {
        for _, col := range fields {
            if !strings.HasPrefix(line, col) { continue }

            sections    := strings.Split(line, " ")
            sectionLen  := len(sections)
            if sectionLen <= 2 { return nil, errors.New("Sorry, the memery value of key '" + col + "' is wrong!") }
            stats       := sections[sectionLen-2:sectionLen-1]
            key         := strings.TrimRight(sections[0], ":")

            curInt, err := strconv.ParseInt(stats[0], 10, 0)
            if err != nil { return nil, err }

            memS.FieldByName(key).SetInt( curInt )
            putFieldCnt++
        }

        if putFieldCnt == fieldCnt { break }
    }

    return &mem, nil
}

// Get the Cpu status.
func GetCpu() (*StatStatus, error) {
    info, err    := ioutil.ReadFile("/proc/stat")
    if err != nil { return nil, err }

    stat    := StatStatus{}

    lines   := bytes.Split(info, []byte("\n"))
    cpus    := []CpuStatus{}

    for _, line := range lines {
        if line == nil { continue }

        // Get Cpus stat
        if bytes.HasPrefix(line, []byte("cpu")) {
            cs          := CpuStatus{}
            hasAdded    := 0
            sections    := bytes.Split(line, []byte(" "))

            for i, section  := range sections {
                curString   := string(section)
                if i == 0 || curString == "" { continue }
                curInt, _   := strconv.ParseInt(curString, 10, 0)

                switch hasAdded {
                    case 0:
                        cs.User     = curInt
                    case 1:
                        cs.Nice     = curInt
                    case 2:
                        cs.System   = curInt
                    case 3:
                        cs.Idel     = curInt
                    case 4:
                        cs.IoWait   = curInt
                    case 5:
                        cs.Irq      = curInt
                    case 6:
                        cs.SoftIrq  = curInt
                        break
                }

                hasAdded++
            }


            cpus    = append(cpus, cs)
        } else {
            // get other system stat
            sections    := strings.Split(string(line), " ")
            if len(sections) != 2 { continue }

            curInt, err   := strconv.ParseInt(sections[1], 10, 0)
            if err !=nil { continue }

            switch sections[0] {
                case "ctxt":
                    stat.CTxt   = curInt
                case "btime":
                    stat.BTime  = curInt
                case "processes":
                    stat.Processes  = curInt
                case "procs_running":
                    stat.ProcsRunning   = curInt
                case "procs_blocked":
                    stat.ProcsBlocked   = curInt
            }
        }

        stat.Cpu    = cpus[0]
        stat.SubCpu = cpus[1:]
    }

    return &stat, nil
}

func GetOSStatus() (*OsStatus, error) {
    os  := OsStatus{}
    cpuInfo, err  := GetCpu();
    if err != nil { return nil, err }
    os.Cpu  = cpuInfo

    memInfo, err  := GetMem();
    if err != nil { return nil, err }
    os.Mem  = memInfo

    hdInfo, err   := GetHd();
    if err != nil { return nil, err }
    os.Hd   = hdInfo

    netInfo, err  := GetNet();
    if err != nil { return nil, err }
    os.Net  = netInfo

    sysInfo, err  := GetSys();
    if err != nil { return nil, err }
    os.Sys  = sysInfo

    return &os, nil
}
