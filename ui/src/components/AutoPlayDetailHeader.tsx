// import { useState } from "react";

import { useEffect, useState } from "react";
import NavigateButton from "./Button/NavigateButton";
import AutoPlay from "../types/AutoPlay";
import { FetchAutoPlayByID } from "../functions/AutoPlay";
import { useTranslation } from "react-i18next";
import Validator from "../types/validator";
import { FetchValidatorByIDKind } from "../functions/Validator";

interface props {
    id: string;
    storyID: string;
}

const AutoPlayDetailHeader = ({ id, storyID }: props) => {
  const [autoPlay, setAutoPLay] = useState<AutoPlay>();
  const { t } = useTranslation(['home', 'main']);
  const [validator, setValidator] = useState<Validator>();
  const kind = "autoplay";

  useEffect(() => {
    FetchAutoPlayByID(id, setAutoPLay);
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
          <NavigateButton link={`/autoplay/${id}/story/${storyID}/next`} variant="primary">
            {t("encounter.add-next-encounter", {ns: ['main', 'home']})}
          </NavigateButton>{" "}

        </div>
      </div>
    );
  };
  
  export default AutoPlayDetailHeader;