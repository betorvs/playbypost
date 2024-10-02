package types

type Activity struct {
	ID          int     `json:"id"`
	UpstreamID  int     `json:"upstream_id"`
	Actions     Actions `json:"actions"`
	EncounterID int     `json:"encounter_id"`
	Processed   bool    `json:"processed"`
}
