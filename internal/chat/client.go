package chat

import (
	"context"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID  string
	Conn    *websocket.Conn
	Service Service
	Hub     *Hub
	Ctx     context.Context
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.RemoveClient(c.UserID)
		c.Conn.Close()
	}()

	for {
		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			break
		}

		// Save message to DB
		msg.Delivered = true
		_ = c.Service.SaveMessage(c.Ctx, &msg)
		// Send to receiver if online
		if c.Hub.IsOnline(msg.ReceiverID) {
			c.Hub.SendToUser(msg.ReceiverID, &msg)
			c.Service.MarkRead(c.Ctx, msg.ReceiverID, msg.SenderID)

		}
	}
}

func (c *Client) Send(msg *Message) {
	c.Conn.WriteJSON(msg)
}
