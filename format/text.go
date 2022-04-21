package format

import (
	"io"
)

// Indent вставляет prefix в начале каждой непустой строки s.
func Indent(s, prefix string) string {
	return string(IndentBytes([]byte(s), []byte(prefix)))
}

// IndentBytes вставляет prefix в начале каждой непустого среза b.
func IndentBytes(b, prefix []byte) []byte {
	var res []byte
	bol := true
	for _, c := range b {
		if bol && c != '\n' {
			res = append(res, prefix...)
		}
		res = append(res, c)
		bol = c == '\n'
	}
	return res
}

// Добавляет отступ pre для каждой строки.
type indentWriter struct {
	w   io.Writer
	bol bool
	pre [][]byte
	sel int
	off int
}

// NewIndentWriter задает для io.Writer префикс pre для каждой новой строки.
func NewIndentWriter(w io.Writer, pre ...[]byte) io.Writer {
	return &indentWriter{
		w:   w,
		pre: pre,
		bol: true,
	}
}

// Возвращаются ошибки из заданного потока io.Writer.
func (w *indentWriter) Write(p []byte) (n int, err error) {
	for _, c := range p {
		if w.bol {
			var i int
			i, err = w.w.Write(w.pre[w.sel][w.off:])
			w.off += i
			if err != nil {
				return n, err
			}
		}

		_, err = w.w.Write([]byte{c})
		if err != nil {
			return n, err
		}

		n++
		w.bol = c == '\n'

		if w.bol {
			w.off = 0
			if w.sel < len(w.pre)-1 {
				w.sel++
			}
		}
	}

	return n, nil
}
