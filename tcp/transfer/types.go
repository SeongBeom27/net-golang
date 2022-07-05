package transfer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	// 정의할 메시지 타입을 나타내는 상수
	BinaryType uint8 = iota + 1
	StringType

	// 보안상의 문제로 인해 최대 페이로드 크기를 지정해줄 필요가 있음
	MaxPayloadSize uint32 = 10 << 20 // 10MB
)

var ErrMaxOayloadSize = errors.New("maximum payload size exceeded")

// 각 타입별 메시지들이 구현해야 하는 payload라는 이름의 인터페이스 정의
type Payload interface {
	fmt.Stringer
	// 각 타입별 메시지를 reader로부터 읽을 수 있음
	io.ReaderFrom
	// 각 타입별 메시지를 writer에 쓸 수 있음
	io.WriterTo
	Bytes() []byte
}

type Binary []byte

func (m Binary) Bytes() []byte  { return m }
func (m Binary) String() string { return string(m) }

func (m Binary) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, BinaryType) // 1byte 타입
	if err != nil {
		return 0, err
	}
	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(m))) // 4byte 크기
	if err != nil {
		return n, err
	}
	n += 4

	o, err := w.Write(m) // payload

	return n + int64(o), err
}
