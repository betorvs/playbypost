import { useEffect, useState } from "react";
import NavigateButton from "./Button/NavigateButton";
import { CloseStage, FetchStage } from "../functions/Stages";
import StageAggregated from "../types/StageAggregated";
import { useTranslation } from "react-i18next";

interface props {
  id: string;
  storyID: string;
  detail: boolean;
}

const StageDetailHeader = ({ id, storyID, detail }: props) => {
  const [stage, setStage] = useState<StageAggregated | undefined>();
   const { t } = useTranslation(['home', 'main']);

  useEffect(() => {
    FetchStage(id, setStage);
  }, []);
  const handleClose = () => {
    console.log("Close stage");
    CloseStage(id);
  }

  return (
    <div
      className="p-4 p-md-5 mb-4 rounded text-body-emphasis bg-body-secondary"
      key="1"
    >
      <div className="col-lg-6 px-0" key="1">
        <h1 className="display-4 fst-italic">
          {stage?.stage.text || t("stage.not-found", {ns: ['main', 'home']})}
        </h1>
        <p className="lead my-3">
          {stage?.story.announcement || t("common.announce-not-found", {ns: ['main', 'home']})}
        </p>
        <p className="lead mb-0">{t("common.notes", {ns: ['main', 'home']})}: {stage?.story.notes || t("common.notes-not-found", {ns: ['main', 'home']})}</p>
        <br />
        <p className="lead mb-0">{t("stage.running-channel", {ns: ['main', 'home']})}: {stage?.channel.channel || "Stage not started yet"}</p>
        <br />
        <NavigateButton link="/stages" variant="secondary">
        {t("stage.back-button", {ns: ['main', 'home']})}
        </NavigateButton>{" "}
        {detail === true ? (
          <>
            
            <NavigateButton link={`/stages/start/${id}`} disabled={stage?.channel.active} variant="primary" >
            {t("stage.start", {ns: ['main', 'home']})}
            </NavigateButton>{" "}
            <NavigateButton link={`/stages/${id}/story/${storyID}/players`} variant="primary">
            {t("player.list", {ns: ['main', 'home']})}
            </NavigateButton>{" "}

            <span>
              <button className="btn btn-secondary" onClick={handleClose}>
                {t("stage.close", {ns: ['main', 'home']})}
              </button>{" "}
            </span>
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

export default StageDetailHeader;
