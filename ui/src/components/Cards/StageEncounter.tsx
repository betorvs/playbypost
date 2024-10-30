import { Button } from "react-bootstrap";
import Encounter from "../../types/Encounter";
import NavigateButton from "../Button/NavigateButton";
import { useTranslation } from "react-i18next";
import { DeleteStageEncounterByID } from "../../functions/Stages";

interface props {
  encounter: Encounter;
  disable_footer?: boolean;
  stageID: string;
  storyId: string;
}

const StageEncounterCards = ({ encounter, disable_footer, stageID, storyId }: props) => {
  const { t } = useTranslation(['home', 'main']);

  const handleDelete = (id: number) => {
    console.log("Deleting assigned encounter " + id);
    DeleteStageEncounterByID(id);
  }

  return (
    <>
      <div className="col-md-6">
        <div className="card mb-4">
          <div className="card-header">{t("stage.text", {ns: ['main', 'home']})}: {encounter.text}  ({encounter.id})
            {encounter.phase !== 0 ? (
              encounter.phase === 1 ? (
                <span className="badge bg-warning ms-2">{t("encounter.phase-started", {ns: ['main', 'home']})}</span>
              ) : (
                encounter.phase === 2 ? (
                  <span className="badge bg-success ms-2">{t("encounter.phase-running", {ns: ['main', 'home']})}</span>
                ) : (
                  <span className="badge bg-danger ms-2">{t("encounter.phase-finished", {ns: ['main', 'home']})}</span>
                )
              )
            ) : (
              <span className="badge bg-primary ms-2">{t("encounter.phase-waiting", {ns: ['main', 'home']})}</span>
            )}
          </div>
          <div className="card-body">
            <h5 className="card-title">{t("encounter.this", {ns: ['main', 'home']})} - {t("common.title", {ns: ['main', 'home']})}</h5>
            <p className="card-text">{encounter.title}</p>
            <h6 className="card-title">{t("common.announce", {ns: ['main', 'home']})}</h6>
            <p className="card-text">{encounter.announcement}</p>
            <h6 className="card-title">{t("common.notes", {ns: ['main', 'home']})}</h6>
            <p className="card-text">{encounter.notes}</p>
          </div>
          <div className="card-footer text-body-secondary" hidden={disable_footer} >
          <NavigateButton link={`/stages/${stageID}/story/${storyId}/encounter/${encounter.id}`} variant="primary">
            {t("encounter.manage-button", {ns: ['main', 'home']})}
          </NavigateButton>{" "}
          <Button variant="warning" size="sm" onClick={() => handleDelete(encounter.id)}>{t("common.remove", {ns: ['main', 'home']})}</Button>
          </div>
        </div>
      </div>
    </>
  );
};

export default StageEncounterCards;
