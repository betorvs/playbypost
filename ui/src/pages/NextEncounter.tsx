import { useContext, useEffect, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import { useNavigate, useParams } from "react-router-dom";
import { Form } from "react-bootstrap";
import Layout from "../components/Layout";
import { Button } from "react-bootstrap";
import { FetchStageEncountersByID } from "../functions/Stages";
import Encounter from "../types/Encounter";
import UseLocation from "../context/UseLocation";
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import NextEncounterType from "../types/Next";
import { useTranslation } from "react-i18next";

const NextEncounter = () => {
  const { Logoff } = useContext(AuthContext);
  const { id, story, encounterid } = useParams();
  const { t } = useTranslation(['home', 'main']);
  const stageID: string = id ?? "";
  const safeID: string = encounterid ?? "";
  const encounteridNumber: number = Number(safeID);
  const [encounters, setEncounters] = useState<Encounter[]>([]);

  const [formData, setFormData] = useState<NextEncounterType[]>([{ upstream_id: Number(id), encounter_id: encounteridNumber, next_encounter_id: 0, text: '', objective: { kind: '', values: [] } }]);

  const handleDropdownChange = (index: number, event: React.ChangeEvent<HTMLSelectElement>) => {
    setFormData((prevData) => {
      const newData = [...prevData];
      newData[index].next_encounter_id = Number(event.target.value);
      return newData;
    });
  };

  const handleDropdownObjectiveChange = (index: number, event: React.ChangeEvent<HTMLSelectElement>) => {
    setFormData((prevData) => {
      const newData = [...prevData];
      newData[index].objective = { kind: event.target.value, values: [0] };
      return newData;
    });
  };

  const handleInputChange = (index: number, event: any) => {
    setFormData((prevData) => {
      const newData = [...prevData];
      newData[index].text = event.target.value;
      return newData;
    });
  };

  const addGroup = () => {
    setFormData((prevData) => [...prevData, { upstream_id: Number(id), encounter_id: encounteridNumber, next_encounter_id: 0, text: '', objective: { kind: '', values: [] } }]);
  };

  const removeGroup = (index: number) => {
    setFormData((prevData) => prevData.filter((_, i) => i !== index));
  };

  const navigate = useNavigate();

  useEffect(() => {
    FetchStageEncountersByID(stageID, setEncounters);
  }, []);

  const cancelButton = () => {
    navigate(`/stages/${id}/story/${story}/encounter/${encounterid}`);
  };

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const apiURL = UseLocation();
    const urlAPI = new URL("api/v1/stage/encounter/next", apiURL);
    const response = await fetch(urlAPI, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Username": GetUsername(),
        "X-Access-Token": GetToken(),
      },
      body: JSON.stringify(formData),
    });
    if (response.ok) {
      alert("Next Encounter added! Have a great session with your friends!");
      navigate(`/stages/${id}/story/${story}/encounter/${encounterid}`);
    } else {
      alert("Something goes wrong. Next Encounter was not added to encounter");
    };
  };

  return (
      <>
        <div className="container mt-3" key="1">
            <Layout Logoff={Logoff} />
            <h2>{t("encounter.add-next-encounter", {ns: ['main', 'home']})}</h2>
            <h4>{t("encounter.this", {ns: ['main', 'home']})} ID: {encounteridNumber}</h4>
        </div>
        <div className="container mt-3" key="2">
            <Form onSubmit={handleSubmit}>
            {formData.map((group, index) => (
              <div key={index}>
                <Form.Group>
                  <Form.Label>{t("encounter.next", {ns: ['main', 'home']})}</Form.Label>
                  <Form.Select name="next_encounter_id" value={group.next_encounter_id} onChange={(e) => handleDropdownChange(index, e)}>
                    <option value="-1">{t("stage.select-encounter", {ns: ['main', 'home']})}</option>
                        {
                          encounters != null ? (
                            encounters.filter(encounter => encounter.id !== encounteridNumber)
                              .map((encounter) => (
                                <option value={encounter.id}>{encounter.title}</option>
                            ))) : (
                              <option>{t("encounter.not-found", {ns: ['main', 'home']})}</option>
                            )
                        }
                  </Form.Select>
                  <Form.Text className="text-muted">
                  {t("encounter.next-text1", {ns: ['main', 'home']})}
                  </Form.Text>
                </Form.Group>
                <Form.Group>
                  <Form.Label>{t("common.name-optional", {ns: ['main', 'home']})}</Form.Label>
                  <Form.Control type="text" name="text" value={group.text} onChange={(e) => handleInputChange(index, e)} />
                  <Form.Text className="text-muted">
                  {t("encounter.next-text2", {ns: ['main', 'home']})}
                  </Form.Text>
                </Form.Group>
                <Form.Group>
                    <Form.Label>{t("encounter.next-automatic", {ns: ['main', 'home']})}</Form.Label>
                    <Form.Select name="objective" value={group.objective.kind} onChange={(e) => handleDropdownObjectiveChange(index, e)}>
                      <option value="invalid">{t("encounter.next-objective", {ns: ['main', 'home']})}</option>
                      <option value="no_action">{t("encounter.select-objective-1", {ns: ['main', 'home']})}</option>
                      <option value="dice_roll">{t("encounter.select-objective-2", {ns: ['main', 'home']})}</option>
                      <option value="task_okay">{t("encounter.select-objective-3", {ns: ['main', 'home']})}</option>
                      <option value="victory">{t("encounter.select-objective-4", {ns: ['main', 'home']})}</option>
                    </Form.Select>
                  </Form.Group>
                <Button variant="danger" onClick={() => removeGroup(index)}>{t("encounter.remove", {ns: ['main', 'home']})}</Button>
              </div>

            ))}
                <Button variant="primary" onClick={addGroup}>{t("encounter.add-more", {ns: ['main', 'home']})}</Button>{" "}
                <Button variant="primary" type="submit">{t("common.submit", {ns: ['main', 'home']})}</Button>{" "}
                <Button variant="secondary" onClick={() => cancelButton()}>
                  {t("common.cancel", {ns: ['main', 'home']})}
                </Button>{" "}
            </Form>
        </div>
      </>
  );
};
  
export default NextEncounter;