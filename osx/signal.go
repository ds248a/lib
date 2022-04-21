package osx

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"unsafe"

	"go.uber.org/zap"
)

var (
	// support to override setting SIG_DFL so tests don't terminate early
	setDflSignal = dflSignal
)

func dflSignal(sig syscall.Signal) {
	// clearing out the sigact sets the signal to SIG_DFL
	var sigactBuf [32]uint64
	ptr := unsafe.Pointer(&sigactBuf)
	syscall.Syscall6(uintptr(syscall.SYS_RT_SIGACTION), uintptr(sig), uintptr(ptr), 0, 8, 0, 0)
}

// InterruptHandler функция, вызываемая при получении согнала SIGTERM или SIGINT.
type InterruptHandler func()

var (
	interruptRegisterMu, interruptExitMu sync.Mutex
	interruptHandlers                    = []InterruptHandler{}
)

// RegisterInterruptHandler registers a new InterruptHandler. Handlers registered
// after interrupt handing was initiated will not be executed.
func RegisterInterruptHandler(h InterruptHandler) {
	interruptRegisterMu.Lock()
	defer interruptRegisterMu.Unlock()
	interruptHandlers = append(interruptHandlers, h)
}

// HandleInterrupts calls the handler functions on receiving a SIGINT or SIGTERM.
func HandleInterrupts(lg *zap.Logger) {
	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-notifier

		interruptRegisterMu.Lock()
		ihs := make([]InterruptHandler, len(interruptHandlers))
		copy(ihs, interruptHandlers)
		interruptRegisterMu.Unlock()

		interruptExitMu.Lock()

		if lg != nil {
			lg.Info("received signal; shutting down", zap.String("signal", sig.String()))
		}

		for _, h := range ihs {
			h()
		}
		signal.Stop(notifier)
		pid := syscall.Getpid()
		// exit directly if it is the "init" process, since the kernel will not help to kill pid 1.
		if pid == 1 {
			os.Exit(0)
		}
		setDflSignal(sig.(syscall.Signal))
		syscall.Kill(pid, sig.(syscall.Signal))
	}()
}

// Exit relays to os.Exit if no interrupt handlers are running, blocks otherwise.
func Exit(code int) {
	interruptExitMu.Lock()
	os.Exit(code)
}
