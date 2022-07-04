package connection

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestDialContextCancel(t *testing.T) {
	// context.WithCancel 함수를 이용하여 context와 context를 취소할 수 있는 함수를 받음
	ctx, cancel := context.WithCancel(context.Background())
	sync := make(chan struct{})

	// 수동으로 context를 취소하기 때문에 closure를 만들어서 별도로 연결 시도를 처리 하기 위한 고루틴 시작
	go func() {
		defer func() { sync <- struct{}{} }()

		var d net.Dialer
		d.Control = func(_, _ string, _ syscall.RawConn) error {
			time.Sleep(time.Second)
			return nil
		}

		conn, err := d.DialContext(ctx, "tcp", "10.0.0.1:80")
		if err != nil {
			t.Log(err)
			return
		}

		conn.Close()
		t.Error("connection did not time out")
	}()

	// dialer가 연결 시도를 하고 원격 노드의 handshake가 끝나면,
	// context를 취소하기 위해 cancel 함수를 호출
	cancel()
	<-sync

	// 결과 : DialContext 메서드는 즉시 nil이 아닌 에러를 반환하고 고루틴을 종료

	// Err 메서드는 context.Canceled를 반환
	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual: %q", ctx.Err())
	}
}
