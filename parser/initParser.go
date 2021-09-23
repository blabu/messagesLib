package parser

func InitParser(rec []byte, size uint64) (IParser, error) {
	return CreateEmptyParser(size), nil
}
