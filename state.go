package main

import (
	"fmt"
	"sync"

	"github.com/xzeldon/whisper-api-server/pkg/whisper"
)

type WhisperState struct {
	model   *whisper.Model
	context *whisper.IContext
	media   *whisper.IMediaFoundation
	params  *whisper.FullParams
	mutex   sync.Mutex
}

func InitializeWhisperState(modelPath string) (*WhisperState, error) {
	lib, err := whisper.New(whisper.LlDebug, whisper.LfUseStandardError, nil)
	if err != nil {
		return nil, err
	}

	model, err := lib.LoadModel(modelPath)
	if err != nil {
		return nil, err
	}

	context, err := model.CreateContext()
	if err != nil {
		return nil, err
	}

	media, err := lib.InitMediaFoundation()
	if err != nil {
		return nil, err
	}

	params, err := context.FullDefaultParams(whisper.SsBeamSearch)
	if err != nil {
		return nil, err
	}

	params.AddFlags(whisper.FlagNoContext)
	params.AddFlags(whisper.FlagTokenTimestamps)

	fmt.Printf("Params CPU Threads : %d\n", params.CpuThreads())

	return &WhisperState{
		model:   model,
		context: context,
		media:   media,
		params:  params,
	}, nil
}

func getResult(ctx *whisper.IContext) (string, error) {
	results := &whisper.ITranscribeResult{}
	ctx.GetResults(whisper.RfTokens|whisper.RfTimestamps, &results)

	length, err := results.GetSize()
	if err != nil {
		return "", err
	}

	segments := results.GetSegments(length.CountSegments)

	var result string

	for _, seg := range segments {
		result += seg.Text()
	}

	return result, nil
}
