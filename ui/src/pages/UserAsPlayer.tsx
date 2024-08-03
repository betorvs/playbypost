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

const UserAsPlayer = () => {
  const { Logoff } = useContext(AuthContext);
  const { id } = useParams();
  const [stages, setStages] = useState<Stage[]>([]);
  const [name, setName] = useState("");
  const [stageID, setStageID] = useState(0);
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
      alert("Player created! Have a great session with your friends!");
      navigate("/users");
    } else {
      alert("Something goes wrong. Player was not create.");
    }
  }
  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>Add as Player</h2>
        <h3>It will create a Player from the Stage of your choice</h3>
        <hr />
      </div>
      <div className="container mt-3" key="2">
        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3" controlId="formStage">
            <Form.Label>Stage</Form.Label>
            <Form.Select as="select"
              value={stageID}
              onChange={e => {
                console.log("set e.target.value", e.target.value);
                setStageID(Number(e.target.value));
              }}>
                <option value="-1">Select one</option>
              {
                stages != null ? (
                  stages.map((stage) => (
                    <option value={stage.id}>{stage.text}</option>
                  ))) : (
                    <option>stages not found</option>
                  )
              }
            </Form.Select>
            <Form.Text className="text-muted">
              Choose a Stage to create a Player
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formName">
            <Form.Label>Name</Form.Label>
            <Form.Control
              type="text"
              placeholder="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
            <Form.Text className="text-muted">
              Character's Name.
            </Form.Text>
          </Form.Group>
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

export default UserAsPlayer;
