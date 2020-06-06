package tenet

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/maciej-zuk/tecligo/common"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
	sendQueueSize  = 256
)

type websocketClient struct {
	hub  *websocketHub
	conn *websocket.Conn
	ip   string
	send chan interface{}
}

type websocketHub struct {
	clients       map[*websocketClient]bool
	broadcast     chan interface{}
	register      chan *websocketClient
	unregister    chan *websocketClient
	botConnection *Connection
}

type websocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (h *websocketHub) close() {
	for c := range h.clients {
		h.unregister <- c
		c.conn.Close()
	}
}

func (h *websocketHub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (c *websocketClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.conn.ReadMessage()
		if err == nil {
			var msgStruct websocketMessage
			err = json.Unmarshal(msg, &msgStruct)
			if err == nil {
				if msgStruct.Type == "chat" {
					text, ok := msgStruct.Data.(string)
					if ok {
						c.hub.botConnection.Send(
							82,
							TnShort(1),
							TnString("Say"),
							TnString(fmt.Sprintf("[c/FF00FF:%s]: [c/FFFF00:%s]", c.ip, text)),
						)
					}
				}
			}
		} else {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WS read error: %v", err)
			}
			break
		}
	}
}

func (c *websocketClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(message); err != nil {
				log.Println("WS write error:", err)
				return
			}
			n := len(c.send)
			for i := 0; i < n; i++ {
				if err := c.conn.WriteJSON(<-c.send); err != nil {
					log.Println("WS write error:", err)
					return
				}
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func serveWs(hub *websocketHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS upgrade error:", err)
		return
	}
	client := &websocketClient{
		hub:  hub,
		conn: conn,
		send: make(chan interface{}, sendQueueSize),
		ip:   r.Header.Get("x-forwarded-for"),
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()

	time.Sleep(1 * time.Second)

	for slot, p := range hub.botConnection.Players {
		if p.Active {
			client.send <- hub.botConnection.getPlayerUpdateMessage(TnByte(slot))
		}
	}
	client.send <- hub.botConnection.getWorldUpdateMessage()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var connections map[*websocket.Conn]bool

func websocketRoutine(c *Connection) {
	c.hub = &websocketHub{
		broadcast:     make(chan interface{}),
		register:      make(chan *websocketClient),
		unregister:    make(chan *websocketClient),
		clients:       make(map[*websocketClient]bool),
		botConnection: c,
	}
	go c.hub.run()
	mux := http.NewServeMux()
	mux.HandleFunc(common.Settings.SocketEndpoint, func(w http.ResponseWriter, r *http.Request) {
		serveWs(c.hub, w, r)
	})
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", common.Settings.SocketPort),
		Handler: mux,
	}
	go srv.ListenAndServe()
	// wait for main to finish
	c.enterWg.Wait()
	c.hub.close()
	c.hub.botConnection = nil
	c.hub = nil
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	cancel()
	<-ctx.Done()
	c.exitWg.Done()
}
