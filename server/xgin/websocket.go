package xgin

import (
	"io"
	"net"
	"net/http"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/gorilla/websocket"
)

// WebSocketConn 抽象一条已建立的 WebSocket 连接，覆盖子协议查询、读写消息、控制帧发送、
// 超时设置、关闭/Ping/Pong 处理器注册以及底层连接访问等操作，方法集合与 websocket.Conn 保持一致。
type WebSocketConn interface {
	Subprotocol() string
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	WriteControl(messageType int, data []byte, deadline time.Time) error
	NextWriter(messageType int) (io.WriteCloser, error)
	WritePreparedMessage(pm *websocket.PreparedMessage) error
	WriteMessage(messageType int, data []byte) error
	SetWriteDeadline(t time.Time) error
	NextReader() (messageType int, r io.Reader, err error)
	ReadMessage() (messageType int, p []byte, err error)
	SetReadDeadline(t time.Time) error
	SetReadLimit(limit int64)
	CloseHandler() func(code int, text string) error
	SetCloseHandler(h func(code int, text string) error)
	PingHandler() func(appData string) error
	SetPingHandler(h func(appData string) error)
	PongHandler() func(appData string) error
	SetPongHandler(h func(appData string) error)
	UnderlyingConn() net.Conn
	EnableWriteCompression(enable bool)
	SetCompressionLevel(level int) error
}

// WebSocketFunc 为 WebSocket 回调函数；conn 为升级成功后的连接，err 非空表示协议升级失败。
type WebSocketFunc func(WebSocketConn, error)

// WebSocket 描述一条 WebSocket 路由的模式、升级器、请求头与业务处理回调。
//
// 零值不可直接使用，须通过 WebSocketOptions 结合 WebSocketOption 构造。
type WebSocket struct {
	// Pattern 为路由匹配路径。
	Pattern string
	// Handler 为连接建立或失败后的回调。
	Handler WebSocketFunc
	*websocket.Upgrader
	// Header 为升级握手时附加的响应请求头。
	Header http.Header
}

// Upgrade 执行 HTTP 到 WebSocket 的协议升级；升级成功后调用 Handler 处理连接，结束时关闭连接。
func (ws *WebSocket) Upgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrader.Upgrade(w, r, ws.Header)
	if err == nil {
		defer func() {
			if err := conn.Close(); err != nil {
				logger.Logger.Error("Error closing conn: ", err)
			}
		}()
	}
	ws.Handler(conn, err)
}

// WebSocketOption 是 WebSocket 的配置选项函数，在 WebSocketOptions 构造时修改 WebSocket。
type WebSocketOption func(*WebSocket)

// WebSocketOptions 按路由模式与处理器构造 WebSocket，并依次套用传入的配置选项。
func WebSocketOptions(pattern string, handler WebSocketFunc, opts ...WebSocketOption) *WebSocket {
	ws := &WebSocket{
		Pattern:  pattern,
		Handler:  handler,
		Upgrader: &websocket.Upgrader{},
	}
	for _, opt := range opts {
		opt(ws)
	}
	return ws
}
