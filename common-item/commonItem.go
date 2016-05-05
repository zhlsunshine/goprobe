package commonItem 

import (
    // "os"
    // "bytes"
    "bufio"
    "strconv"
    "strings"
    "goprobe/model"
    "github.com/codeskyblue/go-sh"
    log "github.com/cihub/seelog"
)

func MonitorCommonItem(flag bool) {
    MonitorCPU(flag)
    // MonitorMemory()
    // MonitorDisk()
} 

func MonitorCPU(flag bool) {
    log.Info("MonitorCPU start!")
    coreNum := -1 
    session := sh.NewSession()
    out, err := session.Command("cat", "/proc/stat").Output()
    if err != nil {
        panic(err)
    }
    var arrCPU []model.CPURaw
    scanner := bufio.NewScanner(strings.NewReader(string(out)))
    for scanner.Scan() {
        line :=  scanner.Text()
        fields := strings.Fields(line)
        if len(fields) > 0 && strings.Contains(fields[0], "cpu") {
            var cpuitme model.CPURaw
            coreNum++
            cpuitme.CPUNum = string(fields[0])
            parseCPUFields(fields, &cpuitme)
            arrCPU = append(arrCPU, cpuitme)
        }
    }

    if !flag {
        model.PreCPU = arrCPU
        return
    }
    var arrCPUInfo []model.CPUInfo
    for i := 0; i < len(arrCPU); i++ {
        stat := arrCPU[i]
        var cpuitem model.CPUInfo
        cpuitem.CPUNum = stat.CPUNum
        total := float32(stat.Total - model.PreCPU[i].Total)
        cpuitem.User = float32(stat.User -  model.PreCPU[i].User) / total * 100
        cpuitem.Nice = float32(stat.Nice -  model.PreCPU[i].Nice) / total * 100
        cpuitem.System = float32(stat.System -  model.PreCPU[i].System) / total * 100
        cpuitem.Idle = float32(stat.Idle -  model.PreCPU[i].Idle) / total * 100
        arrCPUInfo = append(arrCPUInfo, cpuitem)
    }
    model.PreCPU = arrCPU

    log.Info("The number of terminal cpu core is ", coreNum,)
    for l := 0; l < len(arrCPUInfo); l++ {
        cpuInfo := arrCPUInfo[l]
        log.Info("--------------------------------")
        log.Info("CPU Name: ", cpuInfo.CPUNum)
        log.Info("User: ", cpuInfo.User)
        log.Info("Nice: ", cpuInfo.Nice)
        log.Info("System: ", cpuInfo.System)
        log.Info("Idle: ", cpuInfo.Idle)
        log.Info("--------------------------------")
    }
    log.Info("MonitorCPU end!")
}

func parseCPUFields(fields []string, stat *model.CPURaw) {
    numFields := len(fields)
    for i := 1; i < numFields; i++ {
        val, err := strconv.ParseUint(fields[i], 10, 64)
        if err != nil {
            continue
        }
        stat.Total += val
        switch i {
            case 1:
                stat.User = val
            case 2:
                stat.Nice = val
            case 3:
                stat.System = val
            case 4:
                stat.Idle = val
            case 5:
                stat.Iowait = val
            case 6:
                stat.Irq = val
            case 7:
                stat.SoftIrq = val
            case 8:
                stat.Steal = val
            case 9:
                stat.Guest = val
        }
    }
}
