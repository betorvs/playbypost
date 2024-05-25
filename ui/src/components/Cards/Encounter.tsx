import Encounter from "../../types/Encounter";

interface props {
  encounter: Encounter;
}

const EncounterCards = ({ encounter }: props) => {
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
        </div>
      </div>
    </>
  );
};

export default EncounterCards;
