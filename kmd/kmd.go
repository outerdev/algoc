package kmd

import (
	. "fmt"
	. "strconv"
	. "strings"

	"os"
	"os/exec"

	ps "github.com/mitchellh/go-ps"

	. "github.com/outerdev/algoc/errors"
	"github.com/outerdev/algoc/path"
)

const (
	kmdLogFileName = "kmd.log"
	kmdLogFilePerm = 0640

	algocKmdPIDFile = "algoc.kmd.pid"
)

func isKmd() bool {
	if len(os.Args) <= 1 {
		return false
	}

	_, err := NewVersionFromString(os.Args[1])
	return err == nil
}

func kmdPidFromFile(dataDir string) (int, bool) {

	pidFile := dataDir + "/" + algocKmdPIDFile

	pidStr, err := path.ReadStringFromFile(pidFile)
	if err != nil {
		return 0, false
	}

	pid, err := Atoi(TrimSpace(pidStr))
	if err != nil {
		return 0, false
	}

	return pid, true
}

func IsFileNotExist(err error) bool {
	return err == ErrFileNotFound
}

func removeKmdPidFile(dataDir string) error {
	return path.DeleteFile(dataDir + "/" + algocKmdPIDFile)
}

func removeKmdVersionFile(dataDir string) error {
	return path.DeleteFile(dataDir + "/" + algocKmdVersionFile)
}

func startDaemon(dataDir string) {
	if !path.DirExists(dataDir) {
		if err := os.Mkdir(dataDir, 0755); err != nil {
			panic(err)
		}
	}

	if err := path.WriteStringToFile(dataDir+"/"+algocKmdVersionFile, algocKmdVersion); err != nil {
		panic(err)
	}

	cmd := exec.Command(os.Args[0], "-"+algocKmdVersionFlag+"="+algocKmdVersion)
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	pid := cmd.Process.Pid
	if err := path.WriteStringToFile(dataDir+"/"+algocKmdPIDFile, Itoa(pid)); err != nil {
		cmd.Process.Kill()
		panic(err)
	}
}

func isKmdRunning() bool {
	return true
}

func kmdPids() ([]int, error) {

	procs, err := ps.Processes()
	if err != nil {
		return []int{}, nil
	}

	kmdPids := []int{}
	for _, proc := range procs {
		if proc.Executable() == "algoc" && proc.Pid() != os.Getpid() {
			kmdPids = append(kmdPids, proc.Pid())
		}
	}

	return kmdPids, nil
}

func isKmdFilePidAmongPids(kmdPids []int, currentKmdPid int) bool {
	for _, kmdPid := range kmdPids {
		if kmdPid == currentKmdPid {
			return true
		}
	}
	return false
}

func Start(configDir string) {

	dataDir := configDir + "/data"

	if isKmd() {
		// TODO: Setup logging, use it here
		// Println("Starting kmd...")
		startKmd(dataDir, uint64(0))
		return
	}

	pids, err := kmdPids()
	if err != nil {
		panic(err)
	}

	killProcesses := func(pids []int) {
		for _, pid := range pids {
			if proc, _ := os.FindProcess(pid); proc != nil {
				proc.Signal(os.Kill)
			}
		}
	}

	removeKmdFilePid := func(pids []int, filePid int) []int {
		newPids := []int{}
		for _, pid := range pids {
			if pid != filePid {
				newPids = append(newPids, pid)
			}
		}
		return newPids
	}

	removeKmdDaemonFiles := func(dataDir string) {
		if err := removeKmdPidFile(dataDir); err != nil {
			if !IsFileNotExist(err) {
				// TODO: Setup logging, use it here
				Printf("Error: %v\n", err)
			}
		}
		if err := removeKmdVersionFile(dataDir); err != nil {
			if !IsFileNotExist(err) {
				// TODO: Setup logging, use it here
				Printf("Error: %v\n", err)
			}
		}

	}

	if kmdFilePid, ok := kmdPidFromFile(dataDir); ok {
		isKmdCurrent := isKmdFilePidAmongPids(pids, kmdFilePid)
		if isKmdCurrent {
			pids = removeKmdFilePid(pids, kmdFilePid)
		} else {
			removeKmdDaemonFiles(dataDir)
		}

		killProcesses(pids)
		// If kmd is already running then we are done
		// otherwise we'll have to start it
		if isKmdCurrent {
			return
		}
	} else {
		killProcesses(pids)
		removeKmdDaemonFiles(dataDir)
	}

	// TODO: Setup logging, use it here
	// Println("Done checks. Starting the daemon...")
	startDaemon(dataDir)
}
