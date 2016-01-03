// Shared code between blackguns
package common

import "log"

const SampleLen = 15

func TextToSample(txt string) []float64 {
	bs := []byte(txt)
	if len(bs) >= SampleLen {
		bs = bs[:SampleLen]
	}

	sample := make([]float64, SampleLen)

	for i, b := range bs {
		sample[i] = float64(b)
	}

	log.Println(sample)

	return sample
}
