package test

import (
	"testing"
	"voice-helper-go/api"
)

func TestStt(t *testing.T) {
	api.SttTest()
}

func TestAi(t *testing.T) {
	api.AiTest()
}

func TestTts(t *testing.T) {
	api.TtsTest()
}
