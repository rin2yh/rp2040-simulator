//go:build !tinygo

package bridge

import (
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
	"sync"
	"time"
)

var (
	mu         sync.Mutex
	connection *rpc.Client
	launching  bool
)

func Call(method string, args, reply any) error {
	mu.Lock()
	defer mu.Unlock()
	if connection == nil {
		var (
			err      error
			launched bool
		)
		for range 30 {
			connection, err = rpc.Dial("tcp", Address)
			if err == nil {
				break
			}
			if !launched {
				if err := launchEmulator(); err != nil {
					return err
				}
				launched = true
			}
			time.Sleep(100 * time.Millisecond)
		}
		if err != nil {
			return err
		}
	}
	if err := connection.Call(method, args, reply); err != nil {
		_ = connection.Close()
		connection = nil
		return err
	}
	return nil
}

func launchEmulator() error {
	if launching {
		return nil
	}
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locate application executable: %w", err)
	}
	cmd := exec.Command(executable)
	cmd.Env = append(os.Environ(), EmulatorProcess+"=1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start emulator: %w", err)
	}
	launching = true
	go func() {
		_ = cmd.Wait()
		mu.Lock()
		launching = false
		mu.Unlock()
	}()
	return nil
}
