import { useParams } from "react-router-dom";
import { useContext, useEffect, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import StageDetailHeader from "../components/StageDetailHeader";
import { EncounterList } from "../types/Next";
import { useTranslation } from "react-i18next";
import { DeleteStageNextEncounter, FetchEncounterListStage } from "../functions/Stages";
import { Button } from "react-bootstrap";
import NavigateButton from "../components/Button/NavigateButton";

const StageManageNextEncounter = () => {
  const { id, story } = useParams();
  const { Logoff } = useContext(AuthContext);

  const safeID: string = id ?? "";

  const storySafeID: string = story ?? "";

  const { t } = useTranslation(["home", "main"]);

  const [encountersList, setEncountersList] = useState<EncounterList>();

  const handleDelete = (id: number) => {
    console.log("Deleting next encounter " + id);
    DeleteStageNextEncounter(id);
  }

  useEffect(() => {
    FetchEncounterListStage(safeID, setEncountersList);
  }, []);

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        {<StageDetailHeader detail={true} id={safeID} storyID={storySafeID} disableManageNextEncounter={true} backButtonLink={`/stages/${safeID}/story/${storySafeID}`} />}
        <div className="container mt-3" key="3">
          <div className="container mt-3" key="4">
            <hr/>
            <NavigateButton link={`/stages/${id}/story/${storySafeID}/addnext`} variant="primary">
            {t("encounter.add-next-encounter", {ns: ['main', 'home']})}
            </NavigateButton>{" "}
            <hr/>
          </div>
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
                          <ul>
                            <li key={next.id}>{t("encounter.link", {ns: ['main', 'home']})}: {next.next_encounter} ({next.next_id}) <Button variant="danger" size="sm" onClick={() => handleDelete(next.id)}>{t("common.delete", {ns: ['main', 'home']})}</Button> </li>
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
      </div>
    </>
  );
};

export default StageManageNextEncounter;
