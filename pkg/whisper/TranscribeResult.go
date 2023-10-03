package whisper

import (
	"C"
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type eTokenFlags uint32

const (
	TfNone    eTokenFlags = 0
	TfSpecial             = 1
)

type sTranscribeLength struct {
	CountSegments uint32
	CountTokens   uint32
}

type sTimeSpan struct {

	// The value is expressed in 100-nanoseconds ticks: compatible with System.Timespan, FILETIME, and many other things
	Ticks uint64

	/*
		operator sTimeSpanFields() const
		{
			return sTimeSpanFields{ ticks };
		}
		void operator=( uint64_t tt )
		{
			ticks = tt;
		}
		void operator=( int64_t tt )
		{
			assert( tt >= 0 );
			ticks = (uint64_t)tt;
		} */
}

type sTimeInterval struct {
	Begin sTimeSpan
	End   sTimeSpan
}

type sSegment struct {
	// Segment text, null-terminated, and probably UTF-8 encoded
	text *C.char

	// Start and end times of the segment
	Time sTimeInterval

	// These two integers define the slice of the tokens in this segment, in the array returned by iTranscribeResult.getTokens method
	FirstToken  uint32
	CountTokens uint32
}

func (this *sSegment) Text() string {
	return C.GoString(this.text)
}

type sSegmentArray []sSegment

type SToken struct {
	// Token text, null-terminated, and usually UTF-8 encoded.
	// I think for Chinese language the models sometimes outputs invalid UTF8 strings here, Unicode code points can be split between adjacent tokens in the same segment
	// More info: https://github.com/ggerganov/whisper.cpp/issues/399
	text *C.char

	// Start and end times of the token
	Time sTimeInterval
	// Probability of the token
	Probability float32

	// Probability of the timestamp token
	ProbabilityTimestamp float32

	// Sum of probabilities of all timestamp tokens
	Ptsum float32

	// Voice length of the token
	Vlen float32

	// Token id
	Id int32

	Flags eTokenFlags
}

func (this *SToken) Text() string {
	return C.GoString(this.text)
}

type sTokenArray []SToken

type iTranscribeResultVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	getSize     uintptr // ( sTranscribeLength& rdi ) HRESULT
	getSegments uintptr // () getTokens
	getTokens   uintptr // () getToken*
}

type ITranscribeResult struct {
	lpVtbl *iTranscribeResultVtbl
}

func (this *ITranscribeResult) AddRef() int32 {
	ret, _, _ := syscall.Syscall(
		this.lpVtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(this)),
		0,
		0)
	return int32(ret)
}

func (this *ITranscribeResult) Release() int32 {
	ret, _, _ := syscall.Syscall(
		this.lpVtbl.Release,
		1,
		uintptr(unsafe.Pointer(this)),
		0,
		0)
	return int32(ret)
}

func (this *ITranscribeResult) GetSize() (*sTranscribeLength, error) {

	var result sTranscribeLength

	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.getSize,
		uintptr(unsafe.Pointer(this)),
		uintptr(unsafe.Pointer(&result)),
	)

	if windows.Handle(ret) != windows.S_OK {
		fmt.Printf("iTranscribeResult.GetSize failed: %s\n", syscall.Errno(ret).Error())
		return nil, errors.New(syscall.Errno(ret).Error())
	}

	return &result, nil

}

func (this *ITranscribeResult) GetSegments(len uint32) []sSegment {

	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.getSegments,
		uintptr(unsafe.Pointer(this)),
	)

	data := unsafe.Slice((*sSegment)(unsafe.Pointer(ret)), len)

	return data
}

func (this *ITranscribeResult) GetTokens(len uint32) []SToken {

	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.getTokens,
		uintptr(unsafe.Pointer(this)),
	)

	if unsafe.Pointer(ret) != nil {
		return unsafe.Slice((*SToken)(unsafe.Pointer(ret)), len)
	} else {
		return []SToken{}
	}
}
