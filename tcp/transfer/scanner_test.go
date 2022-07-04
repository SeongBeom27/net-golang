package transfer

import (
	"bufio"
	"net"
	"reflect"
	"testing"
)

/**
bufio.Scanner : 구분자로 구분된 데이터를 읽어 들일 수 있는 Go의 표준 라이브러리
-> 바이너리 데이터의 경우 데이터 간 구분을 하기 위하여 여러 방법을 사용한다. 구분자, 특정 크기의 헤더 등 하지만 그것을 실제로 구현하는 것은 쉽지 않음
**/

const payload = "The bigger the interface, the weaker the abstraction."

func TestScanner(t *testing.T) {
	// 리스너가 하는 역할은 payload를 제공하는 일
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		_, err = conn.Write([]byte(payload))
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// 공백으로 구분된 데이터를 읽음

	// TODO : 서버에서 문자열을 읽고 있으므로(?) <- 어디서?
	// 네트워크 연결에서 데이터를 읽어들일 bufio.Scanner 생성
	scanner := bufio.NewScanner(conn)
	// ScanWords로 설정
	// ScanWords의 경우 공백, 마침표 등의 단어 경계를 ㄷ구분하는 구분자를 만날 경우 데이터를 분할
	scanner.Split(bufio.ScanWords)

	var words []string

	// 읽을 데이터가 있는 한 계속 읽음
	for scanner.Scan() {
		// 구분자로 구분된 scanner.Text()를 words string 배열에 append
		words = append(words, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		t.Error(err)
	}

	expected := []string{"The", "bigger", "the", "interface,", "the",
		"weaker", "the", "abstraction."}

	if !reflect.DeepEqual(words, expected) {
		t.Fatal("inaccurate scanned word list")
	}
	t.Logf("Scanned words: %#v", words)
}
