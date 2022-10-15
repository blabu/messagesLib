package dto

import "time"

//MessageContent - Содержимое сообщения безотносительно от того кому и куда оно отправлено.
//Для исключения дубликатов ключ - это хеш сумма содержимого,
//таким образом, сообщение не дублируется если его пересылать или повторно отправить
type MessageContent struct {
	Hash        string    `json:"hash" db:"Hash"`
	ContentType string    `json:"contentType" db:"ContentType"`
	Data        []byte    `json:"data,omitempty" db:"Data"`
	CreatedDate time.Time `json:"createdDate" db:"CreatedDate"`
	ModifDate   time.Time `json:"modifDate,omitempty" db:"ModifDate"`
}

//MessageMetaInf - Хранит мета информацию от кого куда во сколько.
//История сообщений между пользователями.
type MessageMetaInf struct {
	ID          uint8  `json:"-" db:"-"`
	UID         string `json:"uid" db:"UID"`
	ContentHash string `json:"contentHash" db:"ContentHash"`
	Proto       uint16 `json:"proto,omitempty" db:"Proto"`
	Command     uint16 `json:"cmd,omitempty" db:"Command"`
	Channel     string `json:"channel,omitempty" db:"Channel"`
	Name        string `json:"name,omitempty" db:"Name"`
	AddedTime   int64  `json:"addedTime" db:"AddedTime"`
	SendedTime  int64  `json:"sendedTime,omitempty" db:"SendedTime"`
	From        string `json:"from,omitempty" db:"FromName"`
	To          string `json:"to,omitempty" db:"ToName"`
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
	Token       string    `json:"token,omitempty" db:"Token"`
	ImageURL    string    `json:"image,omitempty" db:"ImageURL"`
	CreatedDate time.Time `json:"created,omitempty" db:"CreatedDate"`
	LastDate    time.Time `json:"activity,omitempty" db:"LastDate"`
}

//Bot is entity with base communication and validation functions
type Bot struct {
	ClientDescriptor
	CreatedBy   string `json:"who,omitempty" db:"CreatedBy"`
	About       string `json:"about,omitempty" db:"About"`
	Endpoint    string `json:"endpoint,omitempty" db:"Endpoint"`  // POST request to send data to bot
	HealthCheck string `json:"health,omitempty" db:"HealthCheck"` // GET request to check health bot with return http.StatusOK if all ok. Return any other close connection destroy bot entity
}

// Channel is entity that store information about channel
type Channel struct {
	ClientDescriptor
	More      string `json:"more,omitempty" db:"More"`
	About     string `json:"about,omitempty" db:"About"`
	CreatedBy string `json:"who,omitempty" db:"CreatedBy"`
}

//ModemState - текущее состояние модема, передается при пинге устройства
type ModemState struct {
	Name     string `json:"name"`
	LastPing int64  `json:"lastPing,omitempty"`
	Voltage  uint16 `json:"voltage,omitempty"`
	Signal   uint8  `json:"signal,omitempty"`
}
