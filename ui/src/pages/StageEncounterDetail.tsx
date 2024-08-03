import Layout from "../components/Layout";
import { AuthContext } from "../context/AuthContext";
import { useContext, useEffect, useState } from "react";
import StageEncounterDetailHeader from "../components/StageEncounterDetailHeader";
import { useParams } from "react-router-dom";
import Encounter from "../types/Encounter";
import { FetchStageEncounterByEncounterID } from "../functions/Stages";
import Activities from "../types/Activities";
import FetchActivities from "../functions/Activities";

const StageEncounterDetail = () => {
    const { Logoff } = useContext(AuthContext);
    const { id, story, encounterid } = useParams();
    const safeID: string = id ?? "";
    const storySafeID: string = story ?? "";
    const encSafeID: string = encounterid ?? "";

    const [encounter, setEncounter] = useState<Encounter>();
    const [activities, setActivities] = useState<Activities[]>();

    useEffect(() => {
      FetchStageEncounterByEncounterID(encSafeID, setEncounter);
      FetchActivities(encSafeID, setActivities);
    }, []);
    

    return (
        <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        {<StageEncounterDetailHeader id={encSafeID} stageID={safeID} storyID={storySafeID} encounter={encounter} detail={true} />}
        <div className="row mb-2" key="2">
          <div className="col-md-6">
            <div className="card mb-4">
              <div className="card-header">
                Players
              </div>
              <ul className="list-group list-group-flush">
                {
                  encounter?.pc != null ? (
                    encounter.pc.map((pc) => (
                      <li className="list-group-item" key={pc.id} >{pc.name}</li>
                    ))) : (
                      <li className="list-group-item">Players not found</li>
                    )
                }
              </ul>
            </div>
          </div>
          <div className="col-md-6">
            <div className="card mb-4">
              <div className="card-header">
                Non Players Characters
              </div>
              <ul className="list-group list-group-flush">
                {
                  encounter?.npc != null ? (
                    encounter.npc.map((npc) => (
                      <li className="list-group-item" key={npc.id}>{npc.name}</li>
                    ))) : (
                      <li className="list-group-item">NPCs not found</li>
                    )
                }
              </ul>
            </div>
          </div>
        </div>
        <div className="container mt-3" key="3">
          <div className="card" >
            <div className="card-header">
              Encounter last activities
            </div>
            <ul className="list-group list-group-flush">
              {
                activities != null ? (
                  activities.map((activity) => (
                    <li className="list-group-item" key={activity.id} >{JSON.stringify(activity.actions)}</li>
                  ))) : (
                  <li className="list-group-item">Activities not found</li>
                  )
              }
            </ul>
          </div>
        </div>
      </div>
        </>
    );
};
  
export default StageEncounterDetail;