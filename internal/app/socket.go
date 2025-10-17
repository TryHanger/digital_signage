package app

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
)

func InitSocketServer() *socketio.Server {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("🔌 Новый монитор подключён:", s.ID())
		return nil
	})

	server.OnEvent("/", "register", func(s socketio.Conn, monitorID string) {
		fmt.Println("📺 Зарегистрировался монитор:", monitorID)
		s.Join("monitor_" + monitorID)
	})

	server.OnEvent("/", "status", func(s socketio.Conn, msg string) {
		fmt.Println("📡 Статус:", msg)
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("❌ Ошибка сокета:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("⚡ Монитор отключился:", s.ID(), reason)
	})

	return server
}
