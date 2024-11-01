import AutoPlay from "../../types/AutoPlay";
import NavigateButton from "../Button/NavigateButton";
import { useTranslation } from "react-i18next";
import { ValidatorPut } from "../../functions/Validator";
import { Button } from "react-bootstrap";

interface Props {
    ID: number;
    autoPlay: AutoPlay;
  }
  
  const AutoPlayCards = ({ ID, autoPlay }: Props) => {
    const { t } = useTranslation(['home', 'main']);
    const validatorPut = (id: number) => {
      ValidatorPut(id, "autoplay");
    };

    return (
      <div className="card" key={ID}>
        <div className="card-header">{t("auto-play.this", {ns: ['main', 'home']})} ID: {autoPlay.id}</div>
        <div className="card-body" key={ID}>
          <h5 className="card-title">
            <strong>{t("common.title", {ns: ['main', 'home']})}: </strong> {autoPlay.text}
          </h5>
          <p className="card-text">
            <strong>{t("story.this", {ns: ['main', 'home']})} ID: </strong> {autoPlay.story_id}
          </p>
          <NavigateButton link={`/autoplay/${autoPlay.id}/story/${autoPlay.story_id}`} variant="primary">
            {t("common.details", {ns: ['main', 'home']})}
          </NavigateButton>{" "}
          <Button variant="secondary" onClick={() => validatorPut(autoPlay.id)}>
          {t("common.validator", {ns: ['main', 'home']})}
          </Button>{" "}
        </div>
        <div className="card-footer text-body-secondary">
          {
            autoPlay.solo ? (
              <p>{t("auto-play.solo", {ns: ['main', 'home']})}</p>
            ) : (
              <p>{t("auto-play.didatic", {ns: ['main', 'home']})}</p>
            )
          }
          <br />
          {
            autoPlay.publish ? (
              <p>{t("auto-play.publish", {ns: ['main', 'home']})}</p>
            ) : (
              <p>{t("auto-play.not-publish", {ns: ['main', 'home']})}</p>
            )
          }
        </div>
      </div>
    );
  };
  
  export default AutoPlayCards;