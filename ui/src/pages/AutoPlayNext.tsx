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



const AutoPlayNext = () => {
    const { Logoff } = useContext(AuthContext);
    const { id, story } = useParams();
    const [text, setText] = useState("");
    const safeID: string = story ?? "";
    const [encounters, setEncounters] = useState<Encounter[]>([]);
    const [encounterID, setEncounterID] = useState(0);
    const [nextEncounterID, setNextEncounterID] = useState(0);
  
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
        body: JSON.stringify({
          upstream_id: Number(id),
          encounter_id: encounterID,
          next_encounter_id: nextEncounterID,
          text: text,
        }),
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
                            encounters.filter(encounter => encounter.id !== nextEncounterID)
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
                      <Form.Text className="text-muted">
                        Choose a Encounter to link to the next Encounter
                      </Form.Text>
                  </Form.Group>
                  <Form.Group className="mb-3" controlId="formNextEncounter">
                      <Form.Label>Next Encounter</Form.Label>
                      <Form.Select as="select"
                        value={nextEncounterID}
                        onChange={e => {
                          console.log("set e.target.value", e.target.value);
                          setNextEncounterID(Number(e.target.value));
                        }}>
                          <option value="-1">Select a Encounter</option>
                        {
                          encounters != null ? (
                            encounters.filter(encounter => encounter.id !== encounterID)
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
                  <Form.Group className="mb-3" controlId="formName">
                      <Form.Label>Optional Name</Form.Label>
                      <Form.Control
                        type="text"
                        placeholder="name"
                        value={text}
                        onChange={(e) => setText(e.target.value)}
                      />
                      <Form.Text className="text-muted">
                      It can be assigned multiple times, each of it will enable a option for player to choose his destiny. Keep option name different per encounter.
                      </Form.Text>
                  </Form.Group>
                  
                  <Button variant="primary" type="submit">
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