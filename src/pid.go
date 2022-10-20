package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(pid())
}

// Function to open the /proc dir, read its dirs name and make a slice of it, return an array of pid or null
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

// Function to get the link value of cwd, return a string or error
func get_cwd(pid int) (string, error) {
	filename := "/proc/" + strconv.FormatInt(int64(pid), 10) + "/cwd"
	return os.Readlink(filename)
}

// Function to get the link value of exe, return a string or error
func get_exe(pid int) (string, error) {
	filename := "/proc/" + strconv.FormatInt(int64(pid), 10) + "/exe"
	return os.Readlink(filename)
}

// Function to get the mem value of a pid
func get_all_mem(pid int) (string, error) {
	// We convert the pid to int64 to suit the FormatInt
	file, err := os.Open("/proc/" + strconv.FormatInt(int64(pid), 10) + "/smaps")
	// Control that the file was open with no error or exiting
	if err != nil {
		fmt.Println("Could not open the memory file - Exiting - : ", err)
		return "Permission denied", err
	}
	// Used to close the file at the end of main()
	defer file.Close()

	// Define a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	// Scan all the line and retrieve our targets
	array := [5]string{"Size", "Rss", "Pss", "Swap", "SwapPss"}
	result := make([]string, 5)
	i := 0
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), array[i]) && i < len(array) {
			result = append(result, scanner.Text())
			i++
			if i == len(array) {
				break
			}
		}
	}
	string_result := strings.Join(result[:], ",")
	strings.ReplaceAll(string_result, " ", "")
	return string_result, err
}

// Function to display cwd and exe link of a /proc/dir
func pid() []string {
	// Create the slice of PID or exiting if error
	pids, err := pids()
	if err != nil {
		fmt.Println("Error of pids:", err)
		os.Exit(1)
	}
	// For all the PIDS we get, try to retrieve the cwd, exe & memory
	// If there is an error of permission, we continue or display an error then we go to the next PID
	// If no errors, we retrieve the cwd, exe & memory values
	// The return value is an array of string that we display after to the stdout
	result := make([]string, 1)
	for i := 0; i < len(pids); i++ {
		pid := pids[i]
		pid_val, err_cwd := get_cwd(pid)
		exe_val, err_exe := get_exe(pid)
		all_mem, err_all := get_all_mem(pid)
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
		if err_all != nil {
			continue
		} else {
			result = append(result, "all memory values :"+string(all_mem)+" ko")
			result = append(result, "\n")
		}
	}
	return result
}
