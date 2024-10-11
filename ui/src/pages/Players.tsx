import { useContext } from "react";
import Layout from "../components/Layout";
import { AuthContext } from "../context/AuthContext";
import PlayersList from "../components/PlayersList";
import { useParams } from "react-router-dom";
import NavigateButton from "../components/Button/NavigateButton";
import { useTranslation } from "react-i18next";

const PlayersPage = () => {
  const { id, story } = useParams();
  const { Logoff } = useContext(AuthContext);
  const { t } = useTranslation(["home", "main"]);

  const safeID: string = id ?? "";
  const storySafeID: string = story ?? "";

  return (
    <>
      <>
        <div className="container mt-3" key="1">
          <Layout Logoff={Logoff} />
          <h2>{t("player.list", {ns: ['main', 'home']})}</h2>
          <NavigateButton link={`/stages/${safeID}/story/${storySafeID}`} variant="secondary">
            {t("stage.back-button", {ns: ['main', 'home']})}
          </NavigateButton>{" "}
          <hr />
        </div>
        <div className="container mt-3" key="2">
          {<PlayersList id={safeID} />}
        </div>
      </>
    </>
  );
};

export default PlayersPage;
