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
import FetchStages from "../functions/Stages";
import { useTranslation } from "react-i18next";

const UserAsPlayer = () => {
  const { Logoff } = useContext(AuthContext);
  const { id } = useParams();
  const [stages, setStages] = useState<Stage[]>([]);
  const [name, setName] = useState("");
  const [stageID, setStageID] = useState(0);
  const { t } = useTranslation(['home', 'main']);

  useEffect(() => {
    FetchStages(setStages);
  }, []);
  const navigate = useNavigate();

  const cancelButton = () => {
    navigate("/users");
  };

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const apiURL = UseLocation();
    const urlAPI = new URL("api/v1/player", apiURL);
    const response = await fetch(urlAPI, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Username": GetUsername(),
        "X-Access-Token": GetToken(),
      },
      body: JSON.stringify({
        stage_id: stageID,
        name: name,
        user_id: id,
      }),
    });
    if (response.ok) {
      alert(t("alert.player", {ns: ['main', 'home']}));
      navigate("/users");
    } else {
      const data = await response.text();
      let error = JSON.parse(data);
      console.log(error);
      alert(t("alert.player-wrong", {ns: ['main', 'home']}) + "\n" + error.msg);
    }
  }
  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>{t("user.add-as-player", {ns: ['main', 'home']})}</h2>
        <h3>{t("user.add-as-player-description", {ns: ['main', 'home']})}</h3>
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
              {t("user.add-as-player-text1", {ns: ['main', 'home']})}
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formName">
            <Form.Label>{t("common.name", {ns: ['main', 'home']})}</Form.Label>
            <Form.Control
              type="text"
              placeholder="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
            <Form.Text className="text-muted">
            {t("user.character", {ns: ['main', 'home']})}
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

export default UserAsPlayer;
