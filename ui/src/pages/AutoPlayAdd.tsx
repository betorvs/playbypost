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
import { useTranslation } from "react-i18next";

const AutoPlayAdd = () => {
  const { Logoff } = useContext(AuthContext);
  const [stories, setStory] = useState<Story[]>([]);
  const [storyID, setStoryID] = useState(0);
  const [text, setText] = useState("");
  const [solo, setSolo] = useState(false);
  const user_id = GetUserID();
  const { t } = useTranslation(["home", "main"]);

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
        <h2>{t("auto-play.add-auto-play-header", {ns: ['main', 'home']})}</h2>
        <h3>{t("auto-play.add-auto-play-description", {ns: ['main', 'home']})}</h3>
        <hr />
      </div>
      <div className="container mt-3" key="2">
        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3" controlId="formStory">
            <Form.Label>{t("story.this", {ns: ['main', 'home']})}</Form.Label>
            <Form.Select as="select"
              value={storyID}
              onChange={e => {
                console.log("e.target.value", e.target.value);
                setStoryID(Number(e.target.value));
              }}>
                <option value="-1">{t("common.select-one", {ns: ['main', 'home']})}</option>
              {
                stories != null ? (
                  stories.map((story) => (
                    <option value={story.id}>{story.title}</option>
                  ))) : (
                    <option>{t("story.not-found", {ns: ['main', 'home']})}</option>
                  )
              }
            </Form.Select>
            <Form.Text className="text-muted">
            {t("auto-play.add-auto-play-text1", {ns: ['main', 'home']})}
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formText">
            <Form.Label>{t("common.title", {ns: ['main', 'home']})}</Form.Label>
            <Form.Control
              type="text"
              placeholder="Solo Adventure Title"
              value={text}
              onChange={(e) => setText(e.target.value)}
            />
            <Form.Text className="text-muted">
            {t("auto-play.add-auto-play-text2", {ns: ['main', 'home']})}
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formSolo">
            <Form.Check
              type="checkbox"
              label="Solo Adventure"
              onChange={(e) => setSolo(e.target.checked)}
            />
            <Form.Text className="text-muted">
            {t("auto-play.add-auto-play-text3", {ns: ['main', 'home']})}
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

export default AutoPlayAdd;
