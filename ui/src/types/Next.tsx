
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

export default NextEncounter;