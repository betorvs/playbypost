import Players from "../../types/Players";

interface props {
  player: Players;
}

const PlayerCards = ({ player }: props) => {
  let abilities = "abilities not found";
  if (player.abilities != null ) {
    abilities = JSON.stringify(player.abilities)
  }
  let skills = "skills not found";
  if (player.skills != null ) {
    skills = JSON.stringify(player.skills)
  }
  let extensions = "extensions not found";
  if (player.extensions != null ) {
    extensions = JSON.stringify(player.extensions)
  }
  return (
    <>
    <div className="col-md-6">
      <div className="card mb-4">
        <div className="card-header">Player: {player.name} </div>
        <div className="card-body">
          <h6 className="card-title">Abilities</h6>
          <p className="card-text">
            {abilities}
          </p>
          <h6 className="card-title">Skills</h6>
          <p className="card-text">
            {skills}
          </p>
          <h6 className="card-title">Others</h6>
          <p className="card-text">
            {extensions}
          </p>
          {/* <a href="#" className="btn btn-primary">
            Go somewhere
          </a> */}
        </div>
      </div>
    </div>
  </>
  );
};

export default PlayerCards;
