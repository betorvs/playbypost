// import { useEffect, useState } from "react";
import NavigateButton from "./Button/NavigateButton";
// import { FetchStageEncounterByEncounterID } from "../functions/Stages";
import Encounter from "../types/Encounter";

interface props {
  id: string;
  stageID: string
  storyID: string;
  encounter: Encounter | undefined;
  detail: boolean;
}

const StageEncounterDetailHeader = ({ stageID, storyID, encounter, detail }: props) => {
  // const [encounter, setEncounter] = useState<Encounter>();

  // useEffect(() => {
  //   FetchStageEncounterByEncounterID(id, setEncounter);
  // }, []);

  return (
    <div
      className="p-4 p-md-5 mb-4 rounded text-body-emphasis bg-body-secondary"
      key="1"
    >
      <div className="col-lg-6 px-0" key="1">
        <h1 className="display-4 fst-italic">
          {encounter?.text || "encounter not found"}
        </h1>
        <p className="lead my-3">
          {encounter?.announcement || "Announcement not found"}
        </p>
        <p className="lead mb-0">Notes: {encounter?.notes || "Notes not found"}</p>
        <br />
        {/* <p className="lead mb-0">Running on channel: {encounter?.channel || "Stage not started yet"}</p> */}
        <br />
        <NavigateButton link={`/stages/${stageID}/story/${storyID}`} variant="secondary">
          Back to Stage
        </NavigateButton>{" "}
        {detail === true ? (
          <>
            {/* <NavigateButton link={`/stages/start/${id}`} disabled={stage?.channel.active} variant="primary" >
              Start this Stage
            </NavigateButton>{" "} */}
            <NavigateButton link={`/stages/${stageID}/story/${storyID}/encounter/${encounter?.id}/players`} variant="primary">
              Add Players
            </NavigateButton>{" "}
            <NavigateButton link={`/stages/${stageID}/story/${storyID}/encounter/${encounter?.id}/task/${encounter?.storyteller_id}`} variant="primary">
              Assign Task
            </NavigateButton>{" "}
            <NavigateButton link={`/stages/${stageID}/story/${storyID}/encounter/${encounter?.id}/encounter`} variant="primary">
              Next Encounter
            </NavigateButton>{" "}
            <NavigateButton link={`/stages/${stageID}/story/${storyID}/encounter/${encounter?.id}/npc/${encounter?.storyteller_id}`} variant="primary">
              Add NPC
            </NavigateButton>{" "}
          </>
        ) : (
          <>
            <br />
          </>
        )}
      </div>
    </div>
  );
};

export default StageEncounterDetailHeader;
