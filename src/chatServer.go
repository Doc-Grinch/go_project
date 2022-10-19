package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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

// Variable de nouvelle session
var newSession session

// Variable du nombre de session crées
var nbConn int = 0

// Variable du nombre de connexion maximum
var maxConn int = 0

func main() {
	server, _ := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if server == nil {
		panic("couldn't start listening....")
	}

	// Demande du nombre de connexion maximales
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Veuillez entrer un nombre de connexion maximale au serveur :")
	scanner.Scan()

	rep1, err := strconv.Atoi(scanner.Text())

	if err != nil {
		fmt.Println("Il faut entrer un nombre !")
		main()
	}

	// Nombre de connexion max set
	maxConn = rep1

	newSession = session{
		connections: []net.Conn{},
		names:       []string{},
	}
	conns := clientConns(server)
	for {

		// ****** A checker
		nbConn++
		fmt.Println(nbConn)
		fmt.Println(maxConn)
		// ******
		go handleConnection(<-conns)
	}

}

/*
 * Ecoute et accepte les connexions clients
 */
func clientConns(listener net.Listener) chan net.Conn {
	channel := make(chan net.Conn)

	// ***** A checker
	if nbConn < maxConn {

		// ----- code déja présent avant
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
		// ------

	} else {

		fmt.Println("Nombre de session maximale atteint")
	}
	// *****

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
