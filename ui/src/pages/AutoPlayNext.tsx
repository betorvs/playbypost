import { useContext, useEffect, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import { useNavigate, useParams } from "react-router-dom";
import Encounter from "../types/Encounter";
import FetchEncounters from "../functions/Encounters";
import UseLocation from "../context/UseLocation";
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import Layout from "../components/Layout";
import { Button, Form } from "react-bootstrap";
import NextEncounterType from "../types/Next";
import { useTranslation } from "react-i18next";


const AutoPlayNext = () => {
    const { Logoff } = useContext(AuthContext);
    const { id, story } = useParams();
    const safeID: string = story ?? "";
    const [encounters, setEncounters] = useState<Encounter[]>([]);
    const [encounterID, setEncounterID] = useState(0);

    const [formData, setFormData] = useState<NextEncounterType[]>([{ upstream_id: Number(id), encounter_id: 0, next_encounter_id: 0, text: '', objective: { kind: '', values: [] } }]);
  

    const handleDropdownChange = (index: number, event: React.ChangeEvent<HTMLSelectElement>) => {
      setFormData((prevData) => {
        const newData = [...prevData];
        newData[index].next_encounter_id = Number(event.target.value);
        newData[index].encounter_id = encounterID;
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
        newData[index].encounter_id = encounterID;
        return newData;
      });
    };

    const addGroup = () => {
      setFormData((prevData) => [...prevData, { upstream_id: Number(id), encounter_id: encounterID, next_encounter_id: 0, text: '', objective: { kind: '', values: [] } }]);
    };
  
    const removeGroup = (index: number) => {
      setFormData((prevData) => prevData.filter((_, i) => i !== index));
    };

    const navigate = useNavigate();

    const { t } = useTranslation(["home", "main"]);
  
    useEffect(() => {
        FetchEncounters(safeID, setEncounters);
    }, []);
  
    const cancelButton = () => {
      navigate(`/autoplay/${id}/story/${story}`);
    };
  
    async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
      e.preventDefault();
      const apiURL = UseLocation();
      const urlAPI = new URL("api/v1/autoplay/next", apiURL);
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
        alert(t("alert.next-encounter", {ns: ['main', 'home']}));
        navigate(`/autoplay/${id}/story/${story}`);
      } else {
        alert(t("alert.next-encounter-wrong", {ns: ['main', 'home']}));
      };
    };
  
    return (
        <>
          <div className="container mt-3" key="1">
              <Layout Logoff={Logoff} />
              <h2>{t("encounter.add-next-encounter", {ns: ['main', 'home']})}</h2>
          </div>
          <div className="container mt-3" key="2">
              <Form onSubmit={handleSubmit}>
              <Form.Group className="mb-3" controlId="formNextEncounter">
                  <Form.Label>{t("encounter.this", {ns: ['main', 'home']})}</Form.Label>
                  <Form.Select key={1} as="select"
                    value={encounterID}
                    onChange={e => {
                      console.log("set e.target.value", e.target.value);
                      setEncounterID(Number(e.target.value));
                    }}>
                      <option value="-1">{t("encounter.select-encounter", {ns: ['main', 'home']})}</option>
                    {
                      encounters != null ? (
                          encounters.map((encounter: Encounter) => (
                            <option value={encounter.id}>{encounter.title}
                            {
                                encounter.first_encounter ? " (First Encounter)" : null
                            }
                            </option>
                        ))) : (
                          <option>{t("encounter.not-found", {ns: ['main', 'home']})}</option>
                        )
                    }
                  </Form.Select>
                  <Form.Text className="text-muted">
                  {t("encounter.next-text1", {ns: ['main', 'home']})}
                  </Form.Text>
              </Form.Group>
              {formData.map((group, index) => (
                <div key={index}>
                  <Form.Group>
                    <Form.Label>{t("encounter.next", {ns: ['main', 'home']})}</Form.Label>
                    <Form.Select key={index} name="next_encounter_id" value={group.next_encounter_id} onChange={(e) => handleDropdownChange(index, e)}>
                      <option value="-1">{t("encounter.select-encounter", {ns: ['main', 'home']})}</option>
                      {
                        encounters != null ? (
                          encounters.filter(encounter => encounter.id !== encounterID)
                            .map((encounter) => (
                              <option value={encounter.id}>{encounter.title}
                              {
                                  encounter.first_encounter ? " (First Encounter)" : null
                              }
                              </option>
                          ))) : (
                            <option>{t("encounter.not-found", {ns: ['main', 'home']})}</option>
                          )
                      }
                    </Form.Select>
                    <Form.Text className="text-muted">{t("encounter.next-text3", {ns: ['main', 'home']})}</Form.Text>
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
                      <option value="invalid">{t("encounter.select-objective", {ns: ['main', 'home']})}</option>
                      <option value="no_action">{t("encounter.select-objective-1", {ns: ['main', 'home']})}</option>
                      <option value="dice_roll">{t("encounter.select-objective-2", {ns: ['main', 'home']})}</option>
                    </Form.Select>
                  </Form.Group>
                  <Button variant="danger" onClick={() => removeGroup(index)}>{t("encounter.remove", {ns: ['main', 'home']})}</Button>
                </div>
              ))}    
                  <Button variant="primary" onClick={addGroup}>{t("encounter.add-more", {ns: ['main', 'home']})}</Button>{" "}
                  <Button variant="primary" type="submit" >
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
    
  export default AutoPlayNext;