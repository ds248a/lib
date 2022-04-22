package uuid

import (
	"github.com/google/uuid"

	"github.com/ds248a/lib/strconv"
)

// V4 возвращает строковое значение идентификатора.
func V4() string {
	return uuid.New().String()
}

// V4 возвращает значение идентификатора в виде среза байт.
func V4Bytes() []byte {
	return strconv.S2B(uuid.New().String())
}
