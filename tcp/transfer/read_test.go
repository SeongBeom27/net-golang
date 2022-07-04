package transfer

import (
	"crypto/rand"
	"io"
	"net"
	"testing"
)

func TestReadIntoBuffer(t *testing.T) {

	// Client가 읽어들일 16MB 랜덤 데이터를 가진 payload 생성
	payload := make([]byte, 1<<24) // 16MB
	_, err := rand.Read(payload)
	if err != nil {
		t.Fatal(err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	// 리스너를 시작하교 연결 수신을 대기하기 위한 고루틴 생성
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}
		defer conn.Close()

		// 연결 수신 후 서버는 네트워크 연결로 payload를 전부 write
		_, err = conn.Write(payload)
		if err != nil {
			t.Error(err)
		}
	}()

	// TODO : Dial 동작 확인 필요
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Error(err)
	}

	// Client에는 512KB의 버퍼를 가지고 있음
	buf := make([]byte, 1<<19) // 512KB

	// 512KB 보다 데이터가 클 경우 계속해서 512KB 씩 데이터를 받다가
	// err가 nil이 아닌 경우, 즉 데이터의 마지막 부분의 경우 break로 빠져나옴
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}

		t.Logf("read %d bytes", n) // buf[:n]은 conn 객체에서 읽은 데이터
	}

	conn.Close()
}
