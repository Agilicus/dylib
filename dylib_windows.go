package dylib

import (
	"syscall"
)

type LazyDLL struct {
	*syscall.LazyDLL
}

func NewLazyDLL(name string) *LazyDLL {
	l := new(LazyDLL)
	l.LazyDLL = syscall.NewLazyDLL(name)
	if err := l.Load(); err != nil {
		return l
	}
	return l
}

// 定义一个相同的
type LazyProc struct {
	lzProc *syscall.LazyProc
	lzdll  *LazyDLL
}

func (d *LazyDLL) NewProc(name string) *LazyProc {
	l := new(LazyProc)
	l.lzProc = d.LazyDLL.NewProc(name)
	l.lzdll = d
	return l
}

func (d *LazyDLL) Close() {
	if d.Handle() != 0 {
		syscall.FreeLibrary(syscall.Handle(d.Handle()))
	}
}

func (d *LazyDLL) call(proc *LazyProc, a ...uintptr) (r1, r2 uintptr, lastErr error) {
	return proc.CallOriginal(a...)
}

func (p *LazyProc) Addr() uintptr {
	return p.lzProc.Addr()
}

func (p *LazyProc) Find() error {
	return p.lzProc.Find()
}

func (p *LazyProc) Call(a ...uintptr) (r1, r2 uintptr, lastErr error) {
	return p.lzdll.call(p, a...)
}

func (p *LazyProc) CallOriginal(a ...uintptr) (r1, r2 uintptr, lastErr error) {
	return p.lzProc.Call(a...)
}
