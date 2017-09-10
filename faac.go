package faac

/*

#include "faac.h"
#include <stdio.h>

// "int32_t" makes cgo crazy, "signed int" works ...
int wrappedFaacEncEncode(faacEncHandle hEncoder, signed int *inputBuffer, unsigned int samplesInput, unsigned char *outputBuffer, unsigned int bufferSize) {
  return faacEncEncode(hEncoder, inputBuffer, samplesInput, outputBuffer, bufferSize);
}

#cgo LDFLAGS: -lfaac
*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

//Encoder export
type Encoder struct {
	handle C.faacEncHandle

	inputSamples   int
	maxOutputBytes int

	sampleWidth int
}

/*

pConfiguration->shortctl = SHORTCTL_NORMAL;

    pConfiguration->useTns=true;
*/
const (
	InputFloat  = int(C.FAAC_INPUT_FLOAT)
	Input16bits = int(C.FAAC_INPUT_16BIT)
	Input32bits = int(C.FAAC_INPUT_32BIT)
	shortctl    = int(C.SHORTCTL_NORMAL)
	Main        = int(C.MAIN)
)

//Open func
func Open(sampleRate int, channelCount int) *Encoder {
	var inputSamples C.ulong
	var maxOutputBytes C.ulong
	//inputSamples = 960
	handle := C.faacEncOpen(
		C.ulong(sampleRate),
		C.uint(channelCount),
		&inputSamples,
		&maxOutputBytes)

	encoder := &Encoder{
		handle:         handle,
		inputSamples:   int(inputSamples),
		maxOutputBytes: int(maxOutputBytes),
	}

	runtime.SetFinalizer(encoder, finalizeEncoder)

	return encoder
}

type EncoderConfiguration struct {
	// QuantizerQuality int
	BitRate      int
	InputFormat  int
	OutputFormat int
	ObjectType   int
	UseLFE       bool
	UseTNS       bool
}

func (encoder *Encoder) Configuration() *EncoderConfiguration {
	config := C.faacEncGetCurrentConfiguration(encoder.handle)

	return &EncoderConfiguration{
		//	QuantizerQuality: int(config.quantqual),
		BitRate:      int(config.bitRate),
		InputFormat:  int(config.inputFormat),
		OutputFormat: int(config.outputFormat),
		ObjectType:   int(config.aacObjectType),
		UseLFE:       (config.useLfe == 1),
		UseTNS:       (config.useTns == 1),
	}
}

func (encoder *Encoder) SetConfiguration(configuration *EncoderConfiguration) error {
	config := C.faacEncGetCurrentConfiguration(encoder.handle)

	// config.quantqual = C.ulong(configuration.QuantizerQuality)
	config.bitRate = C.ulong(configuration.BitRate)
	config.inputFormat = C.uint(configuration.InputFormat)
	config.outputFormat = C.uint(configuration.OutputFormat)
	config.aacObjectType = C.uint(configuration.ObjectType)

	if configuration.UseLFE {
		config.useLfe = 1
	} else {
		config.useLfe = 0
	}
	if configuration.UseTNS {
		config.useTns = 1
	} else {
		config.useTns = 0
	}

	config.shortctl = C.int(C.SHORTCTL_NORMAL)
	//config.bandWidth = C.uint(10)
	//config.quantqual = C.ulong(64000)

	switch {
	case configuration.InputFormat == Input16bits:
		encoder.sampleWidth = 2
	case configuration.InputFormat == Input32bits:
		encoder.sampleWidth = 4
	}

	if C.faacEncSetConfiguration(encoder.handle, config) == 0 {
		return errors.New("Can't configure Faac encoder")
	}

	return nil
}

func (encoder *Encoder) InputSamples() int {
	return encoder.inputSamples
}

func (encoder *Encoder) MaxOutputBytes() int {
	return encoder.maxOutputBytes
}

//OutputBuffer method
func (encoder *Encoder) OutputBuffer() []byte {
	return make([]byte, encoder.maxOutputBytes)
}

func (encoder *Encoder) EncodeFloats(samples []float32, output []byte) int {
	encodedByteCount := C.wrappedFaacEncEncode(encoder.handle,
		(*C.int)(unsafe.Pointer(&samples[0])),
		C.uint(len(samples)),
		(*C.uchar)(unsafe.Pointer(&output[0])),
		C.uint(len(output)))
	return int(encodedByteCount)
}

func (encoder *Encoder) EncodeBytes(samples []byte, output []byte) int {
	//	log.Println(len())
	encodedByteCount := C.wrappedFaacEncEncode(encoder.handle,
		(*C.int)(unsafe.Pointer(&samples[0])),
		C.uint(len(samples)/encoder.sampleWidth),
		(*C.uchar)(unsafe.Pointer(&output[0])),
		C.uint(len(output)))
	return int(encodedByteCount)
}

func (encoder *Encoder) Close() {
	if encoder.handle != nil {
		C.faacEncClose(encoder.handle)
		encoder.handle = nil
	}
}

func finalizeEncoder(encoder *Encoder) {
	encoder.Close()
}
