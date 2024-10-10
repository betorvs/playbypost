package validator

import (
	"encoding/json"
	"testing"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

var (
	simpleChoices = []byte(`[
  {
    "id": 1,
    "upstream_id": 1,
    "encounter_id": 7,
    "next_encounter_id": 8,
    "text": "If you want A",
    "objective": {
      "id": 0,
      "kind": "no_action",
      "values": null
    }
  },
  {
    "id": 2,
    "upstream_id": 1,
    "encounter_id": 7,
    "next_encounter_id": 9,
    "text": "If you want B",
    "objective": {
      "id": 0,
      "kind": "no_action",
      "values": null
    }
  },
  {
    "id": 3,
    "upstream_id": 1,
    "encounter_id": 8,
    "next_encounter_id": 10,
    "text": "moving forward",
    "objective": {
      "id": 0,
      "kind": "no_action",
      "values": null
    }
  },
  {
    "id": 4,
    "upstream_id": 1,
    "encounter_id": 9,
    "next_encounter_id": 11,
    "text": "moving forward",
    "objective": {
      "id": 0,
      "kind": "no_action",
      "values": null
    }
  },
  {
    "id": 5,
    "upstream_id": 1,
    "encounter_id": 10,
    "next_encounter_id": 12,
    "text": "go to end notes",
    "objective": {
      "id": 0,
      "kind": "no_action",
      "values": null
    }
  }
]`)
	withValues = []byte(`[
  {
    "id": 6,
    "upstream_id": 2,
    "encounter_id": 13,
    "next_encounter_id": 14,
    "text": "you rolled 1,3 or 5",
    "objective": {
      "id": 0,
      "kind": "dice_roll",
      "values": null
    }
  },
  {
    "id": 7,
    "upstream_id": 2,
    "encounter_id": 13,
    "next_encounter_id": 15,
    "text": "you rolled 2,4 or 6",
    "objective": {
      "id": 0,
      "kind": "dice_roll",
      "values": null
    }
  }
]`)
	withValuesBroken = []byte(`[
  {
    "id": 6,
    "upstream_id": 2,
    "encounter_id": 13,
    "next_encounter_id": 14,
    "text": "you rolled 1,3 or 5",
    "objective": {
      "id": 0,
      "kind": "dice_roll",
      "values": null
    }
  },
  {
    "id": 7,
    "upstream_id": 2,
    "encounter_id": 13,
    "next_encounter_id": 15,
    "text": "you rolled 2,4 or 6",
    "objective": {
      "id": 0,
      "kind": "dice_roll",
      "values": null
    }
  },
  {
    "id": 8,
    "upstream_id": 2,
    "encounter_id": 13,
    "next_encounter_id": 16,
    "text": "good choice",
    "objective": {
      "id": 0,
      "kind": "no_action",
      "values": null
    }
  }
]`)
)

func TestParserAutoPlayNext(t *testing.T) {
	table := []struct {
		name          string
		data          []byte
		first         int
		totalPaths    int
		lastEncounter []int
		wantErr       bool
	}{
		{"simpleChoices", simpleChoices, 7, 2, []int{11, 12}, false},
		{"withValues", withValues, 13, 2, []int{14, 15}, false},
		{"withValuesBroken", withValuesBroken, 13, 2, []int{14, 15}, true},
	}
	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			var next []types.Next
			err := json.Unmarshal(tt.data, &next)
			if err != nil {
				t.Fatal(err)
			}
			result := parserAutoPlayNext(next, tt.first, tt.lastEncounter)
			if (result.TotalPaths != tt.totalPaths) != tt.wantErr {
				t.Errorf("expected %d paths, got %d paths", tt.totalPaths, result.TotalPaths)
			}
			if result.LastEncounters == tt.wantErr {
				t.Errorf("expected %d last encounters, got %d from paths", len(tt.lastEncounter), result.TotalPaths)
			}
			if result.LastEncountersUsed == tt.wantErr {
				t.Errorf("expected last encounters %v, got %v", tt.lastEncounter, result.LastEncounterFound)
			}
			if result.ObjectivesMatch == tt.wantErr {
				t.Errorf("expected objectives does not match")
			}
		})
	}

}
