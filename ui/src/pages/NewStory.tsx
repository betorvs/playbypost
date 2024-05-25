import { useContext, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import { Form } from "react-bootstrap";
import Button from "react-bootstrap/Button";
import GetUserID from "../context/GetUserID";
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import { useNavigate } from "react-router-dom";
import UseLocation from "../context/UseLocation";

const NewStory = () => {
  const { Logoff } = useContext(AuthContext);
  const [title, setTitle] = useState("");
  const [announce, setAnnouncement] = useState("");
  const [note, setNotes] = useState("");
  const user_id = GetUserID();
  const navigate = useNavigate();

  const cancelButton = () => {
    navigate("/stories");
  };

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const apiURL = UseLocation();
    const urlAPI = new URL("api/v1/story", apiURL);
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
        master_id: user_id,
      }),
    });
    if (response.ok) {
      alert("Story created! Have a great session with your friends!");
      navigate("/stories");
    } else {
      alert("Something goes wrong. No new story for you.");
    }
  }
  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>Create a New Story</h2>
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
              A great history starts with a great title
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

export default NewStory;
