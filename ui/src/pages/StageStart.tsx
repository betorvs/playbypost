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

const StageStart = () => {
  const { id } = useParams();
  const { Logoff } = useContext(AuthContext);
  const navigate = useNavigate();

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
      alert("Stage started! Have a great session with your friends!");
      navigate(`/stages/${safeID}/story/${stage?.stage.story_id}`);
    } else {
      alert("Something goes wrong. Stage could not start.");
    }
  }

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>Start Stage</h2>
        <hr />
      </div>
      <div className="container mt-3" key="2">
        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3" controlId="formStage">
            <Form.Label>Stage</Form.Label>
            <Form.Select as="select"
              value={channelID}
              onChange={e => {
                console.log("set e.target.value", e.target.value);
                setChannelID(e.target.value);
              }}>
                <option value="-1">Select one</option>
              {
                channel != null ? (
                  channel.map((channel) => (
                    <option value={channel}>{channel}</option>
                  ))) : (
                    <option>channels not found</option>
                  )
              }
            </Form.Select>
            <Form.Text className="text-muted">
              Choose a Channel to Start your Stage
            </Form.Text>
          </Form.Group>
          <h5>Storyteller ID: {stage?.stage.storyteller_id}</h5>
          <h6>In Users, you selected it using "Add as Storyteller"</h6>
          <h5>Stage ID: {safeID}</h5>
          <h6>After choosinga Storyteller from Users, it is created automatically</h6>
          <Button variant="primary" type="submit">
            Submit
          </Button>{" "}
          <Button variant="secondary" onClick={() => cancelButton()}>
            Cancel
          </Button>{" "}
        </Form>
      </div>
    </>
  );
};

export default StageStart;
