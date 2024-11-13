import { useParams } from "react-router-dom";
import EncounterCards from "../components/Cards/Encounter";
import { useContext, useEffect, useState } from "react";
import Encounter from "../types/Encounter";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import StoryDetailHeader from "../components/StoryDetailHeader";
import FetchEncounters from "../functions/Encounters";
import { useTranslation } from "react-i18next";

const StoryDetail = () => {
  const { id } = useParams();
  const { Logoff } = useContext(AuthContext);
  const { t } = useTranslation(['home', 'main']);

  const safeID: string = id ?? "";

  const [encounters, setEncounters] = useState<Encounter[]>([]);

  const [cursor, setCursor] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    FetchEncounters(safeID, cursor, encounters, setLoading, setEncounters, setCursor);
  }, []);

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        {<StoryDetailHeader detail={true} id={safeID} />}
        <div className="row mb-2" key="2">
          {encounters.length !== 0 ? (
            encounters.map((encounter, index) => (
              <EncounterCards encounter={encounter} key={index} disable_footer={false} />
            ))
          ) : (
            <p>{t("story.error", {ns: ['main', 'home']})}</p>
          )}
        </div>
        <div className="row mb-2">
          <button className="btn btn-primary" disabled={loading} onClick={() => FetchEncounters(safeID, cursor, encounters, setLoading, setEncounters, setCursor)}>
            {t("story.load_more", {ns: ['main', 'home']})}
          </button>
        </div>
      </div>
    </>
  );
};

export default StoryDetail;
