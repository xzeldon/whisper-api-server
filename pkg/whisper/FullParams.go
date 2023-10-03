package whisper

import (
	"syscall"
	"unsafe"
)

// https://github.com/Const-me/Whisper/blob/master/Whisper/API/sFullParams.h
// https://github.com/Const-me/Whisper/blob/master/WhisperNet/API/Parameters.cs

type eSamplingStrategy uint32

const (
	SsGreedy eSamplingStrategy = iota
	SsBeamSearch
	SsINVALIDARG
)

type eFullParamsFlags uint32

const (
	FlagNone            eFullParamsFlags = 0
	FlagTranslate                        = 1 << 0
	FlagNoContext                        = 1 << 1
	FlagSingleSegment                    = 1 << 2
	FlagPrintSpecial                     = 1 << 3
	FlagPrintProgress                    = 1 << 4
	FlagPrintRealtime                    = 1 << 5
	FlagPrintTimestamps                  = 1 << 6
	FlagTokenTimestamps                  = 1 << 7 // Experimental
	FlagSpeedupAudio                     = 1 << 8
)

type EWhisperHWND uintptr

const (
	S_OK    EWhisperHWND = 0
	S_FALSE EWhisperHWND = 1
)

type FullParams struct {
	cStruct *_FullParams
}

func (this *FullParams) CpuThreads() int32 {
	if this == nil {
		return 0
	} else if this.cStruct == nil {
		return 0
	}

	return this.cStruct.cpuThreads
}

func (this *FullParams) setCpuThreads(val int32) {
	if this == nil {
		return
	} else if this.cStruct == nil {
		return
	}

	this.cStruct.cpuThreads = val
}

func (this *FullParams) SetMaxTextCTX(val int32) {
	this.cStruct.n_max_text_ctx = val
}

func (this *FullParams) AddFlags(newflag eFullParamsFlags) {
	if this == nil {
		return
	} else if this.cStruct == nil {
		return
	}

	this.cStruct.Flags = this.cStruct.Flags | newflag
}

func (this *FullParams) RemoveFlags(newflag eFullParamsFlags) {
	if this == nil {
		return
	} else if this.cStruct == nil {
		return
	}

	this.cStruct.Flags = this.cStruct.Flags ^ newflag
}

/*using pfnNewSegment = HRESULT( __cdecl* )( iContext* ctx, uint32_t n_new, void* user_data ) noexcept;*/
type NewSegmentCallback_Type func(context *IContext, n_new uint32, user_data unsafe.Pointer) EWhisperHWND

func (this *FullParams) SetNewSegmentCallback(cb NewSegmentCallback_Type) {
	if this == nil {
		return
	} else if this.cStruct == nil {
		return
	}
	this.cStruct.new_segment_callback = syscall.NewCallback(cb)
}

/*
Return S_OK to proceed, or S_FALSE to stop the process
*/
type EncoderBeginCallback_Type func(context *IContext, user_data unsafe.Pointer) EWhisperHWND

func (this *FullParams) SetEncoderBeginCallback(cb EncoderBeginCallback_Type) {
	if this == nil {
		return
	} else if this.cStruct == nil {
		return
	}

	this.cStruct.encoder_begin_callback = syscall.NewCallback(cb)
}

func (this *FullParams) TestDefaultsOK() bool {
	if this == nil {
		return false
	} else if this.cStruct == nil {
		return false
	}

	if this.cStruct.n_max_text_ctx != 16384 {
		return false
	}

	if this.cStruct.Flags != (FlagPrintProgress | FlagPrintTimestamps) {
		return false
	}

	if this.cStruct.thold_pt != 0.01 {
		return false
	}

	if this.cStruct.thold_ptsum != 0.01 {
		return false
	}

	if this.cStruct.Language != English {
		return false
	}

	// Todo ... why do these not line up as expected.. is our struct out of alignment ?
	/*
		if this.cStruct.strategy == ssGreedy {
			if this.cStruct.beam_search.n_past != -1 ||
				this.cStruct.beam_search.beam_width != -1 ||
				this.cStruct.beam_search.n_best != -1 {
				return false
			}

		} else if this.cStruct.strategy == ssBeamSearch {
			if this.cStruct.greedy.n_past != -1 ||
				this.cStruct.beam_search.beam_width != 10 ||
				this.cStruct.beam_search.n_best != 5 {
				return false
			}
		}
	*/

	return true
}

type _FullParams struct {
	strategy       eSamplingStrategy
	cpuThreads     int32
	n_max_text_ctx int32
	offset_ms      int32
	duration_ms    int32
	Flags          eFullParamsFlags
	Language       eLanguage

	thold_pt    float32
	thold_ptsum float32
	max_len     int32
	max_tokens  int32

	greedy      struct{ n_past int32 }
	beam_search struct {
		n_past     int32
		beam_width int32
		n_best     int32
	}

	audio_ctx int32 // overwrite the audio context size (0 = use default)

	prompt_tokens   uintptr
	prompt_n_tokens int32

	new_segment_callback           uintptr
	new_segment_callback_user_data uintptr

	encoder_begin_callback           uintptr
	encoder_begin_callback_user_data uintptr

	// Are these needed ?? Jay
	// setFlag uintptr
}

func NewFullParams(cstruct *_FullParams) *FullParams {
	this := FullParams{}
	this.cStruct = cstruct
	return &this
}

func _newFullParams_cStruct() *_FullParams {
	return &_FullParams{

		strategy:       0,
		cpuThreads:     0,
		n_max_text_ctx: 0,
		offset_ms:      0,
		duration_ms:    0,

		Flags:    0,
		Language: 0,

		thold_pt:    0,
		thold_ptsum: 0,
		max_len:     0,
		max_tokens:  0,

		// anonymous int32
		greedy: struct{ n_past int32 }{n_past: 0},

		// anonymous struct
		beam_search: struct {
			n_past     int32
			beam_width int32
			n_best     int32
		}{
			n_past:     0,
			beam_width: 0,
			n_best:     0,
		},

		audio_ctx: 0,

		prompt_tokens:   0,
		prompt_n_tokens: 0,

		new_segment_callback:           0,
		new_segment_callback_user_data: 0,

		encoder_begin_callback:           0,
		encoder_begin_callback_user_data: 0,
	}
}
