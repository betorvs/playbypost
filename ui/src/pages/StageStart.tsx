import { useNavigate, useParams } from "react-router-dom";
import { useContext, useEffect, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import { Button, Form } from "react-bootstrap";
import UseLocation from "../context/UseLocation";
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import FetchChannel from "../functions/Channels";
import { FetchStage } from "../functions/Stages";
import StageAggregated from "../types/StageAggregated";
import { useTranslation } from "react-i18next";

const StageStart = () => {
  const { id } = useParams();
  const { Logoff } = useContext(AuthContext);
  const navigate = useNavigate();
  const { t } = useTranslation(['home', 'main']);

  const safeID: string = id ?? "";

  const [channel, setChannel] = useState<string[]>([]);
  const [channelID, setChannelID] = useState<string>();
  const [stage, setStage] = useState<StageAggregated>();

  useEffect(() => {
    FetchChannel(setChannel);
    FetchStage(safeID, setStage);
  }, []);

  const cancelButton = () => {
    navigate(`/stages/${safeID}/story/${stage?.stage.story_id}`);
  };

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const apiURL = UseLocation();
    const urlAPI = new URL("api/v1/stage/channel", apiURL);
    const response = await fetch(urlAPI, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Username": GetUsername(),
        "X-Access-Token": GetToken(),
      },
      body: JSON.stringify({
        channel: channelID,
        stage_id: Number(safeID),
      }),
    });
    if (response.ok) {
      alert(t("alert.stage-start", {ns: ['main', 'home']}));
      navigate(`/stages/${safeID}/story/${stage?.stage.story_id}`);
    } else {
      const data = await response.text();
      let error = JSON.parse(data);
      console.log(error);
      alert(t("alert.stage-start-wrong", {ns: ['main', 'home']})+ "\n" + error.msg);
    }
  }

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>{t("stage.start", {ns: ['main', 'home']})}</h2>
        <hr />
      </div>
      <div className="container mt-3" key="2">
        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3" controlId="formStage">
            <Form.Label>{t("stage.this", {ns: ['main', 'home']})}</Form.Label>
            <Form.Select as="select"
              value={channelID}
              onChange={e => {
                console.log("set e.target.value", e.target.value);
                setChannelID(e.target.value);
              }}>
                <option value="-1">{t("common.select-one", {ns: ['main', 'home']})}</option>
              {
                channel != null ? (
                  channel.map((channel) => (
                    <option value={channel}>{channel}</option>
                  ))) : (
                    <option>{t("stage.channel-not-found", {ns: ['main', 'home']})}</option>
                  )
              }
            </Form.Select>
            <Form.Text className="text-muted">
              {t("stage.start-text1", {ns: ['main', 'home']})}
            </Form.Text>
          </Form.Group>
          <h5>{t("common.storyteller", {ns: ['main', 'home']})} ID: {stage?.stage.storyteller_id}</h5>
          <h6>{t("stage.start-text2", {ns: ['main', 'home']})} "{t("user.add-as-storyteller", {ns: ['main', 'home']})}"</h6>
          <h5>{t("stage.this", {ns: ['main', 'home']})} ID: {safeID}</h5>
          <h6>{t("stage.start-text3", {ns: ['main', 'home']})}</h6>
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

export default StageStart;
