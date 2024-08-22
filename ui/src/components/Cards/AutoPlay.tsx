import AutoPlay from "../../types/AutoPlay";
import NavigateButton from "../Button/NavigateButton";

interface Props {
    ID: number;
    autoPlay: AutoPlay;
  }
  
  const AutoPlayCards = ({ ID, autoPlay }: Props) => {
    return (
      <div className="card" key={ID}>
        <div className="card-header">Auto Play ID: {autoPlay.id}</div>
        <div className="card-body" key={ID}>
          <h5 className="card-title">
            <strong>Title: </strong> {autoPlay.text}
          </h5>
          <p className="card-text">
            <strong>Story ID: </strong> {autoPlay.story_id}
          </p>
          <NavigateButton link={`/autoplay/${autoPlay.id}/story/${autoPlay.story_id}`} variant="primary">
            Details
          </NavigateButton>{" "}
        </div>
        <div className="card-footer text-body-secondary">
          {
            autoPlay.solo ? (
              <p>This is a Solo Adventure</p>
            ) : (
              <p>This is Didatic Adventure</p>
            )
          }
        </div>
      </div>
    );
  };
  
  export default AutoPlayCards;