import Stage from "../../types/Stage";
import NavigateButton from "../Button/NavigateButton";

interface Props {
  ID: number;
  stage: Stage;
}

const StageCards = ({ ID, stage }: Props) => {
  return (
    <div className="card" key={ID}>
      <div className="card-header">Stage ID: {stage.id}</div>
      <div className="card-body" key={ID}>
        <h5 className="card-title">
          <strong>Title: </strong> {stage.text}
        </h5>
        <p className="card-text">
          <strong>Story ID: </strong> {stage.story_id}
        </p>
        <NavigateButton link={`/stages/${stage.id}/story/${stage.story_id}`} variant="primary">
          Details
        </NavigateButton>{" "}
      </div>
      <div className="card-footer text-body-secondary">
        Storyteller ID: {stage.storyteller_id}
      </div>
    </div>
  );
};

export default StageCards;
