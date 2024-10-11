import { useContext } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import StageList from "../components/StageList";
import { useTranslation } from "react-i18next";

const StagesPage = () => {
  const { Logoff } = useContext(AuthContext);
  const { t } = useTranslation(['home', 'main']);
  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>{t("stage.this", {ns: ['main', 'home']})}</h2>
        <hr />
      </div>
      {<StageList />}
    </>
  );
};

export default StagesPage;
