import { useState, useEffect } from "react";
import StoryCards from "../components/Cards/Story";
import Story from "../types/Story";
import GetUserID from "../context/GetUserID";
import { FetchStoriesByUserID } from "../functions/Stories";

const StoryList = () => {
  const [stories, setStory] = useState<Story[]>([]);
  const userID = GetUserID();
  useEffect(() => {
    FetchStoriesByUserID(userID, setStory);
  }, []);
  return (
    <div className="container mt-3" key="2">
      {stories != null ? (
        stories.map((story) => (
          <StoryCards
            key={story.id}
            ID={story.id}
            story={story}
            LinkText="Details"
          />
        ))
      ) : (
        <p>no stories for you</p>
      )}
    </div>
  );
};

export default StoryList;
