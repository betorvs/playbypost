import { Button } from "react-bootstrap";
import Encounter from "../../types/Encounter";
import NavigateButton from "../Button/NavigateButton";
import { useTranslation } from "react-i18next";
import { DeleteEncounterByID } from "../../functions/Encounters";

interface props {
  encounter: Encounter;
  disable_footer: boolean;
}

const EncounterCards = ({ encounter, disable_footer }: props) => {
  const { t } = useTranslation(['home', 'main']);

  const handleDelete = (id: number) => {
    console.log("Deleting encounter " + id);
    DeleteEncounterByID(id);
  }

  return (
    <>
      <div className="col-md-6">
        <div className="card mb-4">
          <div className="card-header">{t("encounter.this", {ns: ['main', 'home']})}: {encounter.title} ({encounter.id})</div>
          <div className="card-body">
            <h6 className="card-title">{t("common.announce", {ns: ['main', 'home']})}</h6>
            <p className="card-text">{encounter.announcement}</p>
            <h6 className="card-title">{t("common.notes", {ns: ['main', 'home']})}</h6>
            <p className="card-text">{encounter.notes}</p>
            {/* <a href="#" className="btn btn-primary">
              Go somewhere
            </a> */}
          </div>
          <div className="card-footer text-body-secondary" hidden={disable_footer} >
          <NavigateButton link={`/stories/${encounter.story_id}/encounter/${encounter.id}`} variant="primary">
          {t("encounter.add-to-stage", {ns: ['main', 'home']})}
          </NavigateButton>{" "}
          <NavigateButton link={`/stories/${encounter.story_id}/encounter/edit/${encounter.id}`} variant="warning">
          {t("common.edit", {ns: ['main', 'home']})}
          </NavigateButton>{" "}
          <Button variant="danger" size="sm" onClick={() => handleDelete(encounter.id)}>{t("common.delete", {ns: ['main', 'home']})}</Button>
          </div>
        </div>
      </div>
    </>
  );
};

export default EncounterCards;
