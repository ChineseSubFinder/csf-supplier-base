package kill_process

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	modkernel32          = syscall.NewLazyDLL("kernel32.dll")
	procCreateToolhelp32 = modkernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First   = modkernel32.NewProc("Process32FirstW")
	procProcess32Next    = modkernel32.NewProc("Process32NextW")
	procOpenProcess      = modkernel32.NewProc("OpenProcess")
	procTerminateProcess = modkernel32.NewProc("TerminateProcess")
	procCloseHandle      = modkernel32.NewProc("CloseHandle")
)

const (
	TH32CS_SNAPPROCESS = 0x00000002
	PROCESS_TERMINATE  = 0x0001
)

type PROCESSENTRY32 struct {
	Size              uint32
	CntUsage          uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	CntThreads        uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [syscall.MAX_PATH]uint16
}

// Kill 只在 Windows 10 64 位 上测试通过，32 位未测试，可能出现乱码进程名无法杀死的情况
func Kill(processName string) error {
	processID, err := GetProcessID(processName)
	if err != nil {
		return err
	}

	err = TerminateProcess(processID)
	if err != nil {
		return err
	}

	return nil
}

func GetProcessID(processName string) (uint32, error) {
	hSnapshot, _, _ := procCreateToolhelp32.Call(TH32CS_SNAPPROCESS, 0)
	if hSnapshot == uintptr(syscall.InvalidHandle) {
		return 0, fmt.Errorf("CreateToolhelp32Snapshot failed")
	}
	defer procCloseHandle.Call(hSnapshot)

	var pe32 PROCESSENTRY32
	pe32.Size = uint32(unsafe.Sizeof(pe32))
	ret, _, _ := procProcess32First.Call(hSnapshot, uintptr(unsafe.Pointer(&pe32)))
	if ret == 0 {
		return 0, fmt.Errorf("Process32First failed")
	}

	for {
		exeFile := syscall.UTF16ToString(pe32.ExeFile[:])
		if exeFile == processName {
			return pe32.ProcessID, nil
		}

		ret, _, _ = procProcess32Next.Call(hSnapshot, uintptr(unsafe.Pointer(&pe32)))
		if ret == 0 {
			break
		}
	}

	return 0, fmt.Errorf("process not found")
}

func TerminateProcess(processID uint32) error {
	hProcess, _, _ := procOpenProcess.Call(uintptr(PROCESS_TERMINATE), 0, uintptr(processID))
	if hProcess == 0 {
		return fmt.Errorf("OpenProcess failed")
	}
	defer procCloseHandle.Call(hProcess)

	ret, _, _ := procTerminateProcess.Call(hProcess, 0)
	if ret == 0 {
		return fmt.Errorf("TerminateProcess failed")
	}

	return nil
}
