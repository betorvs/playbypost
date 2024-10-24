import PlayerCards from "./Cards/Player";
import { useEffect, useState } from "react";
import Players from "../types/Players";
import FetchPlayers from "../functions/Players";
import { useTranslation } from "react-i18next";

interface props {
  id: string;
}

const PlayersList = ({ id }: props) => {
  const [players, setPlayer] = useState<Players[]>([]);
  const { t } = useTranslation(['home', 'main']);

  useEffect(() => {
    FetchPlayers(id, setPlayer);
  }, []);

  return (
    <>
      <div className="row mb-2" key="2">
        {players.length !== 0 ? (
          players.map((player, index) => (
            <PlayerCards player={player} key={index} />
          ))
        ) : (
          <p>{t("player.not-found", {ns: ['main', 'home']})}</p>
        )}
      </div>
    </>
  );
};

export default PlayersList;
