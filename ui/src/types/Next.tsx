
interface NextEncounter {
  upstream_id: number;
  encounter_id: number;
  next_encounter_id: number;
  text: string;
  objective: Objective;
}

interface Objective {
    kind: string;
    values: number[];
}

type EncounterWithNext = {
  id: number;
	encounter_id: number;
	next_id: number;
  encounter: string;
  next_encounter: string;
}

type EncounterList = {
  encounter_list: GenericIDName[];
  link: EncounterWithNext[];
  flow_chart_td: string;
}

type GenericIDName = {
  id: number;
  name: string;
}


export default NextEncounter;
export type { EncounterList, EncounterWithNext, GenericIDName };