import Layout from "../components/Layout";
import { AuthContext } from "../context/AuthContext";
import { useContext, useEffect, useState } from "react";
import StageEncounterDetailHeader from "../components/StageEncounterDetailHeader";
import { useParams } from "react-router-dom";
import Encounter from "../types/Encounter";
import { FetchStageEncounterByEncounterID } from "../functions/Stages";
import Activities from "../types/Activities";
import FetchActivities from "../functions/Activities";
import { useTranslation } from "react-i18next";

const StageEncounterDetail = () => {
    const { Logoff } = useContext(AuthContext);
    const { id, story, encounterid } = useParams();
    const safeID: string = id ?? "";
    const storySafeID: string = story ?? "";
    const encSafeID: string = encounterid ?? "";
    const { t } = useTranslation(['home', 'main']);

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
                {t("player.this", {ns: ['main', 'home']})}
              </div>
              <ul className="list-group list-group-flush">
                {
                  encounter?.pc != null ? (
                    encounter.pc.map((pc) => (
                      <li className="list-group-item" key={pc.id} >{pc.name}</li>
                    ))) : (
                      <li className="list-group-item">{t("player.not-found", {ns: ['main', 'home']})}</li>
                    )
                }
              </ul>
            </div>
          </div>
          <div className="col-md-6">
            <div className="card mb-4">
              <div className="card-header">
              {t("player.npc", {ns: ['main', 'home']})}
              </div>
              <ul className="list-group list-group-flush">
                {
                  encounter?.npc != null ? (
                    encounter.npc.map((npc) => (
                      <li className="list-group-item" key={npc.id}>{npc.name}</li>
                    ))) : (
                      <li className="list-group-item">{t("player.npc-not-found", {ns: ['main', 'home']})}</li>
                    )
                }
              </ul>
            </div>
          </div>
        </div>
        <div className="container mt-3" key="3">
          <div className="card" >
            <div className="card-header">
              {t("stage.encounter-last-activities", {ns: ['main', 'home']})}
            </div>
            <ul className="list-group list-group-flush">
              {
                activities != null ? (
                  activities.map((activity) => (
                    <li className="list-group-item" key={activity.id} >{JSON.stringify(activity.actions)}</li>
                  ))) : (
                  <li className="list-group-item">{t("stage.encounter-last-activities-not-found", {ns: ['main', 'home']})}</li>
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