package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/blabu/messagesLib/dto"
)

const (
	headerParamSize       = 8
	startSymb        byte = '$'
	versionAttribute byte = 'V'
)
const unsupportedSymb = "~!$%^&*)(+=-}{][,;\\/'\""

// BeginHeader - Початок повідомлення
const BeginHeader = "$V"

// EndHeader - кінець повідомлення
const EndHeader = "###"

var delim = []byte(";")

type header struct {
	protocolVer uint64 // Версия протокола
	command     uint64 // Команда
	mType       string // Тип сообщения (смотри клиента)
	headerSize  int    // Размер заголовка
	contentSize int    // Размер данных
	id          uint8  // id сообщения

	channel string // channel name
	from    string
	to      string
}

// C2cParser - Парсер разбирает сообщения по протоколу
// 1 - клиент-клиент
type C2cParser struct {
	maxPackageSize uint64
	head           header
}

func checksumCustom(arr []byte) uint32 {
	var crc uint32
	for _, val := range arr {
		add := uint16(val)
		if (crc & 0x80_00_00_00) > 0 {
			add++
		}
		crc &= 0x7F_FF_FF_FF
		crc += uint32(add)
		crc <<= 1
		crc++
	}
	return crc
}

//CreateEmptyParser - создает интерфейс парсера с ограничением максимального размера сообщения maxSize
// Кусок принятого сообщения нужен для создания других видов парсера в будущем
func CreateEmptyParser(maxSize uint64) IParser {
	c2c := new(C2cParser)
	c2c.maxPackageSize = maxSize
	return c2c
}

func (c2c *C2cParser) addChecksum(arr []byte) []byte {
	var checksum = make([]byte, 4)
	binary.LittleEndian.PutUint32(checksum, checksumCustom(arr))
	return append(arr, checksum...)
}

//FormMessage - from - Content[0], to - Content[1], data - Content[2]
func (c2c *C2cParser) FormMessage(msg *dto.Message) ([]byte, error) {
	if msg == nil {
		return []byte{}, errors.New("Message nil")
	}
	if msg.Proto == 0 {
		msg.Proto = 1
	}
	res := make([]byte, 0, 128+len(msg.Data))
	res = append(res, []byte(BeginHeader)...)
	res = append(res, []byte(strings.ToUpper(strconv.FormatUint(uint64(msg.Proto), 16)))...)
	res = append(res, ';')
	res = append(res, msg.From...)
	res = append(res, ';')
	res = append(res, msg.To...)
	res = append(res, ';')
	res = append(res, []byte(strings.ToUpper(strconv.FormatUint(uint64(msg.Command), 16)))...)
	res = append(res, ';')
	res = append(res, strings.ToUpper(msg.ContentType)[0]) // convert "text" to T, "binary" to "B" device-to-device protocol specific
	res = append(res, ';')
	res = append(res, []byte(msg.Channel)...) // add name of channel
	res = append(res, ';')
	res = append(res, []byte(strings.ToUpper(strconv.FormatUint(uint64(msg.ID), 10)))...)
	res = append(res, ';')
	res = append(res, []byte(strings.ToUpper(strconv.FormatUint(uint64(len(msg.Data)+4), 16)))...) // plus 4 in message length is add crc calculation
	res = append(res, []byte(EndHeader)...)
	res = append(res, msg.Data...)
	return c2c.addChecksum(res), nil
}

// return position for start header or/and error if not find header or parsing error
func (c2c *C2cParser) parseHeader(data []byte) (int, error) {
	if data == nil {
		return -1, errors.New("Input is empty, nothing to be parsed")
	}
	index := bytes.Index(data, []byte(BeginHeader))
	if index < 0 {
		return index, fmt.Errorf("Package must be started from %s", BeginHeader)
	}
	c2c.head.headerSize = bytes.Index(data, []byte(EndHeader)) // Поиск конца заголовка
	if c2c.head.headerSize < index || c2c.head.headerSize >= len(data) {
		return index, fmt.Errorf("Undefined end header %s in message %s", EndHeader, string(data))
	}
	parsed := bytes.Split(data[index+2:c2c.head.headerSize], delim) // index+2 - пропускаем $V
	if len(parsed) < headerParamSize {
		return index, errors.New("Incorrect header")
	}
	var err error
	if c2c.head.protocolVer, err = strconv.ParseUint(string(parsed[0]), 16, 64); err != nil { //Версия протокола
		return index, errors.New("Icorrect protocol version, it must be a number")
	}
	if c2c.head.protocolVer != 1 {
		return index, errors.New("Icorrect protocol version number. Or parser is invalid")
	}
	c2c.head.from = string(parsed[1])                                                     // от кого
	c2c.head.to = string(parsed[2])                                                       //кому
	if c2c.head.command, err = strconv.ParseUint(string(parsed[3]), 16, 64); err != nil { //команда
		return index, errors.New("Error usuported porotocol")
	}
	c2c.head.mType = string(parsed[4]) //тип сообщения
	switch c2c.head.mType {
	case "T":
		c2c.head.mType = "text"
	case "B":
		c2c.head.mType = "binary"
	case "A":
		c2c.head.mType = "audio"
	case "V":
		c2c.head.mType = "video"
	case "F":
		c2c.head.mType = "file"
	}
	c2c.head.channel = string(parsed[5])
	var s uint64
	if s, err = strconv.ParseUint(string(parsed[6]), 16, 8); err != nil { //id сообщения
		c2c.head.id = 0
	} else {
		c2c.head.id = uint8(s)
	}
	if s, err = strconv.ParseUint(string(parsed[7]), 16, 64); err != nil { //размер сообщения
		return index, errors.New("Icorrect message size, it must be a number")
	}
	if s > c2c.maxPackageSize {
		return index, fmt.Errorf("Income package is too big parsed %s to %d. Overflow internal buffer %d", string(parsed[4]), s, c2c.maxPackageSize)
	}
	c2c.head.contentSize = int(s)
	c2c.head.headerSize += len(EndHeader) // Add endHeader
	return index, nil
}

//ParseMessage - from - Content[0], to - Content[1], data - Content[2]
func (c2c *C2cParser) ParseMessage(data []byte) (dto.Message, error) {
	var err error
	var i int
	if i, err = c2c.parseHeader(data); err != nil {
		return dto.Message{}, err
	}
	if len(data) < i+c2c.head.headerSize+c2c.head.contentSize {
		return dto.Message{}, errors.New("Not full message")
	}
	defer func() {
		c2c.head = header{}
	}()
	content := make([]byte, c2c.head.contentSize-4) // Delete crc32 sum from end of package
	copy(content, data[i+c2c.head.headerSize:i+c2c.head.headerSize+c2c.head.contentSize-4])
	crc := checksumCustom(data[i : i+c2c.head.headerSize+c2c.head.contentSize-4])
	if crc != binary.LittleEndian.Uint32(data[i+c2c.head.headerSize+c2c.head.contentSize-4:]) {
		return dto.Message{}, errors.New("Invalid checksum")
	}
	var result dto.Message
	result.MessageMetaInf = dto.MessageMetaInf{
		Command: uint16(c2c.head.command),
		Proto:   uint16(c2c.head.protocolVer),
		ID:      c2c.head.id,
		Channel: c2c.head.channel,
		From:    c2c.head.from,
		To:      c2c.head.to,
	}
	result.MessageContent = dto.MessageContent{
		ContentType: c2c.head.mType,
		Data:        content,
	}
	return result, nil
}

// IsFullReceiveMsg - Проверка пришел полный пакет или нет
// TODO каждый раз парсить заголовок не эффективно надо будет переписать
func (c2c *C2cParser) IsFullReceiveMsg(data []byte) (int, error) {
	if _, err := c2c.parseHeader(data); err != nil {
		return -1, err
	}
	lastBytes := c2c.head.contentSize + c2c.head.headerSize - len(data)
	if lastBytes < 0 {
		return 0, nil
	}
	return lastBytes, nil
}

//ReadPacketHeader - Читает заголовок и возвращает полученный результат
func (c2c *C2cParser) ReadPacketHeader(r io.Reader) ([]byte, error) {
	buf := make([]byte, len(BeginHeader)+headerParamSize*(len(delim)+1)+len(EndHeader))
	header := make([]byte, 0, len(buf))
	for {
		if n, err := r.Read(buf); err == nil {
			header = append(header, buf[:n]...)
		} else {
			return nil, err
		}
		if bytes.Index(header, []byte(EndHeader)) >= 0 {
			break
		}
	}
	return header, nil
}
