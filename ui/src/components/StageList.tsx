import { useState, useEffect } from "react";
import StageCards from "./Cards/Stage";
import FetchStages from "../functions/Stages";
import Stage from "../types/Stage";

const StageList = () => {
  const [stages, setStage] = useState<Stage[]>([]);

  useEffect(() => {
    FetchStages(setStage);
  }, []);
  return (
    <div className="container mt-3" key="2">
      {stages != null ? (
        stages.map((stage) => (
          <StageCards
            key={stage.id}
            ID={stage.id}
            stage={stage}
          />
        ))
      ) : (
        <p>no stages for you</p>
      )}
    </div>
  );
};

export default StageList;
