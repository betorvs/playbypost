import { useContext, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Layout from "../components/Layout";
import { AuthContext } from "../context/AuthContext";
import Button from "react-bootstrap/Button";
import { Form } from "react-bootstrap";
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import GetUserID from "../context/GetUserID";
import UseLocation from "../context/UseLocation";

const NewEncounter = () => {
  const { id } = useParams();
  const { Logoff } = useContext(AuthContext);
  const navigate = useNavigate();
  const [title, setTitle] = useState("");
  const [announce, setAnnouncement] = useState("");
  const [note, setNotes] = useState("");
  const user_id = GetUserID();

  const safeID: string = id ?? "";

  const cancelButton = () => {
    if (safeID === "") {
      navigate("/stories");
    }
    navigate(`/stories/${safeID}`);
  };

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const apiURL = UseLocation();
    const urlAPI = new URL("api/v1/encounter", apiURL);
    const response = await fetch(urlAPI, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Username": GetUsername(),
        "X-Access-Token": GetToken(),
      },
      body: JSON.stringify({
        title: title,
        announcement: announce,
        notes: note,
        story_id: Number(safeID),
        storyteller_id: user_id,
      }),
    });
    if (response.ok) {
      alert("Encounter created! Have a great session with your friends!");
      navigate(`/stories/${safeID}`);
    } else {
      alert("Something goes wrong. No new encounter for you.");
    }
  }

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>Create a New Encounter</h2>
        <hr />
      </div>
      <div className="container mt-3" key="2">
        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3" controlId="formTitle">
            <Form.Label>Title</Form.Label>
            <Form.Control
              type="text"
              placeholder="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
            />
            <Form.Text className="text-muted">
              A great encounter starts with a great title
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formAnnouncement">
            <Form.Label>Announcement</Form.Label>
            <Form.Control
              type="text"
              placeholder="annoucement"
              value={announce}
              onChange={(e) => setAnnouncement(e.target.value)}
            />
            <Form.Text className="text-muted">
              It will be post to all players. Think about all senses and create
              a great description!
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formNotes">
            <Form.Label>Notes</Form.Label>
            <Form.Control
              type="text"
              placeholder="notes"
              value={note}
              onChange={(e) => setNotes(e.target.value)}
            />
            <Form.Text className="text-muted">
              It will be used only for you. Keep here good notes about what
              should happen.
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

export default NewEncounter;
