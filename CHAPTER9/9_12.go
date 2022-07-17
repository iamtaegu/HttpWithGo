
type Hub struct {
	// 등록된 클라이언트
	clients map[*Client]bool
	// 클라이언트가 보내는 메시지
	broadcast chan []byte
	// 클라이언트의 등록 요청
	register chan *Client
	// 클라이언트의 말소 요청
	unregister chan *Client
}

// 초기화 후에 goroutine으로 실행되는 메서드
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client: = <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range.h.clients {
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