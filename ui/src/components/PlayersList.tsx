import PlayerCards from "./Cards/Player";
import { useEffect, useState } from "react";
import Players from "../types/Players";
import FetchPlayers from "../functions/Players";

interface props {
  id: string;
}

const PlayersList = ({ id }: props) => {
  const [players, setPlayer] = useState<Players[]>();

  useEffect(() => {
    FetchPlayers(id, setPlayer);
  }, []);

  return (
    <>
      <div className="row mb-2" key="2">
        {players != null ? (
          players.map((player, index) => (
            <PlayerCards player={player} key={index} />
          ))
        ) : (
          <p>no players registered</p>
        )}
      </div>
    </>
  );
};

export default PlayersList;
