package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

// Constante du port utlilisé
const PORT = 8080

// Sessions de tchat
type session struct {
	// Connexion enregistrées
	connections []net.Conn
	// Entrer le nom de l'utilisateur
	names []string
}

// Nouvelle structure pour chaques sesssions
type server struct {
	currentSession session
}

// Varible de nouvelle session
var newSession session

func main() {
	server, _ := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if server == nil {
		panic("couldn't start listening....")
	}
	newSession = session{
		connections: []net.Conn{},
		names:       []string{},
	}
	conns := clientConns(server)
	for {
		go handleConnection(<-conns)
	}
}

/*
 * Ecoute et accepte les connexion clients
 */
func clientConns(listener net.Listener) chan net.Conn {
	channel := make(chan net.Conn)
	go func() {
		for {
			client, _ := listener.Accept()
			if client == nil {
				fmt.Printf("couldn't accept client connection")
				continue
			}
			channel <- client
		}
	}()
	return channel
}

/*
 * Permet de nouvellles connexions
 * Sauvegarde le nom utilisateur, attends les messages et les diffuses
 */
func handleConnection(client net.Conn) {
	reader := bufio.NewReader(client)
	//Recois le nom utilisateur
	buff := make([]byte, 512)
	clientNameb, _ := client.Read(buff)
	clientName := string(buff[0:clientNameb])
	newSession.names = append(newSession.names, clientName)
	newSession.connections = append(newSession.connections, client)

	for {
		//Recois le message utilisateur
		line, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		//Diffuse le message
		message := clientName + ":" + string(line)
		for _, currentClient := range newSession.connections {
			if currentClient != nil {
				currentClient.Write([]byte(message))
			}
		}

	}
}