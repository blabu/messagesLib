package dto

import (
	"context"
	"io"
	"time"
)

//IClient - интерфейс чтения записи сообщений в систему и из неё
type IClient interface {
	GetName() string
	Read(ctx context.Context, m *Message) error
	Write(ctx context.Context, m *Message) error
	ReadNextNotSended(ctx context.Context, m *Message) error
	IsOnline(ctx context.Context, name string) (bool, error)
	Close()
}

//InitializerIClient - функция инициализатор работы с мессенджером
type InitializerIClient func(uid, name string) IClient

//IContacts - базовый интерфейс по работе со списком контактов абонента
type IContacts interface {
	AddContact(ctx context.Context, from, to string) error             // AddContact - Добавить контакт к списку контактов "From" могут быть повторы (следует исключить возможные повторы в реализации интерфейса)
	GetAllContacts(ctx context.Context, from string) ([]string, error) // GetAllContacts - Получить список всех контактов по имени "From"
}

//IMessageHistoryReader - Интерфейс читатель мета информации сообщения (куда, кому, во сколько)
type IMessageHistoryReader interface {
	GetReceivedMessages(ctx context.Context, self, from string, until time.Time, limit int64) ([]MessageMetaInf, error) // Получить отправленные сообщения из списка "Полученные от"
	GetSendedMessages(ctx context.Context, self, to string, until time.Time, limit int64) ([]MessageMetaInf, error)     // Получить отправленные сообщения из списка "Отправленные кем"
	GetByUID(ctx context.Context, uid string) (MessageMetaInf, error)
}

//IMessageHistoryWriter - Интерфейс писатель добавляет, редактирует (в случае совпадения ключа) и удаляет информацию о сообщении
type IMessageHistoryWriter interface {
	AddTo(ctx context.Context, from string, msg *MessageMetaInf) error // Добавляем сообщение в "Полученые от" список
	AddFrom(ctx context.Context, to string, msg *MessageMetaInf) error // Добавляем сообщение в "Отправленные кем" список
	Delete(ctx context.Context, from, to, id string) error             // Удалить сообщения из списков от кого и кому
}

type IMessageHistory interface {
	IMessageHistoryReader
	IMessageHistoryWriter
}

type IMessageSaver interface {
	SaveMessage(ctx context.Context, msg *MessageContent) error          // Сохранение содержимого сообщения по ключу в хранилище
	GetMessage(ctx context.Context, hash string) (MessageContent, error) // Получить контент сообщения по его ключу (ключ - это комбинация типа контента и хеша содержимого)
}

/*
IMessanger - интерфейс сохранения и обмена сообщениями и контактами
Сообщение хранится отдельно от истории его отправки.
Само сообщение зависит только от контента и не зависит от того кому, куда и во сколько было отправлено.
Поэтому для исключения повторов одних и тех же сообщений сообщение хранится в базе по уникальному ключу (Хеш сумме)
А уже этот ключ сохраняется в списке отправленных и принятых сообщений у каждого клиента.
За это и отвечает интерфейс
*/
type IMessanger interface {
	AddMessageFrom(ctx context.Context, to string, m *Message) error
	AddMessageTo(ctx context.Context, from string, m *Message) error
	GetMessage(ctx context.Context, uid string) (Message, error)
	GetReceivedMessages(ctx context.Context, self, from string, until time.Time, limit int64) ([]Message, error)
	GetSendedMessages(ctx context.Context, self, to string, until time.Time, limit int64) ([]Message, error)
	Delete(ctx context.Context, from, to, uid string) error
}

//ReadWriteCloser - создает интерфейс работы с модемом через tcp или tls соединение
type ReadWriteCloser interface {
	// Write - Передаем данные полученные из сети бизнес логике
	Write(ctx context.Context, msg *Message) error

	//Read - читаем ответ бизнес логики return io.EOF if client never answer
	Read(ctx context.Context, msg *Message) error

	// Close - информирует бизнес логику про разрыв соединения
	io.Closer
}

// Salt - is a random string that must be a uniq in system for all time for one client name in descriptor
type Salt string

//IBgSalt - update salt and return a new value. If value > max integer return max integer
type IBgSalt interface {
	Check(ctx context.Context, name string, s Salt) uint64
}

//IBgClientSaver - Интерфейс работы с модемами
type IBgClientSaver interface {
	GetClient(ctx context.Context, name string) (ClientDescriptor, error)
	SaveClient(ctx context.Context, cl *ClientDescriptor) error
	GenerateClient(ctx context.Context, name string) (ClientDescriptor, error)
}

type IBgMsgSaver interface {
	Get(ctx context.Context, key string) (MessageContent, error)
	Set(ctx context.Context, key string, val *MessageContent) error
}

type IBgMetaSaver interface {
	Get(ctx context.Context, key string) (MessageMetaInf, error)
	Set(ctx context.Context, key string, val *MessageMetaInf) error
	Del(ctx context.Context, key string) error
}

type IBgBotSaver interface {
	Get(ctx context.Context, key string) (Bot, error)
	Set(ctx context.Context, key string, val *Bot) error
}

type IBgHashHistory interface {
	AddHash(ctx context.Context, from, to, hash string, timestamp int64) error
	GetHash(ctx context.Context, from, to string, maxTimestamp, limit int64) ([]string, error)
	RemoveHash(ctx context.Context, from, to, hash string) error
}

type IBgChannels interface {
	Save(ctx context.Context, channel *Channel) error
	Get(ctx context.Context, name string) (Channel, error)
}

type IChannel interface {
	AddMessage(ctx context.Context, msg *Message) error
	GetMessages(ctx context.Context, until time.Time, limit int64) ([]Message, error)
}
