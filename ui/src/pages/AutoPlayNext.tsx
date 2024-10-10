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
        alert("Next Encounter added! Have a great session with your friends!");
        navigate(`/autoplay/${id}/story/${story}`);
      } else {
        alert("Something goes wrong. Next Encounter was not added to encounter");
      };
    };
  
    return (
        <>
          <div className="container mt-3" key="1">
              <Layout Logoff={Logoff} />
              <h2>Add Auto Play Next Encounter</h2>
          </div>
          <div className="container mt-3" key="2">
              <Form onSubmit={handleSubmit}>
              <Form.Group className="mb-3" controlId="formNextEncounter">
                  <Form.Label>Encounter</Form.Label>
                  <Form.Select as="select"
                    value={encounterID}
                    onChange={e => {
                      console.log("set e.target.value", e.target.value);
                      setEncounterID(Number(e.target.value));
                    }}>
                      <option value="-1">Select a Encounter</option>
                    {
                      encounters != null ? (
                          encounters.map((encounter: Encounter) => (
                            <option value={encounter.id}>{encounter.title}
                            {
                                encounter.first_encounter ? " (First Encounter)" : null
                            }
                            </option>
                        ))) : (
                          <option>encounters not found</option>
                        )
                    }
                  </Form.Select>
                  <Form.Text className="text-muted">
                    Choose a Encounter to link to the next Encounter
                  </Form.Text>
              </Form.Group>
              {formData.map((group, index) => (
                <div key={index}>
                  <Form.Group>
                    <Form.Label>Next Encounter</Form.Label>
                    <Form.Select name="next_encounter_id" value={group.next_encounter_id} onChange={(e) => handleDropdownChange(index, e)}>
                      <option value="-1">Select a Encounter</option>
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
                            <option>encounters not found</option>
                          )
                      }
                    </Form.Select>
                    <Form.Text className="text-muted">It should be different. Each pair of Next Encounter with Optional Name should be unique.</Form.Text>
                  </Form.Group>
                  <Form.Group>
                    <Form.Label>Optional Name</Form.Label>
                    <Form.Control type="text" name="text" value={group.text} onChange={(e) => handleInputChange(index, e)} />
                    <Form.Text className="text-muted">
                      It can be assigned multiple times, each of it will enable a option for player to choose his destiny. Keep option name different per encounter.
                    </Form.Text>
                  </Form.Group>
                  <Form.Group>
                    <Form.Label>Automatic Option</Form.Label>
                    <Form.Select name="objective" value={group.objective.kind} onChange={(e) => handleDropdownObjectiveChange(index, e)}>
                      <option value="invalid">Select a Objective</option>
                      <option value="no_action">Free Choice</option>
                      <option value="dice_roll">Dice Roll</option>
                    </Form.Select>
                  </Form.Group>
                  <Button variant="danger" onClick={() => removeGroup(index)}>Remove</Button>
                </div>
              ))}    
                  <Button variant="primary" onClick={addGroup}>Add More Encounters</Button>{" "}
                  <Button variant="primary" type="submit" >
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
    
  export default AutoPlayNext;