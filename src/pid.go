package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println(pid())
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

// Function to get the link value of exe
func get_exe(pid int) (string, error) {
	filename := "/proc/" + strconv.FormatInt(int64(pid), 10) + "/exe"
	return os.Readlink(filename)
}

// Function to display cwd and exe link of a /proc/dir
func pid() []string {
	// Create the slice of PID
	pids, err := pids()
	if err != nil {
		fmt.Println("Error of pids:", err)
		os.Exit(1)
	}
	// For all the PIDS we get, try to retrieve the cwd and exe
	// If there is an error of permission, we continue : go to the next PID
	// If no errors, we retrieve the cwd and exe
	result := make([]string, 1)
	for i := 0; i < len(pids); i++ {
		pid := pids[i]
		pid_val, err_cwd := get_cwd(pid)
		exe_val, err_exe := get_exe(pid)

		if err_cwd != nil {
			continue
		} else {
			str := strconv.Itoa(pid)
			result = append(result, "/proc/"+str+"/cwd:", string(pid_val))
		}
		if err_exe != nil {
			continue
		} else {
			str := strconv.Itoa(pid)
			result = append(result, "/proc/"+str+"/exe:", string(exe_val))
		}
		result = append(result, "\n")

	}
	return result
}
