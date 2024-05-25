type Players = {
  id: number;
  name: string;
  rpg: string;
  abilities: Record<string, string>;
  skills: Record<string, string>;
  destroyed: boolean;
};

export default Players;
