import { useEffect, useState } from "react";
import Story from "../types/Story";
import NavigateButton from "./Button/NavigateButton";
import FetchStory from "../functions/Stories";
import { useTranslation } from "react-i18next";

interface props {
  id: string;
  detail: boolean;
}

const StoryDetailHeader = ({ id, detail }: props) => {
  const [story, setStory] = useState<Story>();
  const { t } = useTranslation(['home', 'main']);

  useEffect(() => {
    FetchStory(id, setStory);
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
          {story?.announcement || t("common.announce-not-found", {ns: ['main', 'home']})}
        </p>
        <p className="lead mb-0">{story?.notes || t("common.notes-not-found", {ns: ['main', 'home']})}</p>
        <br />
        <NavigateButton link="/stories" variant="secondary">
        {t("story.back-button", {ns: ['main', 'home']})}
        </NavigateButton>{" "}
        {detail === true ? (
          <>
            <NavigateButton
              link={`/stories/encounter/new/${id}`}
              variant="primary"
            >
              {t("encounter.button", {ns: ['main', 'home']})}
            </NavigateButton>{" "}
          </>
        ) : (
          <>
            <NavigateButton variant="primary" link={`/stories/${id}`}>
              {t("story.detail", {ns: ['main', 'home']})}
            </NavigateButton>{" "}
            <br />
          </>
        )}
      </div>
    </div>
  );
};

export default StoryDetailHeader;
