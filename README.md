# Go FAAC

[Go](http://www.golang.org) binding for libfaac. See [AudioCoding.com](http://www.audiocoding.com/faac.html) for more info about FAAC.

## Usage

    sampleRate := 48000
    channelCount := 2

    faacEncoder := faac.Open(sampleRate, channelCount)
	config := faac.EncoderConfiguration{
		BitRate:     48000,
		InputFormat: faac.InputFloat,
	}

	err := faacEncoder.SetConfiguration(&config)
	if err != nil {
       // ...
	}

    encodedBytes = faacEncoder.OutputBuffer()
    var interleavedFloats = []float

    // fill interleavedFloats for floats between +/- 32768.0
    // len(interleavedFloats) must be (less than) faacEncoder.InputSamples()
    // ...

	encodedByteCount := faacEncoder.EncodeFloats(
		interleavedFloats,
		encodedBytes)

	if encodedByteCount > 0 {
		writer.write encoder.encodedBytes[0:encodedByteCount]
	}

	faacEncoder.Close()
# go-faac
