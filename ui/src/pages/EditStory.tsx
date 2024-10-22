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
import { useTranslation } from "react-i18next";
import FetchStory from "../functions/Stories";
import Story from "../types/Story";

const EditStory = () => {
  const { Logoff } = useContext(AuthContext);
  const { id } = useParams();
  const [story, setStory] = useState<Story>();
  const user_id = GetUserID();
  const navigate = useNavigate();
  const { t } = useTranslation(['home', 'main']);

  const safeID: string = id ?? "";

  const cancelButton = () => {
    navigate("/stories");
  };

  const onChangeTitle = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (story) {
      setStory({ ...story, title: e.target.value });
    }
  }

  const onChangeAnnouncement = (e: React.ChangeEvent<HTMLInputElement>) => {
      if (story) {
          setStory({ ...story, announcement: e.target.value });
      }
  }

  const onChangeNotes = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
      if (story) {
          setStory({ ...story, notes: e.target.value });
      }
  }

  useEffect(() => {
    FetchStory(safeID, setStory);
  }, []);

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const apiURL = UseLocation();
    const urlAPI = new URL("api/v1/story/" + id, apiURL);
    const response = await fetch(urlAPI, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "X-Username": GetUsername(),
        "X-Access-Token": GetToken(),
      },
      body: JSON.stringify({
        id: Number(id),
        title: story?.title,
        announcement: story?.announcement,
        notes: story?.notes,
        writer_id: user_id,
      }),
    });
    if (response.ok) {
      alert("Story edited! Have a great session with your friends!");
      navigate("/stories");
    } else {
      alert("Something goes wrong. No new story for you.");
    }
  }
  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>{t("story.header-edit", {ns: ['main', 'home']})}</h2>
        <hr />
      </div>
      <div className="container mt-3" key="2">
        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3" controlId="formTitle">
            <Form.Label>{t("common.title", {ns: ['main', 'home']})}</Form.Label>
            <Form.Control
              type="text"
              placeholder={story?.title}
              value={story?.title}
              onChange={onChangeTitle}
            />
            <Form.Text className="text-muted">
            {t("story.form-title-text", {ns: ['main', 'home']})}
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formAnnouncement">
            <Form.Label>{t("common.announce", {ns: ['main', 'home']})}</Form.Label>
            <Form.Control
              type="text"
              as="textarea"
              placeholder={story?.announcement}
              value={story?.announcement}
              onChange={onChangeAnnouncement}
            />
            <Form.Text className="text-muted">
            {t("story.form-announce-text", {ns: ['main', 'home']})}
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formNotes">
            <Form.Label>{t("common.notes", {ns: ['main', 'home']})}</Form.Label>
            <Form.Control
              type="text"
              as="textarea"
              placeholder={story?.notes}
              value={story?.notes}
              onChange={onChangeNotes}
            />
            <Form.Text className="text-muted">
            {t("story.form-notes-text", {ns: ['main', 'home']})}
            </Form.Text>
          </Form.Group>
          <Button variant="primary" type="submit">
          {t("common.submit", {ns: ['main','home']})}
          </Button>{" "}
          <Button variant="secondary" onClick={() => cancelButton()}>
          {t("common.cancel", {ns: ['main','home']})}
          </Button>{" "}
        </Form>
      </div>
    </>
  );
};

export default EditStory;
