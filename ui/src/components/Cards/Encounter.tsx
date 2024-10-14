import Encounter from "../../types/Encounter";
import NavigateButton from "../Button/NavigateButton";
import { useTranslation } from "react-i18next";

interface props {
  encounter: Encounter;
  disable_footer: boolean;
}

const EncounterCards = ({ encounter, disable_footer }: props) => {
  const { t } = useTranslation(['home', 'main']);
  return (
    <>
      <div className="col-md-6">
        <div className="card mb-4">
          <div className="card-header">{t("encounter.this", {ns: ['main', 'home']})}: {encounter.title} </div>
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
          </div>
        </div>
      </div>
    </>
  );
};

export default EncounterCards;
