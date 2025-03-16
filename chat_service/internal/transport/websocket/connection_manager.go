package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
)

type ConnectionManager struct {
	connection map[int]*websocket.Conn
	mu         sync.Mutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connection: make(map[int]*websocket.Conn),
	}
}

func (c *ConnectionManager) AddConnection(userID int, conn *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connection[userID] = conn
}

func (c *ConnectionManager) RemoveConnection(userID int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.connection, userID)
}

func (c *ConnectionManager) GetConnection(userID int) (*websocket.Conn, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	conn, exists := c.connection[userID]
	return conn, exists
}
