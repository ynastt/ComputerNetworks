/*lab1*/
package proto

import "encoding/json"

// Request -- запрос клиента к серверу.
type Request struct {
	// Поле Command может принимать три значения:
	// * "quit" - прощание с сервером (после этого сервер рвёт соединение);
	// * "add" - передача новой пары на сервер;
	// * "avg" - просьба посчитать нод и нок.
	Command string `json:"command"`

	// Если Command == "add", в поле Data должна лежать пара
	// в виде структуры Pair.
	// В противном случае, поле Data пустое.
	Data *json.RawMessage `json:"data"`
}

// Response -- ответ сервера клиенту.
type Response struct {
	// Поле Status может принимать три значения:
	// * "ok" - успешное выполнение команды "quit" или "add";
	// * "failed" - в процессе выполнения команды произошла ошибка;
	// * "result" - нок и нод вычислены.
	Status string `json:"status"`

	// Если Status == "failed", то в поле Data находится сообщение об ошибке.
	// Если Status == "result", в поле Data должна лежать пара чисел - нок и нод
	// в виде структуры Pair.
	// В противном случае, поле Data пустое.
	Data *json.RawMessage `json:"data"`
}

// Pair -- пара чисел.
type Pair struct {
	// 1 число 
	Num1 string `json:"numo"`

	// 2 число 
	Num2 string `json:"numt"`
}

