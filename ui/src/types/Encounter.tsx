type Participant = {
  id: number;
  name: string;
}

type Encounter = {
  id: number;
  title: string;
  story_id: number;
  announcement: string;
  notes: string;
  text: number;
  storyteller_id: number;
  writer_id: number;
  first_encounter: boolean;
  pc: Participant[];
  npc: Participant[]
};

export default Encounter;
