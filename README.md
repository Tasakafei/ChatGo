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
`

Nos utilisateurs sont regroupés dans des chatRooms. Ainsi à l'avenir nous pourrions rajouter des propriétés propres aux chatRooms telles qu'un nombre d'utilisateurs max ou des permissions spéciales (exemple, tout ceux qui s'appellent GéantVert pourraient avoir accès à des commandes spéciales).


## Architecture proposée
Pour vous décrire l'architecture mise en place dans les moindres détails j'ai préféré la dérouler sous forme de scénario étant donné qu'il est compliqué de faire apparaître comme demandé à la fois les goroutines, les canaux de communication entre les goroutines, à quel moment un nouveau module ou go routine est instanciée et par qui etc. 

Du coup nous avons une première étape qui est la connection. La figure 1 détaille cette étape. Au préalable nous avons un serveur qui instancie une `ChatRoom` et lui passe en paramètre un canal de communication `ServerToChatRoomChannel`. Ce serveur écoute sur le port 1234. Un client arrive et se connecte via netcat (par exemple) sur le serveur. Le serveur va ensuite instancié une nouvelle go routine est dédiée à la gestion des utilisateurs : `handleConnection`. Nous lui passons aussi en paramètres la socket de connection (pour qu'il puisse communiquer avec le client et un canal de communication `ServerToChatRoomChannel`. Ce canal de communication était utilisé plus tôt 
