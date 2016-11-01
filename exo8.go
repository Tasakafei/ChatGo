package main

/**
  SUJET : Chat en Go
  AUTEUR : Alexandre Cazala
**/

import (
    "net"
    "bufio"
    "log"
    "fmt"
    "time"
)
type Message struct {
  message_type string
  pseudo string
  content string
}

/* 
 * CLIENT 
 */
type Client struct {
	toHub chan Message
	fromHub chan Message
	pseudo string
	reader *bufio.Reader
	writer *bufio.Writer
	socket net.Conn
}

func (client *Client) Read() {
	for {
		line,err := client.reader.ReadString('\n')
        if err != nil {
           	message := Message {
        		message_type: "DISCONNECTION",
        		pseudo: client.pseudo,
        		content: line,
        	}
        	client.toHub <- message
         	return
        } else {
        	message_type := ""
        	if (client.pseudo == "Unknown") {
        		client.pseudo = line
        		message_type = "CONNECTION"
        	} else if (line == "exit") {
        		message_type = "DISCONNECTION"
        	} else {
        		message_type = "TEXTUAL"
        	}
        	message := Message {
        		message_type: message_type,
        		pseudo: client.pseudo,
        		content: line,
        	}
        	client.toHub <- message
        }

	}
}

func (client *Client) Write() {
	for message := range client.fromHub {
		text := prepareText(message) + "\n"

		client.writer.WriteString(text)
		client.writer.Flush()
	}
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func ClientConstructor(connection net.Conn) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)
	client := &Client {
		toHub: make(chan Message),
		fromHub: make(chan Message),
		pseudo: "Unknown",
		reader: reader,
		writer: writer,
		socket: connection,
	}

	client.Listen()
	return client
}

func prepareText(message Message) string{
	result := ""
	switch {
		case message.message_type == "DISCONNECTION":
			result = "User " + message.pseudo + " has disconnected"
		case message.message_type == "TEXTUAL":
			result = message.pseudo + " : " + message.content
		case message.message_type == "CONNECTION":
			result = message.pseudo + " has joined your chat room"
	}
	return result;
}

/*
 * Chat Room : un hub est une salle de discussion
 */

type PseudoSocket struct {
	pseudo string
	socket net.Conn
}

type ChatRoom struct {
	clients map[string]*Client
	fromClients chan Message
	fromServer chan PseudoSocket
}

func (chatRoom *ChatRoom) Broadcast(message Message) {
	for pseudo,client := range chatRoom.clients {
		if (pseudo != message.pseudo) {
			client.fromHub <- message
		}
	}
}

func (chatRoom *ChatRoom) Join(pseudoSocket PseudoSocket) {
	client := ClientConstructor(pseudoSocket.socket)
	client.pseudo = pseudoSocket.pseudo
	chatRoom.clients[pseudoSocket.pseudo] = client
	// nouveau listener, ce que le client envoie est retransmis dans la channel du chat (donc tout les messages des
	// clients sont regroupés dans la même channel)
	go func() { 
		for {
			fmt.Println("Test "+pseudoSocket.pseudo)
			select {
				case message := <-client.toHub:
					chatRoom.fromClients <- message
				case  <-time.After(time.Second * 20):
					log.Println("TIME OUT de " + client.pseudo )
					chatRoom.Disconnect(client.pseudo)
					return
			} 
		}
	}()
	message := Message {
        		message_type: "CONNECTION",
        		pseudo: pseudoSocket.pseudo,
        		content: "",
    }
    chatRoom.Broadcast(message)
}

func (chatRoom *ChatRoom) Listen() {
	go func() {
		for {
			select {
				case message := <- chatRoom.fromClients:
					chatRoom.HandleMessage(message)
				case message := <- chatRoom.fromServer:
					chatRoom.Join(message)
			}
		}
	}()
}

func (chatRoom *ChatRoom) Disconnect(pseudo string) {
	tmp := chatRoom.clients[pseudo].socket
	delete(chatRoom.clients, pseudo)
	tmp.Close()
	message := Message{
		message_type: "DISCONNECTION",
		pseudo: pseudo,
		content: "",
	}
	chatRoom.Broadcast(message)
}

func (chatRoom *ChatRoom) HandleMessage(message Message) {
	switch {
		case message.message_type == "DISCONNECTION":
			chatRoom.Disconnect(message.pseudo)
		default: 
			chatRoom.Broadcast(message)
	}
}

func ChatRoomConstructor(fromServerChannel chan PseudoSocket) *ChatRoom {
	chatRoom := &ChatRoom{
		clients: make(map[string]*Client),
		fromClients: make(chan Message),
		fromServer: fromServerChannel,
	}

	chatRoom.Listen()

	return chatRoom
}

func handleConnection(conn net.Conn, hub chan PseudoSocket) {
    reader := bufio.NewReader(conn)
    fmt.Fprintln(conn, "Saisissez votre pseudo : ")

    pseudo,err := reader.ReadString('\n')
    if err != nil {
        log.Println(err)
        return
    }
    pseudo = pseudo[0:len(pseudo)-1]
    message := PseudoSocket{
        		pseudo: pseudo,
        		socket: conn,
    }
    hub <- message
}

func main() {
	ServerToChatRoomChannel := make(chan PseudoSocket)
	ChatRoomConstructor(ServerToChatRoomChannel)
    listener, err := net.Listen("tcp", "localhost:1234")
    if err != nil {
        log.Fatal(err)
    }

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println(err)
            continue
        }
        go handleConnection(conn, ServerToChatRoomChannel)    // Rendre la gestion de la connexion concurrente
    }
}