package bytes

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"sync"
)

var (
	spWriter sync.Pool
	spReader sync.Pool
	spBuffer sync.Pool
)

func init() {
	spWriter = sync.Pool{New: func() interface{} {
		return gzip.NewWriter(nil)
	}}
	spReader = sync.Pool{New: func() interface{} {
		return new(gzip.Reader)
	}}
	spBuffer = sync.Pool{New: func() interface{} {
		return bytes.NewBuffer(nil)
	}}
}

// Unzip возвращает не сжатые данные в виде среза байт.
// Использует пакет gzip.
func Unzip(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)

	gr := spReader.Get().(*gzip.Reader)
	defer func() {
		spReader.Put(gr)
	}()

	err := gr.Reset(buf)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	data, err = ioutil.ReadAll(gr)
	if err != nil {
		return nil, err
	}
	return data, err
}

// Zip возвращает сжаные данные в виде среза байт.
// Использует пакет gzip.
func Zip(data []byte) ([]byte, error) {
	buf := spBuffer.Get().(*bytes.Buffer)
	w := spWriter.Get().(*gzip.Writer)
	w.Reset(buf)

	defer func() {
		buf.Reset()
		spBuffer.Put(buf)
		w.Close()
		spWriter.Put(w)
	}()

	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	dec := buf.Bytes()
	return dec, nil
}
