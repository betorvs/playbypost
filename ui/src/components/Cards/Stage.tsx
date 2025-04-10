import { Button } from "react-bootstrap";
import { ValidatorPut } from "../../functions/Validator";
import Stage from "../../types/Stage";
import NavigateButton from "../Button/NavigateButton";
import { useTranslation } from "react-i18next";
import GetUserID from "../../context/GetUserID";

interface Props {
  ID: number;
  stage: Stage;
  creator_id: number;
}

const StageCards = ({ ID, stage, creator_id }: Props) => {
  const user_id = GetUserID();
  const validatorPut = (id: number) => {
    ValidatorPut(id, "stage");
  };
  let isNotCreator = true;
  if (creator_id === user_id) {
    isNotCreator = false;
  }
  const { t } = useTranslation(['home', 'main']);
  return (
    <div className="card" key={ID}>
      <div className="card-header">{t("stage.this", {ns: ['main', 'home']})} ID: {stage.id}</div>
      <div className="card-body" key={ID}>
        <h5 className="card-title">
          <strong>{t("common.title", {ns: ['main', 'home']})}: </strong> {stage.text}
        </h5>
        <p className="card-text">
          <strong>{t("story.this", {ns: ['main', 'home']})} ID: </strong> {stage.story_id}
        </p>
        <NavigateButton link={`/stages/${stage.id}/story/${stage.story_id}`} variant="primary">
        {t("common.details", {ns: ['main', 'home']})}
        </NavigateButton>{" "}
        <Button variant="secondary" disabled={isNotCreator} onClick={() => validatorPut(stage.id)}>
          {t("common.validator", {ns: ['main', 'home']})}
        </Button>{" "}
      </div>
      <div className="card-footer text-body-secondary">
      {t("common.storyteller", {ns: ['main', 'home']})} ID: {stage.storyteller_id}
      </div>
    </div>
  );
};

export default StageCards;
