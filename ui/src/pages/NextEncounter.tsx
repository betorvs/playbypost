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


const NextEncounter = () => {
  const { Logoff } = useContext(AuthContext);
  const { id, story, encounterid } = useParams();
  const stageID: string = id ?? "";
  const safeID: string = encounterid ?? "";
  const encounteridNumber: number = Number(safeID);
  const [encounters, setEncounters] = useState<Encounter[]>([]);

  const [formData, setFormData] = useState<NextEncounterType[]>([{ upstream_id: Number(id), encounter_id: encounteridNumber, next_encounter_id: 0, text: '' }]);

  const handleDropdownChange = (index: number, event: React.ChangeEvent<HTMLSelectElement>) => {
    setFormData((prevData) => {
      const newData = [...prevData];
      newData[index].next_encounter_id = Number(event.target.value);
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
    setFormData((prevData) => [...prevData, { upstream_id: Number(id), encounter_id: encounteridNumber, next_encounter_id: 0, text: '' }]);
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
            <h2>Add Next Encounter</h2>
            <h4>Encounter  ID: {encounteridNumber}</h4>
        </div>
        <div className="container mt-3" key="2">
            <Form onSubmit={handleSubmit}>
            {formData.map((group, index) => (
              <div key={index}>
                <Form.Group>
                  <Form.Label>Next Encounter</Form.Label>
                  <Form.Select name="next_encounter_id" value={group.next_encounter_id} onChange={(e) => handleDropdownChange(index, e)}>
                    <option value="-1">Select a Encounter</option>
                        {
                          encounters != null ? (
                            encounters.filter(encounter => encounter.id !== encounteridNumber)
                              .map((encounter) => (
                                <option value={encounter.id}>{encounter.title}</option>
                            ))) : (
                              <option>encounters not found</option>
                            )
                        }
                  </Form.Select>
                  <Form.Text className="text-muted">
                    Choose the next Encounter
                  </Form.Text>
                </Form.Group>
                <Form.Group>
                  <Form.Label>Optional Name</Form.Label>
                  <Form.Control type="text" name="text" value={group.text} onChange={(e) => handleInputChange(index, e)} />
                  <Form.Text className="text-muted">
                    It can be assigned multiple times, each of it will enable a option for player to choose his destiny. Keep option name different per encounter.
                  </Form.Text>
                </Form.Group>
                <Button variant="danger" onClick={() => removeGroup(index)}>Remove</Button>
              </div>

            ))}
                <Button variant="primary" onClick={addGroup}>Add More Encounters</Button>{" "}
                <Button variant="primary" type="submit">Submit</Button>{" "}
                <Button variant="secondary" onClick={() => cancelButton()}>
                Cancel
                </Button>{" "}
            </Form>
        </div>
      </>
  );
};
  
export default NextEncounter;