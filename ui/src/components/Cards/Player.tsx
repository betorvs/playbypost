import Players from "../../types/Players";
import { useTranslation } from "react-i18next";

interface props {
  player: Players;
}

const PlayerCards = ({ player }: props) => {
  const { t } = useTranslation(['home', 'main']);
  let abilities = "abilities not found";
  if (player.abilities != null ) {
    abilities = JSON.stringify(player.abilities)
  }
  let skills = "skills not found";
  if (player.skills != null ) {
    skills = JSON.stringify(player.skills)
  }
  let extensions = "extensions not found";
  if (player.extensions != null ) {
    extensions = JSON.stringify(player.extensions)
  }
  return (
    <>
    <div className="col-md-6">
      <div className="card mb-4">
        <div className="card-header">{t("player.this", {ns: ['main', 'home']})}: {player.name} </div>
        <div className="card-body">
          <h6 className="card-title">{t("common.ability", {ns: ['main', 'home']})}</h6>
          <p className="card-text">
            {abilities}
          </p>
          <h6 className="card-title">{t("common.skill", {ns: ['main', 'home']})}</h6>
          <p className="card-text">
            {skills}
          </p>
          <h6 className="card-title">{t("common.other", {ns: ['main', 'home']})}</h6>
          <p className="card-text">
            {extensions}
          </p>
        </div>
      </div>
    </div>
  </>
  );
};

export default PlayerCards;
