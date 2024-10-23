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
import { useTranslation } from "react-i18next";

const UserAsStoryteller = () => {
  const { Logoff } = useContext(AuthContext);
  const { id } = useParams();
  const [stories, setStory] = useState<Story[]>([]);
  const [text, setText] = useState("");
  const [storyID, setStoryID] = useState(0);
  const user_id = GetUserID();
  const { t } = useTranslation(['home', 'main']);

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
      alert(t("alert.stage", {ns: ['main', 'home']}));
      navigate("/users");
    } else {
      alert(t("alert.stage-wrong", {ns: ['main', 'home']}));
    }
  }
  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>{t("user.add-as-storyteller", {ns: ['main', 'home']})}</h2>
        <h3>{t("user.add-as-storyteller-description", {ns: ['main', 'home']})}</h3>
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
                <option value="-1">{t("common.select-one", {ns: ['main', 'home']})}</option>
              {
                stories != null ? (
                  stories.map((story) => (
                    <option value={story.id}>{story.title}</option>
                  ))) : (
                    <option>{t("story.error", {ns: ['main', 'home']})}</option>
                  )
              }
            </Form.Select>
            <Form.Text className="text-muted">
            {t("user.add-as-storyteller-text1", {ns: ['main', 'home']})}
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
            {t("user.add-as-storyteller-text2", {ns: ['main', 'home']})}
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

export default UserAsStoryteller;
