package initialize

//
//import (
//	"golang.org/x/net/websocket"
//	"log"
//	"net/http"
//	"sync"
//	"time"
//)
//
//type Connection struct {
//	ws             *websocket.Conn
//	userID         string
//	lastPingTime   time.Time
//	pingTimeoutCnt int
//}
//
//type ConnectionPool struct {
//	connections map[string][]*Connection
//	mutex       sync.Mutex
//}
//
//// 定义 ping 消息内容
//var pingMessage = []byte("ping")
//
//// 定义服务器向客户端发送 ping 消息的间隔时间和 ping 超时时间
//const pingInterval = 30 * time.Second
//const pingTimeout = 10 * time.Second
//
//// 将 ping 消息发送给客户端
//func (c *Connection) sendPingMessage() error {
//	err := c.ws.WriteMessage(websocket.PingMessage, pingMessage)
//	if err == nil {
//		c.lastPingTime = time.Now()
//	}
//	return err
//}
//
//// 检查客户端是否已经下线，如果是则关闭连接
//func (c *Connection) checkTimeout() {
//	if time.Since(c.lastPingTime) > pingTimeout {
//		c.pingTimeoutCnt++
//	} else {
//		c.pingTimeoutCnt = 0
//	}
//	if c.pingTimeoutCnt >= 3 {
//		c.ws.Close()
//	}
//}
//
//func (pool *ConnectionPool) handleConnections() {
//	// 定时向所有客户端发送 ping 消息
//	ticker := time.NewTicker(pingInterval)
//	defer ticker.Stop()
//	for range ticker.C {
//		pool.mutex.Lock()
//		for _, connections := range pool.connections {
//			for _, conn := range connections {
//				conn.sendPingMessage()
//				conn.checkTimeout()
//			}
//		}
//		pool.mutex.Unlock()
//	}
//}
//
//func (pool *ConnectionPool) Register(ws *websocket.Conn, userID string) {
//	conn := &Connection{ws, userID, time.Now(), 0}
//	pool.mutex.Lock()
//	defer pool.mutex.Unlock()
//	pool.connections[userID] = append(pool.connections[userID], conn)
//}
//
//func (pool *ConnectionPool) Unregister(ws *websocket.Conn, userID string) {
//	pool.mutex.Lock()
//	defer pool.mutex.Unlock()
//	connections := pool.connections[userID]
//	for i, conn := range connections {
//		if conn.ws == ws {
//			pool.connections[userID] = append(connections[:i], connections[i+1:]...)
//			break
//		}
//	}
//}
//
//func (pool *ConnectionPool) GetConnections(userID string) []*Connection {
//	pool.mutex.Lock()
//	defer pool.mutex.Unlock()
//	return pool.connections[userID]
//}
//
//func main() {
//	// 创建连接池和 HTTP 服务器
//	connectionPool := ConnectionPool{connections: make(map[string][]*Connection)}
//	http.HandleFunc("/ws", connectionPool.WebsocketHandler)
//
//	// 启动发送 ping 消息的定时器
//	go connectionPool.handleConnections()
//
//	// 启动 HTTP 服务
//	err := http.ListenAndServe(":8080", nil)
//	if err != nil {
//		log.Fatal("ListenAndServe: ", err)
//	}
//}
//
//func (connectionPool *ConnectionPool) websocketHandler(w http.ResponseWriter, r *http.Request) {
//	// 将 HTTP 连接 upgrade 为 WebSocket 连接
//	ws, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		return
//	}
//	defer ws.Close()
//
//	// 从 Query 参数中获取用户 ID
//	userID := r.URL.Query().Get("user_id")
//
//	// 将 WebSocket 连接注册到连接池中，带上用户 ID
//	connectionPool.Register(ws, userID)
//
//	// 持续接收 WebSocket 消息，并转发给目标客户端
//	for {
//		_, message, err := ws.ReadMessage()
//		if err != nil {
//			connectionPool.Unregister(ws)
//			break
//		}
//
//		// 从目标客户端的消息中解析出目标用户 ID
//		var messageObj websocket.Message
//		err = json.Unmarshal(message, &messageObj)
//		if err != nil {
//			continue
//		}
//		targetUserID := messageObj.TargetUserID
//
//		// 在连接池中找到目标用户的连接
//		connections := connectionPool.GetConnections(targetUserID)
//		for _, conn := range connections {
//			err = conn.WriteMessage(websocket.TextMessage, message)
//			if err != nil {
//				connectionPool.Unregister(conn)
//			}
//		}
//	}
//}
