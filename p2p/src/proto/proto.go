package proto

import "encoding/json"

// Response -- ответ от соседнего пира.
type Response struct {
	// Поле Status может принимать значение:
	// * "ls" - просьба вывести список друзей.
	Status string `json:"status"`

	// Если Status == "ls", в поле Data должен лежать элемент
	// в виде структуры Peer.
	// В противном случае, поле Data пустое.
	Data *json.RawMessage `json:"data"`
}

// Request -- сообщение от соседнего пира
type Request struct {
	// Поле Command может принимать три значения:
	// * "friend" - зафрендить другого пира по имени;
	// * "unfriend" - отфрендить другого пира;
	// * "quit" - завершить сввязь с данным пиром;
	// * "list" -  распечатать список френдов.
	Command string `json:"command"`
	// В поле Data должен лежать элемент
	// в виде структуры Elem.
	Data *json.RawMessage `json:"data"`
}

// Elem -- структура элемента.
type Elem struct {
	// IP-адрес пира
	IP string `json:"ip"`

	// имя пира тк мы френдим по имени
	Name string `json:"name"`
}
