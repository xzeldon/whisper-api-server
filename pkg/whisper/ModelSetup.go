package whisper

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// Re-implemented sModelSetup.h

// enum struct eModelImplementation : uint32_t
type eModelImplementation uint32

const (
	// GPGPU implementation based on Direct3D 11.0 compute shaders
	mi_GPU eModelImplementation = 1

	// A hybrid implementation which uses DirectCompute for encode, and decodes on CPU
	// Not implemented in the published builds of the DLL. To enable, change BUILD_HYBRID_VERSION macro to 1
	mi_Hybrid eModelImplementation = 2

	// A reference implementation which uses the original GGML CPU-running code
	// Not implemented in the published builds of the DLL. To enable, change BUILD_BOTH_VERSIONS macro to 1
	mi_Reference eModelImplementation = 3
)

// enum struct eGpuModelFlags : uint32_t
type eGpuModelFlags uint32

const (
	// <summary>Equivalent to <c>Wave32 | NoReshapedMatMul</c> on Intel and nVidia GPUs,<br/>
	// and <c>Wave64 | UseReshapedMatMul</c> on AMD GPUs</summary>
	gmf_None eGpuModelFlags = 0

	// <summary>Use Wave32 version of compute shaders even on AMD GPUs</summary>
	// <remarks>Incompatible with <see cref="Wave64" /></remarks>
	gmf_Wave32 eGpuModelFlags = 1

	// <summary>Use Wave64 version of compute shaders even on nVidia and Intel GPUs</summary>
	// <remarks>Incompatible with <see cref="Wave32" /></remarks>
	gmf_Wave64 eGpuModelFlags = 2

	// <summary>Do not use reshaped matrix multiplication shaders on AMD GPUs</summary>
	// <remarks>Incompatible with <see cref="UseReshapedMatMul" /></remarks>
	gmf_NoReshapedMatMul eGpuModelFlags = 4

	// <summary>Use reshaped matrix multiplication shaders even on nVidia and Intel GPUs</summary>
	// <remarks>Incompatible with <see cref="NoReshapedMatMul" /></remarks>
	gmf_UseReshapedMatMul eGpuModelFlags = 8

	// <summary>Create GPU tensors in a way which allows sharing across D3D devices</summary>
	gmf_Cloneable eGpuModelFlags = 0x10
)

// struct sModelSetup
type sModelSetup struct {
	impl    eModelImplementation
	flags   eGpuModelFlags
	adapter string
}

type _sModelSetup struct {
	impl    eModelImplementation
	flags   eGpuModelFlags
	adapter uintptr
}

func ModelSetup(flags eGpuModelFlags, GPU string) *sModelSetup {
	this := sModelSetup{}
	this.impl = mi_GPU
	this.flags = flags
	this.adapter = GPU

	return &this
}

func (this *sModelSetup) isFlagSet(flag eGpuModelFlags) bool {
	return (this.flags & flag) == 0
}

func (this *sModelSetup) AsCType() *_sModelSetup {
	var err error

	ctype := _sModelSetup{}
	ctype.impl = this.impl
	ctype.flags = this.flags
	ctype.adapter = 0

	// Conver Go String to wchar_t, AKA UTF-16
	if this.adapter != "" {
		var UTF16str *uint16
		UTF16str, err = windows.UTF16PtrFromString(this.adapter)
		ctype.adapter = uintptr(unsafe.Pointer(UTF16str))
	}

	if err != nil {
		return nil
	} else {
		return &ctype
	}
}
