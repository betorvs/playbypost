import { useParams } from "react-router-dom";
import { useContext, useEffect, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import StageDetailHeader from "../components/StageDetailHeader";
import { EncounterList } from "../types/Next";
import { useTranslation } from "react-i18next";
import { DeleteStageNextEncounter, FetchEncounterListStage, FetchStage } from "../functions/Stages";
import { Button, Tab, Tabs } from "react-bootstrap";
import NavigateButton from "../components/Button/NavigateButton";
import StageAggregated from "../types/StageAggregated";
import { MermaidDiagram } from "@lightenna/react-mermaid-diagram";

const StageManageNextEncounter = () => {
  const { id, story } = useParams();
  const { Logoff } = useContext(AuthContext);
  const [stage, setStage] = useState<StageAggregated>();

  const safeID: string = id ?? "";

  const storySafeID: string = story ?? "";

  const { t } = useTranslation(["home", "main"]);

  const [encountersList, setEncountersList] = useState<EncounterList>();

  const [key, setKey] = useState('table');

  const handleDelete = (id: number) => {
    console.log("Deleting next encounter " + id);
    DeleteStageNextEncounter(id);
  }

  useEffect(() => {
    FetchStage(safeID, setStage);
    FetchEncounterListStage(safeID, setEncountersList);
  }, []);

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        {stage && <StageDetailHeader detail={true} disableManageNextEncounter={true} backButtonLink={`/stages/${safeID}/story/${storySafeID}`} stage={stage} />}
        <div className="container mt-3" key="3">
          <div className="container mt-3" key="4">
            <hr/>
            <NavigateButton link={`/stages/${id}/story/${storySafeID}/addnext`} variant="primary">
            {t("encounter.add-next-encounter", {ns: ['main', 'home']})}
            </NavigateButton>{" "}
            <hr/>
          </div>
          <Tabs
          id="controlled-tab"
          activeKey={key}
          onSelect={(k) => k && setKey(k)}
          className="mb-3">

          <Tab eventKey="table" title="Table">
          <div className="card" key={1}>
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
          </Tab>
          <Tab eventKey="flowchart" title="Flowchart">
          <div className="container mt-3" key="2">
          <h3>Flowchart</h3>
            {encountersList?.flow_chart_td && ( 
              <>
              <MermaidDiagram>{encountersList.flow_chart_td}</MermaidDiagram>
              </>
            )}
          </div>
          </Tab>
          </Tabs>
        </div>
      </div>
    </>
  );
};

export default StageManageNextEncounter;
