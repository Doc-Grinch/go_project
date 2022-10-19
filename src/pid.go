package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

// Function to get the mem value of a pid
func get_mem(pid int) (string, error) {
	data, err := os.ReadFile("/proc/" + strconv.FormatInt(int64(pid), 10) + "/stat")
	result := strings.Fields(string(data))
	mem_bytes, err_conv := strconv.Atoi(result[22])
	if err_conv != nil {
		fmt.Println("Error during conversion - Exiting :", err_conv)
		os.Exit(1)
	}
	mem_mb := mem_bytes / (1024 * 1024)
	mem_res := strconv.Itoa(mem_mb)
	return mem_res, err
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
		mem_val, err_mem := get_mem(pid)
		result = append(result, "\n")
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
			result = append(result, "----- /proc/"+str+"/exe:", string(exe_val))
		}
		if err_mem != nil {
			continue
		} else {
			result = append(result, "----- mem :"+string(mem_val)+" ko")
		}

	}
	return result
}
