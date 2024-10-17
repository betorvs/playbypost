import { Button } from "react-bootstrap";
import { ValidatorPut } from "../../functions/Validator";
import Story from "../../types/Story";
import NavigateButton from "../Button/NavigateButton";
import { useTranslation } from "react-i18next";

interface Props {
  ID: number;
  story: Story;
  LinkText: string;
}

const StoryCards = ({ ID, story, LinkText }: Props) => {
  const { t } = useTranslation(['home', 'main']);
  const validatorPut = (id: number) => {
    ValidatorPut(id, "stage");
  };
  return (
    <div className="card" key={ID}>
      <div className="card-header">{t("story.this", {ns: ['main', 'home']})} ID: {story.id}</div>
      <div className="card-body" key={ID}>
        <h5 className="card-title">
          <strong>{t("common.title", {ns: ['main', 'home']})}: </strong> {story.title}
        </h5>
        <p className="card-text">
          <strong>{t("common.announce", {ns: ['main', 'home']})}: </strong> {story.announcement}
        </p>
        <NavigateButton link={`/stories/${ID}`} variant="primary">
          {LinkText}
        </NavigateButton>{" "}
        <Button variant="secondary" onClick={() => validatorPut(story.id)}>
          {t("common.validator", {ns: ['main', 'home']})}
        </Button>{" "}
      </div>
      <div className="card-footer text-body-secondary">
      {t("common.writer", {ns: ['main', 'home']})} ID: {story.writer_id}; {t("common.notes", {ns: ['main', 'home']})}: {story.notes}
      </div>
    </div>
  );
};

export default StoryCards;
