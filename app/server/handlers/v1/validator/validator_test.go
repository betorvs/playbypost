package validator

import (
	"encoding/json"
	"testing"
)

var (
	twoEntries  = []byte(`{"aTTuiiMH":{"id":1,"kind":"story","valid":true,"output":"story validated","checksum":"","updated_at":"2024-10-17T10:23:56.427173-03:00","analise":{"results":null}},"dGfQYAob":{"id":1,"kind":"autoplay","valid":true,"output":"autoplay validated","checksum":"","updated_at":"2024-10-17T10:23:46.818035-03:00","analise":{"results":null}}}`)
	fourEntries = []byte(`{"AiNMALBm":{"id":1,"kind":"autoplay","valid":true,"output":"autoplay validated","checksum":"","updated_at":"2024-10-17T10:24:22.095972-03:00","analise":{"results":null}},"aTTuiiMH":{"id":1,"kind":"story","valid":true,"output":"story validated","checksum":"","updated_at":"2024-10-17T10:23:56.427173-03:00","analise":{"results":null}},"dGfQYAob":{"id":1,"kind":"autoplay","valid":true,"output":"autoplay validated","checksum":"","updated_at":"2024-10-17T10:23:46.818035-03:00","analise":{"results":null}},"dgeXlUGH":{"id":1,"kind":"story","valid":true,"output":"story validated","checksum":"","updated_at":"2024-10-17T10:24:33.573325-03:00","analise":{"results":null}}}`)
)

func TestSlice(t *testing.T) {
	t.Run("Test slice with two entries different kind", func(t *testing.T) {
		result := make(map[string]Request)
		_ = json.Unmarshal(twoEntries, &result)
		v1 := Validator{}
		v1.Request = result

		got := v1.Slice()
		if len(got) != 2 {
			t.Errorf("Expected 2 entries, got %d", len(got))
		}
	})
	t.Run("Test slice with four entries with same id and kind", func(t *testing.T) {
		result := make(map[string]Request)
		_ = json.Unmarshal(fourEntries, &result)
		v1 := Validator{}
		v1.Request = result
		got := v1.Slice()
		if len(got) != 2 {
			t.Errorf("Expected 2 entries, got %d", len(got))
		}
	})
}
