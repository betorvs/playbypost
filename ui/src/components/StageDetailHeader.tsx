import { useEffect, useState } from "react";
import NavigateButton from "./Button/NavigateButton";
import { CloseStage, FetchStage } from "../functions/Stages";
import StageAggregated from "../types/StageAggregated";

interface props {
  id: string;
  storyID: string;
  detail: boolean;
}

const StageDetailHeader = ({ id, storyID, detail }: props) => {
  const [stage, setStage] = useState<StageAggregated | undefined>();

  useEffect(() => {
    FetchStage(id, setStage);
  }, []);
  const handleClose = () => {
    console.log("Close stage");
    CloseStage(id);
  }

  return (
    <div
      className="p-4 p-md-5 mb-4 rounded text-body-emphasis bg-body-secondary"
      key="1"
    >
      <div className="col-lg-6 px-0" key="1">
        <h1 className="display-4 fst-italic">
          {stage?.stage.text || "stage not found"}
        </h1>
        <p className="lead my-3">
          {stage?.story.announcement || "Announcement not found"}
        </p>
        <p className="lead mb-0">Notes: {stage?.story.notes || "Notes not found"}</p>
        <br />
        <p className="lead mb-0">Running on channel: {stage?.channel.channel || "Stage not started yet"}</p>
        <br />
        <NavigateButton link="/stages" variant="secondary">
          Back to Stages
        </NavigateButton>{" "}
        {detail === true ? (
          <>
            
            <NavigateButton link={`/stages/start/${id}`} disabled={stage?.channel.active} variant="primary" >
              Start this Stage
            </NavigateButton>{" "}
            <NavigateButton link={`/stages/${id}/story/${storyID}/players`} variant="primary">
              Players List
            </NavigateButton>{" "}

            <span>
              <button className="btn btn-secondary" onClick={handleClose}>
                Close Stage
                </button>{" "}
            </span>
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

export default StageDetailHeader;
