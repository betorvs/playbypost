type Players = {
  id: number;
  name: string;
  rpg: string;
  abilities: Map<string, string>;
  skills: Record<string, string>;
  extensions: Record<string, string>;
  destroyed: boolean;
};

export default Players;
