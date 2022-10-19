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

	// Définition des valeurs de newSession
	newSession = session{
		connections: []net.Conn{},
		names:       []string{},
	}

	// Fonction de connexion client
	conns := clientConns(server)

	for {
		//Fonction d'envoi/reception de message
		go handleConnection(<-conns)
	}

}

/*
 * Ecoute et accepte les connexions clients
 */
func clientConns(listener net.Listener) chan net.Conn {
	channel := make(chan net.Conn)

	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				fmt.Printf("couldn't accept client connection")
				continue
			}

			//Condition de limite de connexion
			if nbConn >= maxConn {
				fmt.Printf("Rejecting incoming connection, user limit reached : %d sur %d.\n", nbConn, maxConn)
				//fmt.Println(newSession.connections)
				//fmt.Println(newSession.names)
				err = client.Close() // Fermeture de la connexion client
				if err != nil {      // Gestion des erreur de err
					fmt.Println("Error closing connection after max user limit reached:" + err.Error() + "\n")
				}
				//continue
			} else {
				channel <- client //Envoi du client vers le channel
			}

			//channel <- client
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

	nbConn++

	fmt.Println(newSession.connections)
	fmt.Println(newSession.names)

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
