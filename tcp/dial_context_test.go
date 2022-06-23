package tcp

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestDialContext(t *testing.T) {
	// 5초 후 데드라인이 지나는 콘텍스트를 만들기 위해 현재 시간으로부터 5초 뒤의 시간을 저장
	dl := time.Now().Add(5 * time.Second)

	// WithDeadline 함수를 이용하여 context와 cancel 함수를 생성하고 위에서 생성한 데드라인을 설정
	ctx, cancel := context.WithDeadline(context.Background(), dl)

	// context가 가비지컬렉션이 되도록 cancel 함수를 defer로 호출
	defer cancel()

	var d net.Dialer // DialContext는 Dialer의 메서드

	// Dialer의 Control 함수를 Overriding하여 연결을 context의 데드라인을 간신히 초과 (tijme.Milisecond)하는 정도로 지연
	d.Control = func(_, _ string, _ syscall.RawConn) error {
		// context의 데드라인이 지나기 위해 충분히 긴 시간 동안 대기
		time.Sleep(5*time.Second + time.Millisecond)
		return nil
	}

	// DialContext 함수의 첫 번째 매개변수로 위에서 생성한 context(데드라인 5초를 가지고 있는)를 전달
	conn, err := d.DialContext(ctx, "tcp", "10.0.0.0:80")
	if err == nil {
		conn.Close()
		t.Fatal("connection did not time out")
	}

	nErr, ok := err.(net.Error)
	if !ok {
		t.Error(err)
	} else {
		if !nErr.Timeout() {
			t.Errorf("error is not a timeout: %v", err)
		}
	}

	// 테스트의 끝 부분의 에러 처리는 데드라인이 context를 제대로 취소하였는지
	// cancel 함수 호출에 문제는 없었는지 확인
	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("expected deadline exceeded; actual: %v", ctx.Err())
	}
}
