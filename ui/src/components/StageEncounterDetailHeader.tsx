import NavigateButton from "./Button/NavigateButton";
import Encounter from "../types/Encounter";
import { useTranslation } from "react-i18next";

interface props {
  id: string;
  stageID: string
  storyID: string;
  encounter: Encounter | undefined;
  detail: boolean;
}

const StageEncounterDetailHeader = ({ stageID, storyID, encounter, detail }: props) => {
  const { t } = useTranslation(['home', 'main']);

  return (
    <div
      className="p-4 p-md-5 mb-4 rounded text-body-emphasis bg-body-secondary"
      key="1"
    >
      <div className="col-lg-6 px-0" key="1">
        <h1 className="display-4 fst-italic">
          {encounter?.text || t("encounter.not-found", {ns: ['main', 'home']})}
        </h1>
        <p className="lead my-3">
          {encounter?.announcement || t("common.announce-not-found", {ns: ['main', 'home']})}
        </p>
        <p className="lead mb-0">Notes: {encounter?.notes || t("common.notes-not-found", {ns: ['main', 'home']})}</p>
        <br />
          {encounter?.phase !== 0 ? (
              encounter?.phase === 1 ? (
                <p>{t("encounter.phase", {ns: ['main', 'home']})}: {t("encounter.phase-started", {ns: ['main', 'home']})}</p>
              ) : (
                encounter?.phase === 2 ? (
                  <p>{t("encounter.phase", {ns: ['main', 'home']})}: {t("encounter.phase-running", {ns: ['main', 'home']})}</p>
                ) : (
                  <p>{t("encounter.phase", {ns: ['main', 'home']})}: {t("encounter.phase-finished", {ns: ['main', 'home']})}</p>
                )
              )
            ) : (
              <p>{t("encounter.phase", {ns: ['main', 'home']})}: {t("encounter.phase-waiting", {ns: ['main', 'home']})}</p>
            )}
        <br />
        <NavigateButton link={`/stages/${stageID}/story/${storyID}`} variant="secondary">
        {t("stage.back-button", {ns: ['main', 'home']})}
        </NavigateButton>{" "}
        {detail === true ? (
          <>
            {/* <NavigateButton link={`/stages/start/${id}`} disabled={stage?.channel.active} variant="primary" >
              Start this Stage
            </NavigateButton>{" "} */}
            <NavigateButton link={`/stages/${stageID}/story/${storyID}/encounter/${encounter?.id}/players`} variant="primary">
            {t("player.add-player", {ns: ['main', 'home']})}
            </NavigateButton>{" "}
            <NavigateButton link={`/stages/${stageID}/story/${storyID}/encounter/${encounter?.id}/task/${encounter?.storyteller_id}`} variant="primary">
            {t("task.assign-task", {ns: ['main', 'home']})}
            </NavigateButton>{" "}
            <NavigateButton link={`/stages/${stageID}/story/${storyID}/encounter/${encounter?.id}/encounter`} variant="primary">
            {t("encounter.add-next-encounter", {ns: ['main', 'home']})}
            </NavigateButton>{" "}
            <NavigateButton link={`/stages/${stageID}/story/${storyID}/encounter/${encounter?.id}/npc/${encounter?.storyteller_id}`} variant="primary">
            {t("player.add-npc", {ns: ['main', 'home']})}
            </NavigateButton>{" "}
          </>
        ) : (
          <>
            <br />
          </>
        )}
      </div>
    </div>
  );
};

export default StageEncounterDetailHeader;
