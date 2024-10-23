import { useContext, useEffect, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import { Form } from "react-bootstrap";
import Button from "react-bootstrap/Button";
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import { useNavigate, useParams } from "react-router-dom";
import UseLocation from "../context/UseLocation";
import Stage from "../types/Stage";
import React from "react";
import {FetchStageByStoryID} from "../functions/Stages";
import { useTranslation } from "react-i18next";

const EncounterToStage = () => {
  const { Logoff } = useContext(AuthContext);
  const { story_id, enc_id } = useParams();
  const [stages, setStages] = useState<Stage[]>([]);
  const [text, setText] = useState("");
  const [stageID, setStageID] = useState(0);
  const { t } = useTranslation(["home", "main"]);

  const safeID: string = story_id ?? "";

  useEffect(() => {
    FetchStageByStoryID(safeID, setStages);
  }, []);
  const navigate = useNavigate();

  

  const cancelButton = () => {
    navigate(`/stories/${safeID}`);
  };

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const apiURL = UseLocation();
    const urlAPI = new URL("api/v1/stage/encounter", apiURL);
    const response = await fetch(urlAPI, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Username": GetUsername(),
        "X-Access-Token": GetToken(),
      },
      body: JSON.stringify({
        text: text,
        stage_id: stageID,
        story_id: Number(safeID),
        encounter_id: Number(enc_id),
      }),
    });
    if (response.ok) {
      alert(t("alert.encounter-associated", {ns: ['main', 'home']}));
      navigate(`/stories/${safeID}`);
    } else {
      alert(t("alert.encounter-associated-wrong", {ns: ['main', 'home']}));
    }
  }
  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>{t("encounter.encounter-to-stage", {ns: ['main', 'home']})}</h2>
        <h3>{t("encounter.encounter-to-stage-description", {ns: ['main', 'home']})}</h3>
        <hr />
      </div>
      <div className="container mt-3" key="2">
        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3" controlId="formStage">
            <Form.Label>{t("stage.this", {ns: ['main', 'home']})}</Form.Label>
            <Form.Select as="select"
              value={stageID}
              onChange={e => {
                console.log("set e.target.value", e.target.value);
                setStageID(Number(e.target.value));
              }}>
                <option value="-1">{t("common.select-one", {ns: ['main', 'home']})}</option>
              {
                stages != null ? (
                  stages.map((stage) => (
                    <option value={stage.id}>{stage.text}</option>
                  ))) : (
                    <option>{t("stage.not-found", {ns: ['main', 'home']})}</option>
                  )
              }
            </Form.Select>
            <Form.Text className="text-muted">
            {t("encounter.encounter-to-stage-text1", {ns: ['main', 'home']})}
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formName">
            <Form.Label>{t("common.name", {ns: ['main', 'home']})}</Form.Label>
            <Form.Control
              type="text"
              placeholder="name"
              value={text}
              onChange={(e) => setText(e.target.value)}
            />
            <Form.Text className="text-muted">
            {t("encounter.name", {ns: ['main', 'home']})}
            </Form.Text>
          </Form.Group>
          <Button variant="primary" type="submit">
            {t("common.submit", {ns: ['main', 'home']})}
          </Button>{" "}
          <Button variant="secondary" onClick={() => cancelButton()}>
            {t("common.cancel", {ns: ['main', 'home']})}
          </Button>{" "}
        </Form>
      </div>
    </>
  );
};

export default EncounterToStage;

