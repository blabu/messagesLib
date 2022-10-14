package dto

import "time"

//MessageContent - Содержимое сообщения безотносительно от того кому и куда оно отправлено.
//Для исключения дубликатов ключ - это хеш сумма содержимого,
//таким образом, сообщение не дублируется если его пересылать или повторно отправить
type MessageContent struct {
	Hash        string    `json:"hash" db:"Hash"`
	ContentType string    `json:"contentType" db:"ContentType"`
	Data        []byte    `json:"data" db:"Data"`
	CreatedDate time.Time `json:"createdDate" db:"CreatedDate"`
	ModifDate   time.Time `json:"modifDate" db:"ModifDate"`
}

//MessageMetaInf - Хранит мета информацию от кого куда во сколько.
//История сообщений между пользователями.
type MessageMetaInf struct {
	ID          uint8  `json:"-" db:"-"`
	UID         string `json:"uid" db:"UID"`
	ContentHash string `json:"contentHash" db:"ContentHash"`
	Proto       uint16 `json:"proto" db:"Proto"`
	Command     uint16 `json:"cmd" db:"Command"`
	Channel     string `json:"channel" db:"Channel"`
	Name        string `json:"name" db:"Name"`
	AddedTime   int64  `json:"addedTime" db:"AddedTime"`
	SendedTime  int64  `json:"sendedTime" db:"SendedTime"`
	From        string `json:"from" db:"FromName"`
	To          string `json:"to" db:"ToName"`
}

//Message - сообщение между клиентами. Сообщение разделено на мета информации и содержимое сообщения разделение позволяет исключить дубликаты содержимого сообщений
//(особенно актуальна на больших сообщениях)
type Message struct {
	MessageMetaInf
	MessageContent
}

//ClientDescriptor - предоставляет базовую информацию о пользователе (человек, участник общения) через веб интерфейс
type ClientDescriptor struct {
	Name        string    `json:"name" db:"Name"`
	Token       string    `json:"token" db:"Token"`
	ImageURL    string    `json:"image" db:"ImageURL"`
	CreatedDate time.Time `json:"created" db:"CreatedDate"`
	LastDate    time.Time `json:"activity" db:"LastDate"`
}

//Bot is entity with base communication and validation functions
type Bot struct {
	ClientDescriptor
	CreatedBy   string `json:"who" db:"CreatedBy"`
	About       string `json:"about" db:"About"`
	Endpoint    string `json:"endpoint" db:"Endpoint"`  // POST request to send data to bot
	HealthCheck string `json:"health" db:"HealthCheck"` // GET request to check health bot with return http.StatusOK if all ok. Return any other close connection destroy bot entity
}

// Channel is entity tha store information about channel
type Channel struct {
	ClientDescriptor
	More      string `json:"more" db:"More"`
	About     string `json:"about" db:"About"`
	CreatedBy string `json:"who" db:"CreatedBy"`
}

//ModemState - текущее состояние модема, передается при пинге устройства
type ModemState struct {
	Name     string `json:"name"`
	LastPing int64  `json:"lastPing"`
	Voltage  uint16 `json:"voltage"`
	Signal   uint8  `json:"signal"`
}
