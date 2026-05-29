package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
	sendBufSize    = 256
)

type Client struct {
	hub             *Hub
	conn            *websocket.Conn
	userID          string
	role            ClientRole
	activeRequests  map[string]bool
	send            chan []byte
}

func NewClient(hub *Hub, conn *websocket.Conn, userID string, role ClientRole) *Client {
	return &Client{
		hub:            hub,
		conn:           conn,
		userID:         userID,
		role:           role,
		activeRequests: make(map[string]bool),
		send:           make(chan []byte, sendBufSize),
	}
}

func (c *Client) isSubscribed(requestID string) bool {
	return c.activeRequests[requestID]
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("ws read error: %v", err)
			}
			break
		}

		c.handleMessage(msg)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("ws write error: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(msg []byte) {
	var req struct {
		Event     string `json:"event"`
		RequestID string `json:"request_id,omitempty"`
	}

	if err := json.Unmarshal(msg, &req); err != nil {
		log.Printf("ws unmarshal error: %v", err)
		return
	}

	switch req.Event {
	case EventSubscribe:
		if req.RequestID != "" {
			c.activeRequests[req.RequestID] = true
		}

	case EventUnsubscribe:
		delete(c.activeRequests, req.RequestID)

	case EventVetAcceptReq, EventVetRejectReq:
		c.hub.BroadcastToRequest(req.RequestID, Event{
			Event:     req.Event,
			RequestID: req.RequestID,
		})
	}
}
