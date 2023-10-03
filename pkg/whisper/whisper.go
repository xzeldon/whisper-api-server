//go:build windows
// +build windows

package whisper

import (
	"C"
	"errors"
	"syscall"
	"unsafe"

	// Using lxn/win because its COM functions expose raw HRESULTs
	"golang.org/x/sys/windows"
)
import (
	"fmt"
)

/*
	eModelImplementation - TranscribeStructs.h

	// GPGPU implementation based on Direct3D 11.0 compute shaders
	GPU = 1,

	// A hybrid implementation which uses DirectCompute for encode, and decodes on CPU
	// Not implemented in the published builds of the DLL. To enable, change BUILD_HYBRID_VERSION macro to 1
	Hybrid = 2,

	// A reference implementation which uses the original GGML CPU-running code
	// Not implemented in the published builds of the DLL. To enable, change BUILD_BOTH_VERSIONS macro to 1
	Reference = 3,
*/

// https://learn.microsoft.com/en-us/windows/win32/seccrypto/common-hresult-values
// https://pkg.go.dev/golang.org/x/sys/windows
const (
	E_INVALIDARG                      = 0x80070057
	ERROR_HV_CPUID_FEATURE_VALIDATION = 0xC0350038

	DLLName = "whisper.dll"
)

type Libwhisper struct {
	dll            *syscall.LazyDLL
	ver            WinVersion
	existing_model map[string]*Model

	proc_setupLogger         *syscall.LazyProc
	proc_loadModel           *syscall.LazyProc
	proc_initMediaFoundation *syscall.LazyProc
	// proc_findLanguageKeyW      *syscall.LazyProc
	// proc_findLanguageKeyA      *syscall.LazyProc
	// proc_getSupportedLanguages *syscall.LazyProc
}

var singleton_whisper *Libwhisper = nil

func New(level eLogLevel, flags eLogFlags, cb *any) (*Libwhisper, error) {
	if singleton_whisper != nil {
		return singleton_whisper, nil
	}

	var err error
	this := &Libwhisper{}

	this.ver, err = GetFileVersion(DLLName)
	if err != nil {
		return nil, err
	}

	if this.ver.Major < 1 && this.ver.Minor < 9 {
		return nil, errors.New("This library requires whisper.dll version 1.9 or higher.") // or less than 1.11 for now .. because the API changed
	}

	this.dll = syscall.NewLazyDLL(DLLName) // Todo wrap this in a class, check file exists, handle errors ... you know, just a few things.. AKA Stop being lazy

	this.proc_setupLogger = this.dll.NewProc("setupLogger")
	this.proc_loadModel = this.dll.NewProc("loadModel")
	this.proc_initMediaFoundation = this.dll.NewProc("initMediaFoundation")
	/*
		this.proc_findLanguageKeyW = this.dll.NewProc("findLanguageKeyW")
		this.proc_findLanguageKeyA = this.dll.NewProc("findLanguageKeyA")
		this.proc_getSupportedLanguages = this.dll.NewProc("getSupportedLanguages")
	*/

	ok, err := this._setupLogger(level, flags, cb)
	if !ok {
		return nil, errors.New("Logger Error : " + err.Error())
	}

	this.existing_model = make(map[string]*Model)
	singleton_whisper = this

	return singleton_whisper, nil
}

func (this *Libwhisper) Version() string {
	return fmt.Sprintf("%d.%d.%d.%d.", this.ver.Major, this.ver.Minor, this.ver.Patch, this.ver.Build)
}

func (this *Libwhisper) SupportsMultiThread() bool {
	return this.ver.Major >= 1 && this.ver.Minor >= 10
}

func (this *Libwhisper) _setupLogger(level eLogLevel, flags eLogFlags, cb *any) (bool, error) {

	setup := sLoggerSetup{}
	setup.sink = 0
	setup.context = 0
	setup.level = level
	setup.flags = flags

	if cb != nil {
		setup.sink = syscall.NewCallback(cb)
	}

	res, _, err := this.proc_setupLogger.Call(uintptr(unsafe.Pointer(&setup)))

	if windows.Handle(res) == windows.S_OK {
		return true, nil
	} else {
		return false, err
	}
}

func (this *Libwhisper) LoadModel(path string, aGPU ...string) (*Model, error) {
	var modelptr *_IModel

	whisperpath, _ := windows.UTF16PtrFromString(path)

	GPU := ""
	if len(aGPU) == 1 {
		GPU = aGPU[0]
	}

	setup := ModelSetup(gmf_Cloneable, GPU)

	// Construct our map hash
	singleton_hash := GPU + "|" + path
	if this.existing_model[singleton_hash] != nil {
		ClonedModel, err := this.existing_model[singleton_hash].Clone()
		if ClonedModel != nil {
			return NewModel(setup, ClonedModel), nil
		} else {
			return nil, err
		}
	}

	obj, _, _ := this.proc_loadModel.Call(uintptr(unsafe.Pointer(whisperpath)), uintptr(unsafe.Pointer(setup.AsCType())), uintptr(unsafe.Pointer(nil)), uintptr(unsafe.Pointer(&modelptr)))

	if windows.Handle(obj) != windows.S_OK {
		fmt.Printf("loadModel failed: %s\n", syscall.Errno(obj).Error())
		return nil, fmt.Errorf("loadModel failed: %s", syscall.Errno(obj))
	}

	if modelptr == nil {
		return nil, errors.New("loadModel did not return a Model")
	}

	if modelptr.lpVtbl == nil {
		return nil, errors.New("loadModel method table is nil")
	}

	model := NewModel(setup, modelptr)

	this.existing_model[singleton_hash] = model

	return model, nil
}

func (this *Libwhisper) InitMediaFoundation() (*IMediaFoundation, error) {

	var mediafoundation *IMediaFoundation

	// initMediaFoundation( iMediaFoundation** pp );
	obj, _, _ := this.proc_initMediaFoundation.Call(uintptr(unsafe.Pointer(&mediafoundation)))

	if windows.Handle(obj) != windows.S_OK {
		fmt.Printf("initMediaFoundation failed: %s\n", syscall.Errno(obj).Error())
		return nil, fmt.Errorf("initMediaFoundation failed: %s", syscall.Errno(obj))
	}

	if mediafoundation.lpVtbl == nil {
		return nil, errors.New("initMediaFoundation method table is nil")
	}

	return mediafoundation, nil
}
