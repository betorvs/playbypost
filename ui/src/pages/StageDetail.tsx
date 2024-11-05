import { useParams } from "react-router-dom";
import { useContext, useEffect, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import StageDetailHeader from "../components/StageDetailHeader";
import Encounter from "../types/Encounter";
import StageEncounterCards from "../components/Cards/StageEncounter";
import { FetchStage, FetchStageEncountersByID } from "../functions/Stages";
import { useTranslation } from "react-i18next";
import StageAggregated from "../types/StageAggregated";

const StageDetail = () => {
  const { id, story } = useParams();
  const { Logoff } = useContext(AuthContext);
  const [stage, setStage] = useState<StageAggregated>();

  const safeID: string = id ?? "";

  const storySafeID: string = story ?? "";

  const { t } = useTranslation(["home", "main"]);

  const [encounters, setEncounters] = useState<Encounter[]>([]);

  useEffect(() => {
    FetchStage(safeID, setStage);
    FetchStageEncountersByID(safeID, setEncounters);
  }, []);

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        {stage && <StageDetailHeader detail={true} backButtonLink="/stages" stage={stage} />}
        <div className="row mb-2" key="2">
          {encounters.length !== 0 ? (
            encounters.map((encounter, index) => (
              <StageEncounterCards encounter={encounter} key={index} stageID={safeID} storyId={storySafeID} creator_id={stage?.stage.creator_id ?? 0} />
            ))
          ) : (
            <p>{t("stage.no-encounter", {ns: ['main', 'home']})}</p>
          )}
        </div>
      </div>
    </>
  );
};

export default StageDetail;
