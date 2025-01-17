package ratelimiter

import (
	"testing"
	"time"

	"github.com/charmbracelet/wish/testsession"
	"github.com/gliderlabs/ssh"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {
	s := &ssh.Server{
		Handler: Middleware(NewRateLimiter(rate.Limit(10), 4, 1))(func(s ssh.Session) {
			// noop
		}),
	}

	addr := testsession.Listen(t, s)

	g := errgroup.Group{}
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			sess, err := testsession.NewClientSession(t, addr, nil)
			if err != nil {
				t.Fatalf("expected no errors, got %v", err)
			}
			if err := sess.Run(""); err != nil {
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err == nil {
		t.Fatal("expected error, got nil")
	}

	// after some time, it should reset and pass again
	time.Sleep(100 * time.Millisecond)
	sess, err := testsession.NewClientSession(t, addr, nil)
	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}
	if err := sess.Run(""); err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}
}
