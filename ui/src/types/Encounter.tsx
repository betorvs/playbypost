type Encounter = {
  id: number;
  title: string;
  story_id: number;
  announcement: string;
  notes: string;
  phase: string;
  finished: boolean;
  reward: string;
  xp: number;
};

export default Encounter;
