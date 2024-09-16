import { useParams } from "react-router-dom";
// import EncounterCards from "../components/Cards/Encounter";
import { useContext, useEffect, useState } from "react";
// import Encounter from "../types/Encounter";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import StageDetailHeader from "../components/StageDetailHeader";
import Encounter from "../types/Encounter";
import StageEncounterCards from "../components/Cards/StageEncounter";
import { FetchStageEncountersByID } from "../functions/Stages";
// import StoryDetailHeader from "../components/StoryDetailHeader";
// import FetchEncounters from "../functions/Encounters";

const StageDetail = () => {
  const { id, story } = useParams();
  const { Logoff } = useContext(AuthContext);

  const safeID: string = id ?? "";

  const storySafeID: string = story ?? "";

  const [encounters, setEncounters] = useState<Encounter[]>([]);

  useEffect(() => {
    FetchStageEncountersByID(safeID, setEncounters);
  }, []);

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        {<StageDetailHeader detail={true} id={safeID} storyID={storySafeID} />}
        <div className="row mb-2" key="2">
          {encounters != null ? (
            encounters.map((encounter, index) => (
              <StageEncounterCards encounter={encounter} key={index} stageID={safeID} storyId={storySafeID} />
            ))
          ) : (
            <p>no encounters for you</p>
          )}
        </div>
      </div>
    </>
  );
};

export default StageDetail;
