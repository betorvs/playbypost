import { useEffect, useState } from "react";
import NavigateButton from "./Button/NavigateButton";
import { CloseStage, FetchStage } from "../functions/Stages";
import StageAggregated from "../types/StageAggregated";
import { useTranslation } from "react-i18next";
import Validator from "../types/validator";
import { FetchValidatorByIDKind } from "../functions/Validator";
import { Button } from "react-bootstrap";

interface props {
  id: string;
  storyID: string;
  backButtonLink: string;
  detail: boolean;
  disableManageNextEncounter?: boolean;
}

const StageDetailHeader = ({ id, storyID, detail, disableManageNextEncounter, backButtonLink }: props) => {
  const [stage, setStage] = useState<StageAggregated | undefined>();
  const [validator, setValidator] = useState<Validator>();
  const { t } = useTranslation(['home', 'main']);
  const kind = "stage";

  useEffect(() => {
    FetchStage(id, setStage);
    FetchValidatorByIDKind(Number(id), kind, setValidator);
  }, []);
  const handleClose = (id: number) => {
    if (id === 0) return;
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
        <p>
          { 
            validator != null ? (
              validator.valid === true ? (
                <p>{t("common.validator", {ns: ['main', 'home']})}: {t("common.validator-okay", {ns: ['main', 'home']})}</p>
              ) : (
                <p>{t("common.validator", {ns: ['main', 'home']})}: {validator?.analise?.results || t("common.validator-not-found", {ns: ['main', 'home']}) }</p>
              )
            ) : (
              <p>{t("common.validator-not-found", {ns: ['main', 'home']})}</p>
            )
            
          }
        </p>
        <br />
        <NavigateButton link={backButtonLink} variant="secondary">
        {t("stage.back-button", {ns: ['main', 'home']})}
        </NavigateButton>{" "}
        {detail === true ? (
          <>
            
            <NavigateButton link={`/stages/start/${id}`} disabled={stage?.channel.active} variant="primary" >
            {t("stage.start", {ns: ['main', 'home']})}
            </NavigateButton>{" "}
            <NavigateButton link={`/stages/${id}/story/${storyID}/next`} disabled={disableManageNextEncounter} variant="primary">
            {t("stage.manage-next-encounter", {ns: ['main', 'home']})}
            </NavigateButton>{" "}
            <NavigateButton link={`/stages/${id}/story/${storyID}/players`} variant="primary">
            {t("player.list", {ns: ['main', 'home']})}
            </NavigateButton>{" "}
            

            <span>
              <Button variant="danger" size="sm" onClick={() => handleClose(stage?.stage.id || 0)}>{t("stage.close", {ns: ['main', 'home']})}</Button>
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
