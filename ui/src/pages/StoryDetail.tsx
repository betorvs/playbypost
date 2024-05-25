import { useParams } from "react-router-dom";
import EncounterCards from "../components/Cards/Encounter";
import { useContext, useEffect, useState } from "react";
import Encounter from "../types/Encounter";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import StoryDetailHeader from "../components/StoryDetailHeader";
import FetchEncounters from "../functions/Encounters";

const StoryDetail = () => {
  const { id } = useParams();
  const { Logoff } = useContext(AuthContext);

  const safeID: string = id ?? "";

  const [encounters, setEncounters] = useState<Encounter[]>([]);

  useEffect(() => {
    FetchEncounters(safeID, setEncounters);
  }, []);

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        {<StoryDetailHeader detail={true} id={safeID} />}
        <div className="row mb-2" key="2">
          {encounters != null ? (
            encounters.map((encounter, index) => (
              <EncounterCards encounter={encounter} key={index} />
            ))
          ) : (
            <p>no encounters for you</p>
          )}
        </div>
      </div>
    </>
  );
};

export default StoryDetail;
