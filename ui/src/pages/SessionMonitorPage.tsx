import { useContext } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import { useTranslation } from "react-i18next";
import SessionMonitoring from "../components/SessionMonitoring";


const SessionMonitorPage = () => {
  const { Logoff } = useContext(AuthContext);
  const { t } = useTranslation(['home', 'main']);

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>{t("common.sessions", {ns: ['main', 'home']})}</h2>
        <hr />
      </div>
      {<SessionMonitoring />}
    </>
  );
};

export default SessionMonitorPage;
