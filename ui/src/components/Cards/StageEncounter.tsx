import Encounter from "../../types/Encounter";
import NavigateButton from "../Button/NavigateButton";

interface props {
  encounter: Encounter;
  disable_footer?: boolean;
  stageID: string;
  storyId: string;
}

const StageEncounterCards = ({ encounter, disable_footer, stageID, storyId }: props) => {
  return (
    <>
      <div className="col-md-6">
        <div className="card mb-4">
          <div className="card-header">Text: {encounter.text} </div>
          <div className="card-body">
            <h5 className="card-title">Encounter Title</h5>
            <p className="card-text">{encounter.title}</p>
            <h6 className="card-title">Announcement</h6>
            <p className="card-text">{encounter.announcement}</p>
            <h6 className="card-title">Notes</h6>
            <p className="card-text">{encounter.notes}</p>
          </div>
          <div className="card-footer text-body-secondary" hidden={disable_footer} >
          <NavigateButton link={`/stages/${stageID}/story/${storyId}/encounter/${encounter.id}`} variant="primary">
            Manage Encounter
          </NavigateButton>{" "}
          </div>
        </div>
      </div>
    </>
  );
};

export default StageEncounterCards;
