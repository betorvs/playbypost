type AutoPlay = {
  id: number;
  text: string;
  story_id: number;
  solo: boolean;
};

type AutoPlayEncounterWithNext = {
  encounter: string;
  next_encounter: string;
}

type AutoPlayEncounterList = {
  encounter_list: GenericIDName[];
  link: AutoPlayEncounterWithNext[];
}

type GenericIDName = {
  id: number;
  name: string;
}

// require a nexted map of the next encounters
type NextEncounter = {
  id: number; 
  auto_play_id: number;
  encounter_id: number;
  next_encounter_id: number;
  text: string;
};
  
export default AutoPlay;
export type { AutoPlayEncounterList, AutoPlayEncounterWithNext, GenericIDName, NextEncounter };