package parser

import (
	"testing"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func TestParserAutoPlaysSolo(t *testing.T) {
	autoPlays := []types.AutoPlay{}
	result := ParserAutoPlays(autoPlays, "start-solo")
	if len(result) != 0 {
		t.Errorf("ParserAutoPlaysSolo failed, expected %v, got %v", 0, len(result))
	}
}

func TestParserAutoPlaysNext(t *testing.T) {
	autoPlays := []types.Next{}
	result, _ := ParserAutoPlaysNext(autoPlays, true)
	if len(result) != 0 {
		t.Errorf("ParserAutoPlaysNext failed, expected %v, got %v", 0, len(result))
	}
}
