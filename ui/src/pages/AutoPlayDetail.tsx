import { useContext, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import AutoPlayDetailHeader from "../components/AutoPlayDetailHeader";
import { DeleteAutoPlayNextEncounter, FetchAutoPlayByID, FetchEncounterListAutoPlay } from "../functions/AutoPlay";
import { EncounterList } from "../types/Next";
import { useTranslation } from "react-i18next";
import { Button } from "react-bootstrap";
import GetUserID from "../context/GetUserID";
import AutoPlay from "../types/AutoPlay";

const AutoPlayDetail = () => {
  const { id, story } = useParams();
  const [autoPlay, setAutoPLay] = useState<AutoPlay>();
  const { Logoff } = useContext(AuthContext);
  const { t } = useTranslation(['home', 'main']);
  const user_id = GetUserID();
  
  const safeID: string = id ?? "";
  
  const storySafeID: string = story ?? "";
  
  const [encountersList, setEncountersList] = useState<EncounterList>();

  const handleDelete = (id: number) => {
    console.log("Deleting next encounter " + id);
    DeleteAutoPlayNextEncounter(id);
  }

  useEffect(() => {
    FetchAutoPlayByID(safeID, setAutoPLay);
    FetchEncounterListAutoPlay(storySafeID, setEncountersList);
  }, []);
  
    return (
      <>
        <div className="container mt-3" key="1">
          <Layout Logoff={Logoff} />
          {autoPlay && <AutoPlayDetailHeader id={safeID} storyID={storySafeID} autoPlay={autoPlay} />}
        </div>
        <div className="container mt-3" key="3">
          <div className="card" >
            <div className="card-header">
              {t("encounter.list", {ns: ['main', 'home']})}
            </div>
            <ul className="list-group list-group-flush" key={1}>
              {
                encountersList?.encounter_list != null ? (
                  encountersList?.encounter_list.map((encounter) => (
                    <li className="list-group-item" key={encounter.id} >{encounter.name} ({encounter.id})
                    {
                      encountersList.link != null ? (
                        encountersList.link.filter((next) => next.encounter === encounter.name)
                        .map((next) => (
                          <ul key={next.id}>
                            <li key={next.id}>{t("encounter.link", {ns: ['main', 'home']})}: {next.next_encounter} ({next.next_id}) 
                              
                              { user_id === autoPlay?.creator_id && (
                                <Button variant="danger" size="sm" onClick={() => handleDelete(next.id)}>{t("common.delete", {ns: ['main', 'home']})}</Button> 
                              )}
                            </li>
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
  