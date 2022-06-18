package tcp

import (
	"io"
	"net"
	"testing"
)

func TestDial(t *testing.T) {
	// 랜덤 포트에 리스너 생성 ( go에서 사용 가능한 포트를 임의로 설정 )
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})

	// 핸들러
	go func() {
		defer func() { done <- struct{}{} }()

		for {
			conn, err := listener.Accept()

			if err != nil {
				t.Log(err)
				return
			}
			go func(c net.Conn) {
				defer func() {
					c.Close()
					done <- struct{}{}
				}()

				buf := make([]byte, 1024)
				for {
					// Read한 내용은 buf 배열에 저장되고 n은 읽은 크기의 인덱스를 획득
					n, err := c.Read(buf)
					if err != nil {
						// FIN 패킷을 받고나면 Read 메서드는 io.EOF 에러를 반환
						// 리스너 측에서는 반대편 연결이 종료되었다는 의미
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
					t.Logf("received: %q", buf[:n])
				}
			}(conn)
		}
	}()

	// tcp와 같은 네트워크의 종류와 ip 주소, 포트의 조합을 매개변수로 받는다는 점에서 net.Listen 함수와 유사
	// 두 번째 매개변수로 받은 IP주소, 포트를 이용하여 리스너로 연결을 시도
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	// Connection 핸들러는 연결 객체의 Close 메서드를 호출하여 종료
	// Close 메서드는 FIN 패킷을 전송하여 TCP의 우아한 조욜를 마무리
	conn.Close()
	<-done

	// Listener를 종료함과 동시에 Listener의 Accept 메서드는 즉시 블로킹이 해제되고 에러를 반환
	// 실패한 에러가 아닌 그냥 로깅하고 넘어가면 됨
	listener.Close()
	<-done
}
