package mock

import (
	"io"
	"testing"

	goloader "wallet/pkg/core/loader"
)

type Option func(goloader.ILoader, IMockEnv)

func NewMockEnv(t testing.TB, injector goloader.ILoader, opts ...Option) IMockEnv {
	me := &mockenv{t, nil}
	for _, opt := range opts {
		opt(injector, me)
	}
	return me
}

type IMockEnv interface {
	testing.TB
	AppendCloser(closer ...io.Closer)
	Close()
}

type mockenv struct {
	testing.TB
	closers []io.Closer
}

func (m *mockenv) AppendCloser(closer ...io.Closer) {
	if len(closer) > 0 {
		closer = append(m.closers, closer...)
	}
}

func (m *mockenv) Close() {
	for _, closer := range m.closers {
		if err := closer.Close(); err != nil {
			m.Error(err)
		}
	}
}
