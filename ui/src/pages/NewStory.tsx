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
import { useTranslation } from "react-i18next";

const NewStory = () => {
  const { Logoff } = useContext(AuthContext);
  const [title, setTitle] = useState("");
  const [announce, setAnnouncement] = useState("");
  const [note, setNotes] = useState("");
  const user_id = GetUserID();
  const navigate = useNavigate();
  const { t } = useTranslation(['home', 'main']);

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
        writer_id: user_id,
      }),
    });
    if (response.ok) {
      alert(t("alert.story", {ns: ['main', 'home']}));
      navigate("/stories");
    } else {
      alert(t("alert.story-wrong", {ns: ['main', 'home']}));
    }
  }
  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>{t("story.header-new", {ns: ['main', 'home']})}</h2>
        <hr />
      </div>
      <div className="container mt-3" key="2">
        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3" controlId="formTitle">
            <Form.Label>{t("common.title", {ns: ['main', 'home']})}</Form.Label>
            <Form.Control
              type="text"
              placeholder="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
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
              placeholder="annoucement"
              value={announce}
              onChange={(e) => setAnnouncement(e.target.value)}
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
              placeholder="notes"
              value={note}
              onChange={(e) => setNotes(e.target.value)}
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

export default NewStory;
