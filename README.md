# ChatGo
## Description du projet

Projet d'étude dans le cadre de la matière techniques de programmation concurrente modernes à l'école Polytech Nice-Sophia en 5e année des sciences informatiques.

Auteur : Alexandre CAZALA

## Scénario et contraintes
Le projet doit respecter plusieurs contraintes que nous numéroterons pour y faire référence plus tard :

1. Lorsqu'un utilisateur se connecte sur le serveur, on commence par lui demander quel est son pseudo. Cette information est communiquée aux autres participants.
2. Lorsqu'un utilisateur tape une ligne, celle-ci est envoyée (précédée de son pseudo) aux autres participants (mais pas à lui).
3. Si un utilisateur se déconnecte du serveur, les autres sont prévenus.
4. Si un utilisateur est inactif pendant un certain temps, le serveur le déconnectera, en indiquant aux autres utilisateurs le pseudo de l'utilisateur qui vient d'être déconnecté.

## Principe global et propriétés
Pour ce projet, afin de me former j'ai préféré développer un projet complet. Je défend donc la modularisation (et donc la réutilisation de mon code) c'est pourquoi j'ai beaucoup de structures définies.

### ChatRoom
```go
type ChatRoom struct {
	clients map[string]*Client
	fromClients chan Message
	fromServer chan PseudoSocket
}
```

Nos utilisateurs sont regroupés dans des chatRooms. Ces clients sont identifiés par un pseudo dans une hashmap, ce pseudo doit donc être unique. Ainsi à l'avenir nous pourrions rajouter des propriétés propres aux chatRooms telles qu'un nombre d'utilisateurs max ou des permissions spéciales (exemple, tout ceux qui s'appellent GéantVert pourraient avoir accès à des commandes spéciales).

Une chatRoom a deux canaux. Le premier **fromClients** est celui où arrive les messages de nos clients. Le deuxième, **fromServer** est celui où les messages du serveur mère arrive (exemple, nouvel utilisateur à rajouter au salon et à l'avenir potentiellement d'autres messages tels que des "Redémarrage prévu dans deux minutes".


Une chat room a un constructeur **ChatRoomConstructor** mais aussi plusieurs méthodes dont leur nom définissent leur objectif :
1. **Join** : permet de créer un nouvel utilisateur (la structure de données) et de le rajouter au salon. Cette méthode est appelée par le connectionHandler qui a déjà recueillie le nom du visiteur.
2. **Broadcast** : Broadcast un message à tout les utilisateurs du salon sauf celui qui l'a envoyé.
3. **HandleMessage** : Permet de factoriser le traitement des messages, si c'est un message de déconnection on débranche l'utilisateur de la map et on ferme sa socket. Sinon on broadcast le message. 
4. **Disconnect** : Permet de déconnecter un client et de prévenir tout le monde.
5. **Listen** : permet de mettre le serveur sur écoute des messages du client et des messages du serveur.

### Client
```go
type Client struct {
	toHub chan Message
	fromHub chan Message
	pseudo string
	reader *bufio.Reader
	writer *bufio.Writer
	socket net.Conn
}
```
 
Un client est une modélisation du client réel. Cette structure permet de garder plusieurs méthodes effectuable sur un client ainsi que de transporter des attributs. **toHub** est le canal pour envoyer des messages à la **ChatRoom** dans laquelle le client se trouve (un client ne peut être que dans une seule chatRoom). **fromHub** est le canal de réception des messages venant de la ChatRoom. **Pseudo** est le pseudo de l'utilisateur. reader est le moyen d'écriture du client (chaque client peut avoir sa propre interface qui difère à l'avenir). socket est la socket du client (pour nous permettre de la fermer plus facilement).

Un client a aussi des méthodes et un constructeur :
1. **Read** : c'est une go routine qui lira l'input de l'utilsateur physique pour formater un message et l'envoyer à la chatRoom. Nous continuerons ce process infiniment tant que la socket est ouverte.
2. **Write** : go routine pour écrire chez l'utilisateur
3. **Listen** : permettra de mettre le client virtuel sous écoute du client physique (lance simplement les deux go routines read et write).

### Message
```go
type Message struct {
  message_type string
  pseudo string
  content string
}
```
C'est le message qui transite entre le client virtuel et la chatRoom.

### PseudoSocket
```go
type PseudoSocket struct {
	pseudo string
	socket net.Conn
}
```
Ce message est celui envoyé du connectionHandler vers la chatRoom pour demander de rajouter un utilisateur.

## Architecture proposée
De façon synthétisée, voici notre architecture. Le rapport écrit décrit mieux les structures. Le code est compréhensible puisque découpées en méthodes avec des fonctions "raisonnée", de la valeur métier.

![Image : image_generale.png | Erreur, l'image ne se charge pas, allez voir dans le dossier ressource]("./Images/image_generale.png")
