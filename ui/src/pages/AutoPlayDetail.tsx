import { useContext, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import AutoPlayDetailHeader from "../components/AutoPlayDetailHeader";
import { FetchEncountersAutoPlay } from "../functions/AutoPlay";
import { AutoPlayEncounterList } from "../types/AutoPlay";
import { useTranslation } from "react-i18next";

const AutoPlayDetail = () => {
  const { id, story } = useParams();
  const { Logoff } = useContext(AuthContext);
  const { t } = useTranslation(['home', 'main']);
  
  const safeID: string = id ?? "";
  
  const storySafeID: string = story ?? "";
  
  const [encountersList, setEncountersList] = useState<AutoPlayEncounterList>();

  useEffect(() => {
    FetchEncountersAutoPlay(storySafeID, setEncountersList);
  }, []);
  
    return (
      <>
        <div className="container mt-3" key="1">
          <Layout Logoff={Logoff} />
          {<AutoPlayDetailHeader id={safeID} storyID={storySafeID} />}
        </div>
        <div className="container mt-3" key="3">
          <div className="card" >
            <div className="card-header">
              {t("encounter.list", {ns: ['main', 'home']})}
            </div>
            <ul className="list-group list-group-flush">
              {
                encountersList?.encounter_list != null ? (
                  encountersList?.encounter_list.map((encounter) => (
                    <li className="list-group-item" key={encounter.id} >{encounter.name} ({encounter.id})
                    {
                      encountersList.link != null ? (
                        encountersList.link.filter((next) => next.encounter === encounter.name)
                        .map((next) => (
                          <ul>
                            <li key={next.id}>{t("encounter.link", {ns: ['main', 'home']})}: {next.next_encounter} ({next.next_id}) </li>
                          </ul>
                        ))
                      ) : (
                        <p>{t("encounter.link-not-found", {ns: ['main', 'home']})}</p>
                      )
                    }
                    </li>
                  ))) : (
                  <li className="list-group-item">{t("encounter.not-found", {ns: ['main', 'home']})}</li>
                  )
              }
            </ul>
          </div>
        </div>
      </>
    );
  };
  
  export default AutoPlayDetail;
  