package example_test

import (
	"context"
	"log"
	"os"

	"pipelined.dev/pipe"
	"pipelined.dev/signal"
	"pipelined.dev/vst2"
	"pipelined.dev/wav"
)

// This example demonstrates how to process .wav file with
// vst2 plugin and write result to a new .wav file.
func Example_2() {
	// open input file.
	inputFile, err := os.Open("_testdata/sample1.wav")
	if err != nil {
		log.Fatalf("failed to open input file: %v", err)
	}
	defer inputFile.Close()

	// open vst2 library.
	vst, err := vst2.Open("_testdata/Krush.vst")
	if err != nil {
		log.Fatalf("failed to open vst2 plugin: %v", err)
	}
	defer vst.Close()

	// create output file.
	outputFile, err := os.Create("_testdata/out2.wav")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// build a pipe with single line.
	p, err := pipe.New(
		&pipe.Line{
			// wav pump.
			Pump: &wav.Pump{
				ReadSeeker: inputFile,
			},
			// vst2 processor.
			Processors: pipe.Processors(
				&vst2.Processor{
					VST: vst,
				},
			),
			// wav sink.
			Sinks: pipe.Sinks(
				&wav.Sink{
					BitDepth:    signal.BitDepth16,
					WriteSeeker: outputFile,
				},
			),
		},
	)
	if err != nil {
		log.Fatalf("failed to bind pipeline: %v", err)
	}
	defer p.Close()

	// run the pipeline.
	err = pipe.Wait(p.Run(context.Background(), 512))
	if err != nil {
		log.Fatalf("failed to execute pipeline: %v", err)
	}
}
