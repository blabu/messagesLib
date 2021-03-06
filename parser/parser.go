package parser

import (
	"io"

	"github.com/blabu/messagesLib/dto"
)

// IParser - Основной интерфейс для парсинга сообщений (определяет базовый функционал всего приложения)
// Это связующая точка между сервером и бизнес логикой инициализация которой происходит здесь же
// Сделан для того чтобы была возможность модифицировать протокол передачи данных, например добавить "взрослое" шифрование
// При этом достаточно будет сделать делегата, который реализет этот интерфейс, будет расшифровывать сообщение + делегировать остальной функционал обычному парсеру
type IParser interface {
	FormMessage(msg *dto.Message) ([]byte, error)
	ParseMessage(data []byte) (dto.Message, error)
	IsFullReceiveMsg(data []byte) (int, error)
	ReadPacketHeader(r io.Reader) ([]byte, error)
}
