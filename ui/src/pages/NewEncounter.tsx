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
import { useTranslation } from "react-i18next";

const NewEncounter = () => {
  const { id } = useParams();
  const { Logoff } = useContext(AuthContext);
  const navigate = useNavigate();
  const [title, setTitle] = useState("");
  const [announce, setAnnouncement] = useState("");
  const [note, setNotes] = useState("");
  const [firstEncounter, setFirstEncounter] = useState(false);
  const [lastEncounter, setLastEncounter] = useState(false);
  const user_id = GetUserID();
  const { t } = useTranslation(['home', 'main']);

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
        writer_id: user_id,
        first_encounter: firstEncounter,
        last_encounter: lastEncounter,
      }),
    });
    if (response.ok) {
      alert(t("alert.encounter", {ns: ['main', 'home']}));
      navigate(`/stories/${safeID}`);
    } else {
      alert(t("alert.encounter-wrong", {ns: ['main', 'home']}));
    }
  }

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>{t("encounter.header", {ns: ['main', 'home']})}</h2>
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
            {t("encounter.form-title-text", {ns: ['main', 'home']})}
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
            {t("encounter.form-announce-text", {ns: ['main', 'home']})}
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
            {t("encounter.form-notes-text", {ns: ['main', 'home']})}
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formFirstEncounter">
            <Form.Label>{t("encounter.first-encounter", {ns: ['main', 'home']})}</Form.Label>
            <Form.Check
              type="checkbox"
              label={t("encounter.first-encounter-label", {ns: ['main', 'home']})}
              checked={firstEncounter}
              onChange={(e) => setFirstEncounter(e.target.checked)}
              />
            <Form.Text className="text-muted">
            {t("encounter.first-encounter-description", {ns: ['main', 'home']})}
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formFirstEncounter">
            <Form.Label>{t("encounter.last-encounter", {ns: ['main', 'home']})}</Form.Label>
            <Form.Check
              type="checkbox"
              label={t("encounter.last-encounter-label", {ns: ['main', 'home']})}
              checked={lastEncounter}
              onChange={(e) => setLastEncounter(e.target.checked)}
              />
            <Form.Text className="text-muted">
            {t("encounter.last-encounter-description", {ns: ['main', 'home']})}
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

export default NewEncounter;
