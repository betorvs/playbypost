import { useContext, useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Layout from "../components/Layout";
import { AuthContext } from "../context/AuthContext";
import Button from "react-bootstrap/Button";
import { Form } from "react-bootstrap";
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import UseLocation from "../context/UseLocation";
import { useTranslation } from "react-i18next";
import Encounter from "../types/Encounter";
import { FetchEncounter } from "../functions/Encounters";

const EditEncounter = () => {
  const { story_id, enc_id } = useParams();
  const { Logoff } = useContext(AuthContext);
  const navigate = useNavigate();
  const [encounter, setEncounter] = useState<Encounter>();
  const { t } = useTranslation(['home', 'main']);

  const storySafeID: string = story_id ?? "";
  const encSafeID: string = enc_id ?? "";

  const cancelButton = () => {
    if (storySafeID === "") {
      navigate("/stories");
    }
    navigate(`/stories/${storySafeID}`);
  };

  const onChangeTitle = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (encounter) {
      setEncounter({ ...encounter, title: e.target.value });
    }
  }

  const onChangeAnnouncement = (e: React.ChangeEvent<HTMLInputElement>) => {  
    if (encounter) {
      setEncounter({ ...encounter, announcement: e.target.value });
    }
  }

  const onChangeNotes = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    if (encounter) {
      setEncounter({ ...encounter, notes: e.target.value });
    }
  }

  const onChangeFirstEncounter = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (encounter) {
      setEncounter({ ...encounter, first_encounter: e.target.checked });
    }
  }
  
  const onChangeLastEncounter = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (encounter) {
      setEncounter({ ...encounter, last_encounter: e.target.checked });
    }
  }

  useEffect(() => {
    FetchEncounter(encSafeID, setEncounter);
  }, []);

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const apiURL = UseLocation();
    const urlAPI = new URL("api/v1/encounter/" + enc_id, apiURL);
    const response = await fetch(urlAPI, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "X-Username": GetUsername(),
        "X-Access-Token": GetToken(),
      },
      body: JSON.stringify({
        id: Number(encSafeID),
        title: encounter?.title,
        announcement: encounter?.announcement,
        notes: encounter?.notes,
        story_id: Number(storySafeID),
        writer_id: encounter?.writer_id,
        first_encounter: encounter?.first_encounter,
        last_encounter: encounter?.last_encounter,
      }),
    });
    if (response.ok) {
      alert("Encounter edited! Have a great session with your friends!");
      navigate(`/stories/${storySafeID}`);
    } else {
      alert("Something goes wrong. No new encounter for you.");
    }
  }

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>{t("encounter.header-edit", {ns: ['main', 'home']})}</h2>
        <hr />
      </div>
      <div className="container mt-3" key="2">
        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3" controlId="formTitle">
            <Form.Label>{t("common.title", {ns: ['main', 'home']})}</Form.Label>
            <Form.Control
              type="text"
              placeholder="title"
              value={encounter?.title}
              onChange={onChangeTitle}
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
              value={encounter?.announcement}
              onChange={onChangeAnnouncement}
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
              value={encounter?.notes}
              onChange={onChangeNotes}
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
              checked={encounter?.first_encounter}
              onChange={onChangeFirstEncounter}
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
              checked={encounter?.last_encounter}
              onChange={onChangeLastEncounter}
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

export default EditEncounter;
