package whisper

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type uuid [16]byte

type eResultFlags uint32

const (
	RfNone eResultFlags = 0

	// Return individual tokens in addition to the segments
	RfTokens = 1

	// Return timestamps
	RfTimestamps = 2

	// Create a new COM object for the results.
	// Without this flag, the context returns a pointer to the COM object stored in the context.
	// The content of that object is replaced every time you call IContext.getResults method
	RfNewObject = 0x100
)

type IContextVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	RunFull           uintptr
	RunStreamed       uintptr
	RunCapture        uintptr
	GetResults        uintptr
	DetectSpeaker     uintptr
	GetModel          uintptr
	FullDefaultParams uintptr
	TimingsPrint      uintptr
	TimingsReset      uintptr
}

type IContext struct {
	lpVtbl *IContextVtbl
}

//type sFullParams struct{}

// type iAudioBuffer struct{}
type sProgressSink struct {
	pfn uintptr
	pv  uintptr
}

// type iAudioReader struct{}
type sCaptureCallbacks struct{}

// type iAudioCapture struct{}
// type eResultFlags int32
// type iTranscribeResult struct{}
// type sTimeInterval struct{}
type eSpeakerChannel int32

//type eSamplingStrategy int32

// Create a new IContext instance
func newIContext() *IContext {
	return &IContext{
		lpVtbl: &IContextVtbl{
			QueryInterface:    0,
			AddRef:            0,
			Release:           0,
			RunFull:           0,
			RunStreamed:       0,
			RunCapture:        0,
			GetResults:        0,
			DetectSpeaker:     0,
			GetModel:          0,
			FullDefaultParams: 0,
			TimingsPrint:      0,
			TimingsReset:      0,
		},
	}
}

func (context *IContext) TimingsPrint() error {

	//  TimingsPrint();
	ret, _, _ := syscall.SyscallN(
		context.lpVtbl.TimingsPrint,
		uintptr(unsafe.Pointer(context)),
	)

	if windows.Handle(ret) != windows.S_OK {
		fmt.Printf("RunFull failed: %s\n", syscall.Errno(ret).Error())
		return errors.New(syscall.Errno(ret).Error())
	}

	return nil
}

// Run the entire model: PCM -> log mel spectrogram -> encoder -> decoder -> text
// Uses the specified decoding strategy to obtain the text.
func (context *IContext) RunFull(params *FullParams, buffer *iAudioBuffer) error {

	//  runFull( const sFullParams& params, const iAudioBuffer* buffer );
	ret, _, _ := syscall.SyscallN(
		context.lpVtbl.RunFull,
		uintptr(unsafe.Pointer(context)),

		uintptr(unsafe.Pointer(params.cStruct)),
		uintptr(unsafe.Pointer(buffer)),
	)

	if windows.Handle(ret) != windows.S_OK {
		fmt.Printf("RunFull failed: %s\n", syscall.Errno(ret).Error())
		return errors.New(syscall.Errno(ret).Error())
	}

	return nil
}

func (context *IContext) RunStreamed(params *FullParams, reader *iAudioReader) error {

	cb := sProgressSink{}

	//   runStreamed( const sFullParams& params, const sProgressSink& progress, const iAudioReader* reader );
	ret, _, _ := syscall.SyscallN(
		context.lpVtbl.RunStreamed,
		uintptr(unsafe.Pointer(context)),
		uintptr(unsafe.Pointer(params.cStruct)),
		uintptr(unsafe.Pointer(&cb)), // No progress cb yet
		uintptr(unsafe.Pointer(reader)),
	)

	if windows.Handle(ret) != windows.S_OK {
		fmt.Printf("RunStreamed failed: %s\n", syscall.Errno(ret).Error())
		return errors.New(syscall.Errno(ret).Error())
	}

	return nil
}

func (this *IContext) AddRef() int32 {
	ret, _, _ := syscall.Syscall(
		this.lpVtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(this)),
		0,
		0)
	return int32(ret)
}

func (this *IContext) Release() int32 {
	ret, _, _ := syscall.Syscall(
		this.lpVtbl.Release,
		1,
		uintptr(unsafe.Pointer(this)),
		0,
		0)
	return int32(ret)
}

/*
https://github.com/Const-me/Whisper/blob/f6f743c7b3570b85ccf47f74b84e06a73667ef3e/Whisper/Whisper/ContextImpl.misc.cpp

Returns E_POINTER if null pointer provided in params
Initialises params to all 0
sets values in struct, does not malloc
*/
func (context *IContext) FullDefaultParams(strategy eSamplingStrategy) (*FullParams, error) {

	/*
		ERR : unreadable Only part of a ReadProcessMemory or WriteProcessMemory request was completed
		 * not related to stratergy ... tested 0, 1 and 2 ... 2 produced E_INVALIDARG as expected
		 * not a nil ptr to params ... nil poitner produced E_POINTER as expected
		 * params seems to return 0x4000
		 * !!!!!  FullParams is not a com interface !!!
		 *   so no lpVtbl *FullParamsVtbl , no queryinterface, addref etc
	*/

	params := _newFullParams_cStruct()
	//params := &[160]byte{}

	ret, _, _ := syscall.SyscallN(
		context.lpVtbl.FullDefaultParams,
		uintptr(unsafe.Pointer(context)),
		uintptr(strategy),
		uintptr(unsafe.Pointer(params)),
	)

	// nil ptr should be 0x80004003L
	// unsafe.Pointer(0xc00011dc28)
	// unsafe.Pointer(0x4000)

	if windows.Handle(ret) != windows.S_OK {
		fmt.Printf("FullDefaultParams failed: %s\n", syscall.Errno(ret).Error())
		return nil, syscall.Errno(ret)

	}

	if params == nil {
		return nil, errors.New("FullDefaultParams did not return params")
	}
	ParamObj := NewFullParams(params)
	// ParamObj.SetLanguage(Polish)

	if ParamObj.TestDefaultsOK() {
		return ParamObj, nil
	}

	return nil, nil
}

func (context *IContext) GetModel() (*_IModel, error) {

	var modelptr *_IModel

	// getModel( iModel** pp );
	ret, _, _ := syscall.SyscallN(
		context.lpVtbl.GetModel,
		uintptr(unsafe.Pointer(context)),
		uintptr(unsafe.Pointer(&modelptr)),
	)

	if windows.Handle(ret) != windows.S_OK {
		fmt.Printf("FullDefaultParams failed: %s\n", syscall.Errno(ret).Error())
		return nil, syscall.Errno(ret)
	}

	if modelptr == nil {
		return nil, errors.New("loadModel did not return a Model")
	}

	if modelptr.lpVtbl == nil {
		return nil, errors.New("loadModel method table is nil")
	}

	return modelptr, nil
}

// ************************************************************************************************************************************************
// Not really implemented / tested
// ************************************************************************************************************************************************

func (context *IContext) RunCapture(params *FullParams, callbacks *sCaptureCallbacks, reader *iAudioCapture) uintptr {
	ret, _, _ := syscall.SyscallN(
		context.lpVtbl.RunCapture,
		//3,
		uintptr(unsafe.Pointer(context)),
		uintptr(unsafe.Pointer(params)),
		uintptr(unsafe.Pointer(callbacks)),
		uintptr(unsafe.Pointer(reader)),
	)
	return ret
}

func (context *IContext) GetResults(flags eResultFlags, pp **ITranscribeResult) uintptr {
	ret, _, _ := syscall.Syscall(
		context.lpVtbl.GetResults,
		3,
		uintptr(unsafe.Pointer(context)),
		uintptr(flags),
		uintptr(unsafe.Pointer(pp)),
	)
	return ret
}

func (context *IContext) DetectSpeaker(time *sTimeInterval, result *eSpeakerChannel) uintptr {
	ret, _, _ := syscall.Syscall(
		context.lpVtbl.DetectSpeaker,
		3,
		uintptr(unsafe.Pointer(context)),
		uintptr(unsafe.Pointer(time)),
		uintptr(unsafe.Pointer(result)),
	)
	return ret
}
