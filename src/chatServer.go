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
		panic("Ne peut démarrer....")
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
				fmt.Printf("Ne peut accepter la connexion client")
				continue
			}

			//Condition de limite de connexion
			if nbConn >= maxConn {
				fmt.Printf("Rejet de connexion entrante, limite d'utilisateur atteinte : %d sur %d.\n", nbConn, maxConn)
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
		//Si erreur il y à (déconnexion) et retrait des valeurs de l'utilisateur
		if err != nil {
			fmt.Println("Le client " + clientName + " s'est déconnécté")
			for i := 0; i < len(newSession.names); i++ {
				if newSession.names[i] == clientName {
					newSession.names = remove(newSession.names, i)
					newSession.connections = removeConn(newSession.connections, i)
					nbConn--
					fmt.Println(newSession.connections)
					fmt.Println(newSession.names)
				}
			}
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

// Fonction remove de valreur string d'une array
func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

// Fonction remove de valreur net.Conn d'une array
func removeConn(slice []net.Conn, s int) []net.Conn {
	return append(slice[:s], slice[s+1:]...)
}
