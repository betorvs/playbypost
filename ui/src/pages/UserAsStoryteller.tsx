import { useContext, useEffect, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import { Form } from "react-bootstrap";
import Button from "react-bootstrap/Button";
import GetUserID from "../context/GetUserID";
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import { useNavigate, useParams } from "react-router-dom";
import UseLocation from "../context/UseLocation";
import Story from "../types/Story";
import { FetchStoriesByUserID } from "../functions/Stories";

const UserAsStoryteller = () => {
  const { Logoff } = useContext(AuthContext);
  const { id } = useParams();
  const [stories, setStory] = useState<Story[]>([]);
  const [text, setText] = useState("");
  const [storyID, setStoryID] = useState(0);
  const user_id = GetUserID();
  useEffect(() => {
    FetchStoriesByUserID(user_id, setStory);
  }, []);
  const navigate = useNavigate();

  const cancelButton = () => {
    navigate("/users");
  };

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const apiURL = UseLocation();
    const urlAPI = new URL("api/v1/stage", apiURL);
    const response = await fetch(urlAPI, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Username": GetUsername(),
        "X-Access-Token": GetToken(),
      },
      body: JSON.stringify({
        story_id: storyID,
        text: text,
        user_id: id,
        storyteller_id: user_id,
      }),
    });
    if (response.ok) {
      alert("Stage created! Have a great session with your friends!");
      navigate("/users");
    } else {
      alert("Something goes wrong. No new stage for you.");
    }
  }
  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>Add as Storyteller</h2>
        <h3>It will create a Stage from the Story of your choice</h3>
        <hr />
      </div>
      <div className="container mt-3" key="2">
        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3" controlId="formStory">
            <Form.Label>Story</Form.Label>
            <Form.Select as="select"
              value={storyID}
              onChange={e => {
                console.log("e.target.value", e.target.value);
                setStoryID(Number(e.target.value));
              }}>
                <option value="-1">Select one</option>
              {
                stories != null ? (
                  stories.map((story) => (
                    <option value={story.id}>{story.title}</option>
                  ))) : (
                    <option>stories not found</option>
                  )
              }
            </Form.Select>
            <Form.Text className="text-muted">
              A great history starts with a great title
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formText">
            <Form.Label>Text</Form.Label>
            <Form.Control
              type="text"
              placeholder="annoucement"
              value={text}
              onChange={(e) => setText(e.target.value)}
            />
            <Form.Text className="text-muted">
              Text to be used in Stage page about this story.
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

export default UserAsStoryteller;
