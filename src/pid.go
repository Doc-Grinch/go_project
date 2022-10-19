package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	display_infos()
}

// Function to open the /proc dir, read its dirs name and make a slice of it
func pids() ([]int, error) {
	f, err := os.Open(`/proc`)
	if err != nil {
		fmt.Println("Could not open the /proc dir - Exiting - :", err)
		os.Exit(1)
	}
	defer f.Close()

	names, err := f.Readdirnames(0)
	if err != nil {
		fmt.Println("Could not list the directories of /proc - Exiting :", err)
	}
	pids := make([]int, 0, len(names))
	// Convert the dir name into int == PIDs
	for _, name := range names {
		if pid, err := strconv.ParseInt(name, 10, 0); err == nil {
			pids = append(pids, int(pid))
		}
	}
	return pids, nil
}

// Function to get the link value of cwd
func get_cwd(pid int) (string, error) {
	filename := "/proc/" + strconv.FormatInt(int64(pid), 10) + "/cwd"
	return os.Readlink(filename)
}

func display_infos() {
	pids, err := pids()
	if err != nil {
		fmt.Println("Error of pids:", err)
		os.Exit(1)
	}
	for i := 0; i < len(pids); i++ {
		pid := pids[i]
		stat, err := get_cwd(pid)
		if err != nil {
			fmt.Println("pid:", pid, err)
			return
		}
		str := strconv.Itoa(pid)
		fmt.Println("/proc/"+str+"/cwd:", string(stat))
	}
}
