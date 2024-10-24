import { useState, useEffect } from "react";
import StoryCards from "../components/Cards/Story";
import Story from "../types/Story";
import GetUserID from "../context/GetUserID";
import { FetchStoriesByUserID } from "../functions/Stories";
import { useTranslation } from "react-i18next";

const StoryList = () => {
  const [stories, setStory] = useState<Story[]>([]);
  const userID = GetUserID();
  const { t } = useTranslation(['home', 'main']);
  useEffect(() => {
    FetchStoriesByUserID(userID, setStory);
  }, []);
  return (
    <div className="container mt-3" key="2">
      {stories.length !== 0 ? (
        stories.map((story) => (
          <StoryCards
            key={story.id}
            ID={story.id}
            story={story}
            LinkText={t("common.details", {ns: ['main', 'home']})}
          />
        ))
      ) : (
        <p>{t("story.error", {ns: ['main', 'home']})}</p>
      )
      }
    </div>
  );
};

export default StoryList;
