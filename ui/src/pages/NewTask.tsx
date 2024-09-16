import { useContext, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import { Form } from "react-bootstrap";
import Button from "react-bootstrap/Button";
// import GetUserID from "../context/GetUserID";
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import { useNavigate } from "react-router-dom";
import UseLocation from "../context/UseLocation";

const NewTask = () => {
  const { Logoff } = useContext(AuthContext);
  const [description, setDescription] = useState("");
  const [ability, setAbility] = useState("");
  const [skill, setSkill] = useState("");
  const [target, setTarget] = useState(0);
  // const user_id = GetUserID();
  const navigate = useNavigate();

  const cancelButton = () => {
    navigate("/tasks");
  };

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const apiURL = UseLocation();
    const urlAPI = new URL("api/v1/task", apiURL);
    const response = await fetch(urlAPI, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Username": GetUsername(),
        "X-Access-Token": GetToken(),
      },
      body: JSON.stringify({
        description: description,
        ability: ability,
        skill: skill,
        target: target,
      }),
    });
    if (response.ok) {
      alert("Task created! Have a great session with your friends!");
      navigate("/tasks");
    } else {
      alert("Something goes wrong. No new task for you.");
    }
  }
  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>Create a New Task</h2>
        <hr />
      </div>
      <div className="container mt-3" key="2">
        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3" controlId="formTitle">
            <Form.Label>Description</Form.Label>
            <Form.Control
              type="text"
              placeholder="title"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
            <Form.Text className="text-muted">
              short description about this task and how we should use it
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formAbility">
            <Form.Label>Ability</Form.Label>
            <Form.Control
              type="text"
              placeholder="ability"
              value={ability}
              onChange={(e) => setAbility(e.target.value)}
            />
            <Form.Text className="text-muted">
              Which ability should roll for it
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formSkill">
            <Form.Label>Skill</Form.Label>
            <Form.Control
              type="text"
              placeholder="skill"
              value={skill}
              onChange={(e) => setSkill(e.target.value)}
            />
            <Form.Text className="text-muted">
              Which skill should roll for it
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formTarget">
            <Form.Label>Target</Form.Label>
            <Form.Control
              type="number"
              placeholder="target"
              value={target}
              onChange={(e) => setTarget(Number(e.target.value))}
            />
            <Form.Text className="text-muted">
              What number this roll should achieve
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

export default NewTask;
