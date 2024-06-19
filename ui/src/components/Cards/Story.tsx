import Story from "../../types/Story";
import NavigateButton from "../Button/NavigateButton";

interface Props {
  ID: number;
  story: Story;
  LinkText: string;
}

const StoryCards = ({ ID, story, LinkText }: Props) => {
  return (
    <div className="card" key={ID}>
      <div className="card-header">Story ID: {story.id}</div>
      <div className="card-body" key={ID}>
        <h5 className="card-title">
          <strong>Title: </strong> {story.title}
        </h5>
        <p className="card-text">
          <strong>Announcement: </strong> {story.announcement}
        </p>
        <NavigateButton link={`/stories/${ID}`} variant="primary">
          {LinkText}
        </NavigateButton>{" "}
        <NavigateButton link={`/stories/players/${ID}`} variant="primary">
          Stage List
        </NavigateButton>{" "}
      </div>
      <div className="card-footer text-body-secondary">
        Storyteller ID: {story.storyteller_id}; Notes: {story.notes}
      </div>
    </div>
  );
};

export default StoryCards;
