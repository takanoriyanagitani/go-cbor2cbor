package arr2wz2arr

import (
	w0 "github.com/tetratelabs/wazero"
)

type Config struct {
	w0.ModuleConfig

	SetInputSize       string
	EstimateOutputSize string
	SetOutputSize      string
	InputOffset        string
	Convert            string
	OutputOffset       string
}

var ConfigDefault Config = Config{
	ModuleConfig: w0.NewModuleConfig().WithName(""),

	SetInputSize:       "set_input_size",
	EstimateOutputSize: "estimate_output_size",
	SetOutputSize:      "set_output_size",
	InputOffset:        "input_offset",
	Convert:            "convert",
	OutputOffset:       "output_offset",
}
