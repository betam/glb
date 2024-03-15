package lib2

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

type WSServer struct {
	clients       map[*websocket.Conn]bool
	handleMessage func(message []byte) // хандлер новых сообщений
	mu            sync.Mutex
}

func StartWSServer(handleMessage func(message []byte), addr string) *WSServer {
	server := WSServer{
		make(map[*websocket.Conn]bool),
		handleMessage,
		sync.Mutex{},
	}

	http.HandleFunc("/", server.echo)
	go http.ListenAndServe(addr, nil) // Уводим http сервер в горутину

	return &server
}

func (server *WSServer) echo(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer connection.Close()

	server.clients[connection] = true        // Сохраняем соединение, используя его как ключ
	defer delete(server.clients, connection) // Удаляем соединение

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
		}

		go server.handleMessage(message)
	}
}

func (server *WSServer) WriteMessage(message []byte) {
	for conn := range server.clients {
		server.mu.Lock()
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Debugf("Can not process msg to socket: %v", err)
			conn.Close()
			delete(server.clients, conn)
		}
		server.mu.Unlock()
	}
}
