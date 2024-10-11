import { useContext } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import NavigateButton from "../components/Button/NavigateButton";
import AutoPlayList from "../components/AutoPlayList";
import { useTranslation } from "react-i18next";

const AutoPlayPage = () => {
  const { Logoff } = useContext(AuthContext);
  const { t } = useTranslation(["home", "main"]);

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>{t("auto-play.header", {ns: ['main', 'home']})}</h2>
        <NavigateButton link="/autoplay/new" variant="primary">
          {t("auto-play.add-story", {ns: ['main', 'home']})}
        </NavigateButton>{" "}
        <hr />
      </div>
      {<AutoPlayList />}
    </>
  );
};

export default AutoPlayPage;