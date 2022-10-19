package main

import ( // Librairies nécéssaires
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	connection, err := net.Dial("tcp", "192.168.207.130:8080")
	//Ferme la connexion et nettoie le tchat à la sortie de l'utilisateur
	if err != nil {
		fmt.Println("Can not connect to the server")
		return
	}
	defer cleanUp(connection)
	//Envois du nom utilisateur au serveur
	sendUserName(connection)

	fmt.Printf("*****************TChat Botté*****************\n")
	//Lit les messages depuis la console pour les envoyer au serveur
	go messageWriter(connection)

	//Lit les messages reçu depuis le serveur
	messageReader(connection)

}

/*
 * Ferme la connexion en cas d'erreur
 */
func cleanUp(clientConnection net.Conn) {
	clientConnection.Close()
	os.Exit(0)
}

/*
 * Entrer son nom d'utilisateur et l'envois au serveur
 */
func sendUserName(client net.Conn) {
	fmt.Printf("Welcome to chat rooms. Please enter your name.\n")
	inputReader := bufio.NewReader(os.Stdin)
	name, _, error := inputReader.ReadLine()
	if error != nil {
		fmt.Println("Can not read user name")
		cleanUp(client)
		return
	}
	userName := string(name)
	client.Write([]byte(userName))
}

/*
 * Attends l'input du message de l'utilisateur et l'envois au serveur
 */
func messageReader(client net.Conn) {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		message, error := inputReader.ReadString('\n')
		if error != nil {
			fmt.Println("Can not read user message")
			cleanUp(client)
			break
		}
		//Envois le message au serveur
		client.Write([]byte(message))
	}
}

/*
 * Attends que la serveur diffuse un message pour pouvoir l'afficher dans la console
 */
func messageWriter(client net.Conn) {
	reader := bufio.NewReader(client)
	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Error read from server")
			cleanUp(client)
			break
		}
		//Affiche le message reçu
		fmt.Printf(string(message))
	}
}
