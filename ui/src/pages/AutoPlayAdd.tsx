import { useContext, useEffect, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import { Form } from "react-bootstrap";
import Button from "react-bootstrap/Button";
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import { useNavigate } from "react-router-dom";
import UseLocation from "../context/UseLocation";
import Story from "../types/Story";
import GetUserID from "../context/GetUserID";
import { FetchStoriesByUserID } from "../functions/Stories";

const AutoPlayAdd = () => {
  const { Logoff } = useContext(AuthContext);
  const [stories, setStory] = useState<Story[]>([]);
  const [storyID, setStoryID] = useState(0);
  const [text, setText] = useState("");
  const [solo, setSolo] = useState(false);
  const user_id = GetUserID();
  useEffect(() => {
    FetchStoriesByUserID(user_id, setStory);
  }, []);
  const navigate = useNavigate();

  const cancelButton = () => {
    navigate("/autoplay");
  };

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const apiURL = UseLocation();
    const urlAPI = new URL("api/v1/autoplay", apiURL);
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
        solo: solo,
      }),
    });
    if (response.ok) {
      alert("Auto Play created! Now you need to organise the encounters.");
      navigate("/autoplay");
    } else {
      alert("Something goes wrong. No new stage for you.");
    }
  }
  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>Make Story as Auto Play</h2>
        <h3>It will allow you to organise encounters sequence and make it as Solo Adventure (or Didatic)</h3>
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
              A great Solo Adventures starts with a great title
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formText">
            <Form.Label>Text</Form.Label>
            <Form.Control
              type="text"
              placeholder="Solo Adventure Title"
              value={text}
              onChange={(e) => setText(e.target.value)}
            />
            <Form.Text className="text-muted">
              Text to be used in Auto Play options.
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formSolo">
            <Form.Check
              type="checkbox"
              label="Solo Adventure"
              onChange={(e) => setSolo(e.target.checked)}
            />
            <Form.Text className="text-muted">
              Check if it is a Solo Adventure (or Didatic)
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

export default AutoPlayAdd;
