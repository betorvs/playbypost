// import { useState } from "react";

import { useEffect, useState } from "react";
import NavigateButton from "./Button/NavigateButton";
import AutoPlay from "../types/AutoPlay";
import { ChangePublishAutoPlay } from "../functions/AutoPlay";
import { useTranslation } from "react-i18next";
import Validator from "../types/validator";
import { FetchValidatorByIDKind } from "../functions/Validator";
import { Button } from "react-bootstrap";
import GetUserID from "../context/GetUserID";

interface props {
    id: string;
    storyID: string;
    autoPlay: AutoPlay;
}

const AutoPlayDetailHeader = ({ id, storyID, autoPlay }: props) => {
  // const [autoPlay, setAutoPLay] = useState<AutoPlay>();
  const { t } = useTranslation(['home', 'main']);
  const [validator, setValidator] = useState<Validator>();
  const user_id = GetUserID();
  const kind = "autoplay";

  const handlePublish = (id: number, publish: boolean) => {
    console.log("Publish: ", id, publish);
    ChangePublishAutoPlay(id);
  };

  let publishButton = "success";
  let publishText = t("auto-play.publish-button", {ns: ['main', 'home']});
  if (autoPlay?.publish) {
    publishButton = "warning";
    publishText = t("auto-play.unpublish-button", {ns: ['main', 'home']});
  }

  useEffect(() => {
    // FetchAutoPlayByID(id, setAutoPLay);
    FetchValidatorByIDKind(Number(id), kind, setValidator);
  }, []);
    return (
      <div
        className="p-4 p-md-5 mb-4 rounded text-body-emphasis bg-body-secondary"
        key="1"
      >
        <div className="col-lg-6 px-0" key="1">
          <h1 className="display-4 fst-italic">
            { autoPlay?.text || t("common.title-not-found", {ns: ['main', 'home']})}
          </h1>
          <p className="lead my-3">
            ID: { autoPlay?.id || t("auto-play.not-found", {ns: ['main', 'home']})} 
          </p>
          <p className="lead mb-0">{t("story.this", {ns: ['main', 'home']})} ID: { autoPlay?.story_id || t("story.not-found", {ns: ['main', 'home']})}</p>
          <br />
          <p className="lead mb-0">
          {
            autoPlay?.solo ? (
              <p>{t("auto-play.solo", {ns: ['main', 'home']})}</p>
            ) : (
              <p>{t("auto-play.didatic", {ns: ['main', 'home']})}</p>
            )
          }
          </p>
          <p>
          {
            autoPlay?.publish ? (
              <p>{t("auto-play.publish", {ns: ['main', 'home']})}</p>
            ) : (
              <p>{t("auto-play.not-publish", {ns: ['main', 'home']})}</p>
            )
          }
          </p>
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
         
          <NavigateButton link="/autoplay" variant="secondary">
            {t("auto-play.back-button", {ns: ['main', 'home']})}
          </NavigateButton>{" "}

            { autoPlay?.creator_id === user_id && (
              <>
                <NavigateButton link={`/autoplay/${id}/story/${storyID}/next`} variant="primary">
                  {t("encounter.add-next-encounter", {ns: ['main', 'home']})}
                </NavigateButton>{" "}
              
                <span>
                  <Button variant={publishButton} size="sm" onClick={() => handlePublish(autoPlay?.id || 0, false)}>{publishText}</Button>
                </span>
             </>

            )}

         
          

        </div>
      </div>
    );
  };
  
  export default AutoPlayDetailHeader;