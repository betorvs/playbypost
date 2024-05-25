import { useEffect, useState } from "react";
import Story from "../types/Story";
import NavigateButton from "./Button/NavigateButton";
import FetchStories from "../functions/Stories";

interface props {
  id: string;
  detail: boolean;
}

const StoryDetailHeader = ({ id, detail }: props) => {
  const [story, setStory] = useState<Story>();

  useEffect(() => {
    FetchStories(id, setStory);
  }, []);
  return (
    <div
      className="p-4 p-md-5 mb-4 rounded text-body-emphasis bg-body-secondary"
      key="1"
    >
      <div className="col-lg-6 px-0" key="1">
        <h1 className="display-4 fst-italic">
          {story?.title || "story not found"}
        </h1>
        <p className="lead my-3">
          {story?.announcement || "Announcement not found"}
        </p>
        <p className="lead mb-0">{story?.notes || "Notes not found"}</p>
        <br />
        <NavigateButton link="/stories" variant="secondary">
          Back to Stories
        </NavigateButton>{" "}
        {detail === true ? (
          <>
            <NavigateButton link={`/stories/players/${id}`} variant="primary">
              Players List
            </NavigateButton>{" "}
            <NavigateButton
              link={`/stories/encounter/new/${id}`}
              variant="primary"
            >
              New Encounter
            </NavigateButton>
          </>
        ) : (
          <>
            <NavigateButton variant="primary" link={`/stories/${id}`}>
              Story Detail
            </NavigateButton>{" "}
            <br />
          </>
        )}
      </div>
    </div>
  );
};

export default StoryDetailHeader;
