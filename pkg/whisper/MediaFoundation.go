package whisper

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// https://github.com/Const-me/Whisper/blob/843a2a6ca6ea47c5ac4889a281badfc808d0ea01/Whisper/API/IMediaFoundation.h

type IMediaFoundation struct {
	lpVtbl *IMediaFoundationVtbl
}

type IMediaFoundationVtbl struct {
	QueryInterface     uintptr
	AddRef             uintptr
	Release            uintptr
	loadAudioFile      uintptr // ( LPCTSTR path, bool stereo, iAudioBuffer** pp ) const;
	openAudioFile      uintptr // ( LPCTSTR path, bool stereo, iAudioReader** pp );
	loadAudioFileData  uintptr // ( const void* data, uint64_t size, bool stereo, iAudioReader** pp );  HRESULT
	listCaptureDevices uintptr // ( pfnFoundCaptureDevices pfn, void* pv );
	openCaptureDevice  uintptr // ( LPCTSTR endpoint, const sCaptureParams& captureParams, iAudioCapture** pp );
}

func (this *IMediaFoundation) AddRef() int32 {
	ret, _, _ := syscall.Syscall(
		this.lpVtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(this)),
		0,
		0)
	return int32(ret)
}

func (this *IMediaFoundation) Release() int32 {
	ret, _, _ := syscall.Syscall(
		this.lpVtbl.Release,
		1,
		uintptr(unsafe.Pointer(this)),
		0,
		0)
	return int32(ret)
}

// ( LPCTSTR path, bool stereo, iAudioBuffer** pp ) const;
func (this *IMediaFoundation) LoadAudioFile(file string, stereo bool) (*iAudioBuffer, error) {

	var buffer *iAudioBuffer

	UTFFileName, _ := windows.UTF16PtrFromString(file)

	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.loadAudioFile,
		uintptr(unsafe.Pointer(this)),
		uintptr(unsafe.Pointer(UTFFileName)),
		uintptr(1), // Todo ... Stereo !
		uintptr(unsafe.Pointer(&buffer)))

	if windows.Handle(ret) != windows.S_OK {
		fmt.Printf("loadAudioFile failed: %s\n", syscall.Errno(ret).Error())
		return nil, syscall.Errno(ret)
	}

	return buffer, nil
}

func (this *IMediaFoundation) OpenAudioFile(file string, stereo bool) (*iAudioReader, error) {

	var buffer *iAudioReader

	UTFFileName, _ := windows.UTF16PtrFromString(file)

	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.openAudioFile,
		uintptr(unsafe.Pointer(this)),
		uintptr(unsafe.Pointer(UTFFileName)),
		uintptr(1), // Todo ... Stereo !
		uintptr(unsafe.Pointer(&buffer)))

	if windows.Handle(ret) != windows.S_OK {
		fmt.Printf("openAudioFile failed: %s\n", syscall.Errno(ret).Error())
		return nil, syscall.Errno(ret)
	}

	return buffer, nil
}

func (this *IMediaFoundation) LoadAudioFileData(inbuffer *[]byte, stereo bool) (*iAudioReader, error) {

	var reader *iAudioReader

	// loadAudioFileData( const void* data, uint64_t size, bool stereo, iAudioReader** pp );
	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.loadAudioFileData,
		uintptr(unsafe.Pointer(this)),

		uintptr(unsafe.Pointer(&(*inbuffer)[0])),
		uintptr(uint64(len(*inbuffer))),
		uintptr(1), // Todo ... Stereo !
		uintptr(unsafe.Pointer(&reader)))

	if windows.Handle(ret) != windows.S_OK {
		fmt.Printf("LoadAudioFileData failed: %s\n", syscall.Errno(ret).Error())
		return nil, syscall.Errno(ret)
	}

	return reader, nil
}

// ************************************************************

type iAudioBuffer struct {
	lpVtbl *iAudioBufferVtbl
}

type iAudioBufferVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	countSamples   uintptr // returns uint32_t
	getPcmMono     uintptr // returns float*
	getPcmStereo   uintptr // returns float*
	getTime        uintptr // ( int64_t& rdi )
}

func (this *iAudioBuffer) AddRef() int32 {
	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.AddRef,
		uintptr(unsafe.Pointer(this)),
	)
	return int32(ret)
}

func (this *iAudioBuffer) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.Release,
		uintptr(unsafe.Pointer(this)),
	)
	return int32(ret)
}

func (this *iAudioBuffer) CountSamples() (uint32, error) {

	ret, _, err := syscall.SyscallN(
		this.lpVtbl.countSamples,
		uintptr(unsafe.Pointer(this)),
	)

	if err != 0 {
		return 0, errors.New(err.Error())
	}

	return uint32(ret), nil
}

// ************************************************************

type iAudioReader struct {
	lpVtbl *iAudioReaderVtbl
}

type iAudioReaderVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	getDuration     uintptr // ( int64_t& rdi )
	getReader       uintptr // ( IMFSourceReader** pp )
	requestedStereo uintptr // ()
}

func (this *iAudioReader) AddRef() int32 {
	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.AddRef,
		uintptr(unsafe.Pointer(this)),
	)
	return int32(ret)
}

func (this *iAudioReader) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.Release,
		uintptr(unsafe.Pointer(this)),
	)
	return int32(ret)
}

func (this *iAudioReader) GetDuration() (uint64, error) {

	var rdi int64

	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.getDuration,
		uintptr(unsafe.Pointer(this)),
		uintptr(unsafe.Pointer(&rdi)),
	)

	if windows.Handle(ret) != windows.S_OK {
		fmt.Printf("LoadAudioFileData failed: %s\n", syscall.Errno(ret).Error())
		return 0, syscall.Errno(ret)
	}

	return uint64(rdi), nil
}

// ************************************************************

type iAudioCapture struct {
	lpVtbl *iAudioCaptureVtbl
}

type iAudioCaptureVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	getReader      uintptr // ( IMFSourceReader** pp )
	getParams      uintptr // returns sCaptureParams&
}
