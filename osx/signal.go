package osx

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"unsafe"

	"go.uber.org/zap"
)

// Переопределение сигнала по умолчанию.
// Используется в тестах.
var setDflSignal = dflSignal

func dflSignal(sig syscall.Signal) {
	var sigactBuf [32]uint64
	ptr := unsafe.Pointer(&sigactBuf)
	syscall.Syscall6(uintptr(syscall.SYS_RT_SIGACTION), uintptr(sig), uintptr(ptr), 0, 8, 0, 0)
}

// InterruptHandler функция обработчик, вызывается при регистрации согнала SIGTERM или SIGINT.
type InterruptHandler func()

var interruptRegisterMu, interruptExitMu sync.Mutex
var interruptHandlers = []InterruptHandler{}

// RegisterInterruptHandler регистрирует функцию обработчика системных событий.
func RegisterInterruptHandler(h InterruptHandler) {
	interruptRegisterMu.Lock()
	defer interruptRegisterMu.Unlock()
	interruptHandlers = append(interruptHandlers, h)
}

// HandleInterrupts запускает обработку системных событий прерывания работы приложения.
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
		if pid == 1 {
			os.Exit(0)
		}

		setDflSignal(sig.(syscall.Signal))
		syscall.Kill(pid, sig.(syscall.Signal))
	}()
}

func Exit(code int) {
	interruptExitMu.Lock()
	os.Exit(code)
}
