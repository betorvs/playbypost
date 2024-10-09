package types

import "testing"

func TestValidateNextSlice(t *testing.T) {
	type args struct {
		s            []Next
		upstreamKind string
	}
	tests := []struct {
		name    string
		args    args
		want    []Next
		wantErr bool
	}{
		{
			name: "Test validate next slice",
			args: args{
				s: []Next{
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     1,
						NextEncounterID: 2,
						Text:            "text 2",
						Objective:       Objective{},
					},
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     1,
						NextEncounterID: 3,
						Text:            "text 3",
						Objective:       Objective{},
					},
				},
				upstreamKind: UpstreamKindAutoPlay,
			},
			want: []Next{
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     1,
					NextEncounterID: 2,
					Text:            "text 2",
					Objective: Objective{
						ID:     1,
						Kind:   ObjectiveDefault,
						Values: []int{0},
					},
				},
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     1,
					NextEncounterID: 3,
					Text:            "text 3",
					Objective: Objective{
						ID:     1,
						Kind:   ObjectiveDefault,
						Values: []int{0},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test validate next slice one option",
			args: args{
				s: []Next{
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     1,
						NextEncounterID: 2,
						Text:            "text go forward",
						Objective:       Objective{},
					},
				},
				upstreamKind: UpstreamKindAutoPlay,
			},
			want: []Next{
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     1,
					NextEncounterID: 2,
					Text:            "text go forward",
					Objective: Objective{
						ID:     1,
						Kind:   ObjectiveDefault,
						Values: []int{0},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test validate next slice with dice roll objective",
			args: args{
				s: []Next{
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     1,
						NextEncounterID: 2,
						Text:            "text 2",
						Objective: Objective{
							ID:     1,
							Kind:   ObjectiveDiceRoll,
							Values: []int{0},
						},
					},
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     1,
						NextEncounterID: 3,
						Text:            "text 3",
						Objective: Objective{
							ID:     1,
							Kind:   ObjectiveDiceRoll,
							Values: []int{0},
						},
					},
				},
				upstreamKind: UpstreamKindAutoPlay,
			},
			want: []Next{
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     1,
					NextEncounterID: 2,
					Text:            "text 2",
					Objective: Objective{
						ID:     1,
						Kind:   ObjectiveDiceRoll,
						Values: []int{1, 3, 5},
					},
				},
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     1,
					NextEncounterID: 3,
					Text:            "text 3",
					Objective: Objective{
						ID:     1,
						Kind:   ObjectiveDiceRoll,
						Values: []int{2, 4, 6},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test validate next slice empty",
			args: args{
				s:            []Next{},
				upstreamKind: UpstreamKindAutoPlay,
			},
			want:    []Next{},
			wantErr: true,
		},
		{
			name: "Test validate next slice with next encounter id 0",
			args: args{
				s: []Next{
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     1,
						NextEncounterID: 2,
						Text:            "text 2",
						Objective:       Objective{},
					},
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     1,
						NextEncounterID: 0,
						Text:            "text 0",
						Objective:       Objective{},
					},
				},
				upstreamKind: UpstreamKindAutoPlay,
			},
			want: []Next{
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     1,
					NextEncounterID: 2,
					Text:            "text 2",
					Objective: Objective{
						ID:     1,
						Kind:   ObjectiveDefault,
						Values: []int{0},
					},
				},
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     2,
					NextEncounterID: 0,
					Text:            "text 0",
					Objective: Objective{
						ID:     1,
						Kind:   ObjectiveDefault,
						Values: []int{0},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Test validate next slice invalid objective kind",
			args: args{
				s: []Next{
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     1,
						NextEncounterID: 2,
						Text:            "text 2",
						Objective: Objective{
							ID:     1,
							Kind:   "roll_dice",
							Values: []int{1, 2},
						},
					},
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     1,
						NextEncounterID: 3,
						Text:            "text 3",
						Objective: Objective{
							ID:     1,
							Kind:   "roll_dice",
							Values: []int{3, 4},
						},
					},
				},
				upstreamKind: UpstreamKindStage,
			},
			want: []Next{
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     1,
					NextEncounterID: 2,
					Text:            "text 2",
					Objective: Objective{
						ID:     1,
						Kind:   "roll_dice",
						Values: []int{1, 2},
					},
				},
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     2,
					NextEncounterID: 3,
					Text:            "text 3",
					Objective: Objective{
						ID:     1,
						Kind:   "roll_dice",
						Values: []int{3, 4},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Test validate next slice with dice roll objective but one next option",
			args: args{
				s: []Next{
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     1,
						NextEncounterID: 2,
						Text:            "text 2",
						Objective: Objective{
							ID:     1,
							Kind:   ObjectiveDiceRoll,
							Values: []int{0},
						},
					},
				},
				upstreamKind: UpstreamKindAutoPlay,
			},
			want: []Next{
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     1,
					NextEncounterID: 2,
					Text:            "text 2",
					Objective: Objective{
						ID:     1,
						Kind:   ObjectiveDiceRoll,
						Values: []int{1, 3, 5},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Test validate next slice auto play kind and task okay objective",
			args: args{
				s: []Next{
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     1,
						NextEncounterID: 2,
						Text:            "text task okay",
						Objective: Objective{
							ID:     1,
							Kind:   ObjectiveTaskOkay,
							Values: []int{7},
						},
					},
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     1,
						NextEncounterID: 3,
						Text:            "text task okay",
						Objective: Objective{
							ID:     1,
							Kind:   ObjectiveTaskOkay,
							Values: []int{0},
						},
					},
				},
				upstreamKind: UpstreamKindAutoPlay,
			},
			want: []Next{
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     1,
					NextEncounterID: 2,
					Text:            "text task okay",
					Objective: Objective{
						ID:     1,
						Kind:   ObjectiveTaskOkay,
						Values: []int{7},
					},
				},
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     1,
					NextEncounterID: 3,
					Text:            "text task okay",
					Objective: Objective{
						ID:     1,
						Kind:   ObjectiveTaskOkay,
						Values: []int{0},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Test validate next slice stage kind with one next option and task okay objective",
			args: args{
				s: []Next{
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     1,
						NextEncounterID: 2,
						Text:            "text task okay",
						Objective: Objective{
							ID:     1,
							Kind:   ObjectiveTaskOkay,
							Values: []int{7},
						},
					},
				},
				upstreamKind: UpstreamKindStage,
			},
			want: []Next{
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     1,
					NextEncounterID: 2,
					Text:            "text task okay",
					Objective: Objective{
						ID:     1,
						Kind:   ObjectiveTaskOkay,
						Values: []int{7},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Test validate next slice not match encounter id",
			args: args{
				s: []Next{
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     1,
						NextEncounterID: 2,
						Text:            "text 2",
						Objective:       Objective{},
					},
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     2,
						NextEncounterID: 3,
						Text:            "text 3",
						Objective:       Objective{},
					},
				},
				upstreamKind: UpstreamKindAutoPlay,
			},
			want: []Next{
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     1,
					NextEncounterID: 2,
					Text:            "text 2",
					Objective: Objective{
						ID:     1,
						Kind:   ObjectiveDefault,
						Values: []int{0},
					},
				},
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     2,
					NextEncounterID: 3,
					Text:            "text 3",
					Objective: Objective{
						ID:     1,
						Kind:   ObjectiveDefault,
						Values: []int{0},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Test validate next slice not match next encounter id",
			args: args{
				s: []Next{
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     1,
						NextEncounterID: 3,
						Text:            "text 2",
						Objective:       Objective{},
					},
					{
						ID:              1,
						UpstreamID:      1,
						EncounterID:     2,
						NextEncounterID: 3,
						Text:            "text 3",
						Objective:       Objective{},
					},
				},
				upstreamKind: UpstreamKindAutoPlay,
			},
			want: []Next{
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     1,
					NextEncounterID: 3,
					Text:            "text 2",
					Objective: Objective{
						ID:     1,
						Kind:   ObjectiveDefault,
						Values: []int{0},
					},
				},
				{
					ID:              1,
					UpstreamID:      1,
					EncounterID:     2,
					NextEncounterID: 3,
					Text:            "text 3",
					Objective: Objective{
						ID:     1,
						Kind:   ObjectiveDefault,
						Values: []int{0},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateNextSlice(tt.args.s, tt.args.upstreamKind)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNextSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("ValidateNextSlice() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("ValidateNextSlice() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
