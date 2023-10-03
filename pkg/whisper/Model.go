package whisper

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// External - Go version of the struct
type Model struct {
	cStruct *_IModel
	setup   *sModelSetup
}

// Internal - C Version of the structs
type _IModel struct {
	lpVtbl *IModelVtbl
}

// https://github.com/Const-me/Whisper/blob/master/Whisper/API/iContext.cl.h
type IModelVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	createContext    uintptr //( iContext** pp ) = 0;
	tokenize         uintptr /* HRESULT __stdcall tokenize( const char* text, pfnDecodedTokens pfn, void* pv ); */
	isMultilingual   uintptr //() = 0;
	getSpecialTokens uintptr //( SpecialTokens& rdi ) = 0;
	stringFromToken  uintptr //( whisper_token token ) = 0;
	clone            uintptr //( iModel** rdi ) = 0;
}

func NewModel(setup *sModelSetup, cstruct *_IModel) *Model {
	this := Model{}
	this.setup = setup
	this.cStruct = cstruct
	return &this
}

func (this *Model) AddRef() int32 {
	ret, _, _ := syscall.Syscall(
		this.cStruct.lpVtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(this.cStruct)),
		0,
		0)
	return int32(ret)
}

func (this *Model) Release() int32 {
	ret, _, _ := syscall.Syscall(
		this.cStruct.lpVtbl.Release,
		1,
		uintptr(unsafe.Pointer(this.cStruct)),
		0,
		0)
	return int32(ret)
}

func (this *Model) CreateContext() (*IContext, error) {
	var context *IContext

	/*
		ret, _, err := syscall.Syscall(
			this.cStruct.lpVtbl.createContext,
			2, // Why was this 1, rather than 2 ?? 1 seemed to work fine
			uintptr(unsafe.Pointer(this.cStruct)),
			uintptr(unsafe.Pointer(&context)),
			0)*/
	ret, _, err := syscall.SyscallN(
		this.cStruct.lpVtbl.createContext,
		uintptr(unsafe.Pointer(this.cStruct)),
		uintptr(unsafe.Pointer(&context)))

	if windows.Handle(ret) != windows.S_OK {
		fmt.Printf("createContext failed: %w", err.Error())
	}

	if windows.Handle(ret) != windows.S_OK {
		return nil, fmt.Errorf("loadModel failed: %w", err)
	}

	return context, nil
}

func (this *Model) IsMultilingual() bool {
	ret, _, _ := syscall.SyscallN(
		this.cStruct.lpVtbl.isMultilingual,
		uintptr(unsafe.Pointer(this.cStruct)),
	)

	return bool(windows.Handle(ret) == windows.S_OK)
}

func (this *Model) Clone() (*_IModel, error) {

	if this.setup.isFlagSet(gmf_Cloneable) {
		return nil, errors.New("Model is not cloneable")
	}
	//this.Cloneable ?

	var modelptr *_IModel

	ret, _, _ := syscall.SyscallN(
		this.cStruct.lpVtbl.clone,
		uintptr(unsafe.Pointer(this.cStruct)),
		uintptr(unsafe.Pointer(&modelptr)),
	)

	if windows.Handle(ret) == windows.S_OK {
		return modelptr, nil
	} else {
		return nil, errors.New("Model.Clone() failed : " + syscall.Errno(ret).Error())
	}
}
