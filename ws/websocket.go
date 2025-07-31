package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"Taskly.com/m/global"
	model "Taskly.com/m/internal/models"
	"Taskly.com/m/internal/service"
	"Taskly.com/m/package/utils/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client struct được đơn giản hóa, không cần trường Auth
type Client struct {
	UserInfo *model.User
	Conn     *websocket.Conn
}
type RoomMessage struct {
	RoomId  uuid.UUID
	Message []byte
}

type ConnectionManager struct {
	clients    map[*websocket.Conn]*Client
	roomUsers  map[uuid.UUID][]*Client
	broadcast  chan RoomMessage
	register   chan *Client
	unregister chan *websocket.Conn
	mu         sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		clients:    make(map[*websocket.Conn]*Client),
		roomUsers:  make(map[uuid.UUID][]*Client),
		broadcast:  make(chan RoomMessage),
		register:   make(chan *Client),
		unregister: make(chan *websocket.Conn),
	}
}
func (cm *ConnectionManager) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				cm.mu.Lock()
				for conn := range cm.clients {
					conn.Close()
				}
				cm.mu.Unlock()
				return
			case client := <-cm.register:
				cm.clients[client.Conn] = client

			case conn := <-cm.unregister:
				cm.RemoveClient(conn)

			case roomMessage := <-cm.broadcast:
				cm.BroadcastToRoom(roomMessage.RoomId, roomMessage.Message)
			}
		}
	}()
}
func (cm *ConnectionManager) AddClient(client *Client, roomID uuid.UUID) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, existingClient := range cm.roomUsers[roomID] {
		if existingClient.Conn == client.Conn {
			return
		}
	}

	cm.roomUsers[roomID] = append(cm.roomUsers[roomID], client)
}

func (cm *ConnectionManager) RemoveClient(conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if client, exists := cm.clients[conn]; exists {
		delete(cm.clients, conn)
		for roomId, users := range cm.roomUsers {
			for i, c := range users {
				if c.Conn == conn {
					cm.roomUsers[roomId] = append(users[:i], users[i+1:]...)
					if len(cm.roomUsers[roomId]) == 0 {
						delete(cm.roomUsers, roomId)
					}
					break
				}
			}
		}
		client.Conn.Close()
	}
}

func (cm *ConnectionManager) BroadcastToRoom(roomID uuid.UUID, message []byte) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if clients, exists := cm.roomUsers[roomID]; exists {
		for _, client := range clients {
			err := client.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				global.Logger.Sugar().Errorf("Error sending message: %v", err)
				client.Conn.Close()
				cm.unregister <- client.Conn
			}
		}
	} else {
		global.Logger.Sugar().Infof("No clients in room %s to send message", roomID)
	}
}

func HandleConnections(ctx *gin.Context, cm *ConnectionManager) {
	userInfo := auth.GetUserFromContext(ctx)
	if userInfo == nil || userInfo.ID == uuid.Nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		ctx.Abort()
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		global.Logger.Sugar().Errorf("Lỗi khi nâng cấp WebSocket: %v", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			global.Logger.Sugar().Errorf("Panic in websocket handler: %v", r)
		}
		cm.unregister <- conn
	}()

	client := &Client{
		UserInfo: userInfo,
		Conn:     conn,
	}
	cm.register <- client

	const (
		pongWait   = 60 * time.Second
		pingPeriod = (pongWait * 9) / 10
	)

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				global.Logger.Sugar().Warnf("Unexpected WebSocket closure: %v", err)
			}
			cm.unregister <- conn
			return
		}
		cm.handleMessage(message, client, ctx)
	}
}

func (cm *ConnectionManager) handleMessage(message []byte, client *Client, ctx *gin.Context) {
	var msgData map[string]interface{}
	if err := json.Unmarshal(message, &msgData); err != nil {
		client.Conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
		return
	}
	action, ok := msgData["action"].(string)
	if !ok {
		client.Conn.WriteMessage(websocket.TextMessage, []byte("Missing action"))
		return
	}

	switch action {
	case "join":
		roomIDStr, ok := msgData["room_id"].(string)
		if !ok {
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Invalid room_id"))
			return
		}
		roomID, err := uuid.Parse(roomIDStr)
		if err != nil {
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Invalid room_id format"))
			return
		}
		cm.AddClient(client, roomID)

	case "leave":
		roomIDStr, ok := msgData["room_id"].(string)
		if !ok {
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Invalid room_id"))
			return
		}
		_, err := uuid.Parse(roomIDStr)
		if err != nil {
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Invalid room_id format"))
			return
		}
		cm.RemoveClient(client.Conn)

	case "send_message":
		roomIDStr, ok := msgData["room_id"].(string)
		if !ok {
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Invalid room_id"))
			return
		}
		roomID, err := uuid.Parse(roomIDStr)
		if err != nil {
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Invalid room_id format"))
			return
		}

		receiverIDStr, ok := msgData["receiver_id"].(string)
		if !ok {
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Invalid receiver_id"))
			return
		}
		receiverID, err := uuid.Parse(receiverIDStr)
		if err != nil {
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Invalid receiver_id format"))
			return
		}

		content, ok := msgData["content"].(string)
		if !ok {
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Invalid content"))
			return
		}

		chatInput := model.SetChatInput{
			RoomID:     roomID,
			SenderID:   client.UserInfo.ID,
			ReceiverID: receiverID,
			Content:    content,
		}

		savedMessage, err := service.GetChatService().SetChatHistory(ctx, chatInput)
		if err != nil {
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Error saving message"))
			return
		}

		jsonMessage, err := json.Marshal(savedMessage)
		if err != nil {
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Error formatting message"))
			return
		}

		cm.broadcast <- RoomMessage{
			RoomId:  roomID,
			Message: jsonMessage,
		}

	default:
		client.Conn.WriteMessage(websocket.TextMessage, []byte("Invalid action"))
	}
}

func (cm *ConnectionManager) cleanupEmptyRooms() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			cm.mu.Lock()
			for roomID, clients := range cm.roomUsers {
				if len(clients) == 0 {
					delete(cm.roomUsers, roomID)
				}
			}
			cm.mu.Unlock()
		}
	}()
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketMessage struct {
	Action     string `json:"action"`
	RoomID     string `json:"room_id,omitempty"`
	Content    string `json:"content,omitempty"`
	ReceiverID string `json:"receiver_id,omitempty"`
}
