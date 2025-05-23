import { useState, useEffect } from "react";
import StageCards from "./Cards/Stage";
import FetchStages from "../functions/Stages";
import Stage from "../types/Stage";
import { useTranslation } from "react-i18next";

const StageList = () => {
  const [stages, setStage] = useState<Stage[]>([]);
  const { t } = useTranslation(['home', 'main']);

  useEffect(() => {
    FetchStages(setStage);
  }, []);
  return (
    <div className="container mt-3" key="2">
      {stages.length !== 0 ? (
        stages.map((stage) => (
          <StageCards
            key={stage.id}
            ID={stage.id}
            stage={stage}
            creator_id={stage.creator_id}
          />
        ))
      ) : (
        <p>{t("stage.not-found", {ns: ['main', 'home']})}</p>
      )}
    </div>
  );
};

export default StageList;
