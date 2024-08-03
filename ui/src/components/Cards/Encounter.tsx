import Encounter from "../../types/Encounter";
import NavigateButton from "../Button/NavigateButton";

interface props {
  encounter: Encounter;
  disable_footer: boolean;
}

const EncounterCards = ({ encounter, disable_footer }: props) => {
  return (
    <>
      <div className="col-md-6">
        <div className="card mb-4">
          <div className="card-header">Encounter: {encounter.title} </div>
          <div className="card-body">
            <h6 className="card-title">Announcement</h6>
            <p className="card-text">{encounter.announcement}</p>
            <h6 className="card-title">Notes</h6>
            <p className="card-text">{encounter.notes}</p>
            {/* <a href="#" className="btn btn-primary">
              Go somewhere
            </a> */}
          </div>
          <div className="card-footer text-body-secondary" hidden={disable_footer} >
          <NavigateButton link={`/stories/${encounter.story_id}/encounter/${encounter.id}`} variant="primary">
            Add to Stage
          </NavigateButton>{" "}
          </div>
        </div>
      </div>
    </>
  );
};

export default EncounterCards;
