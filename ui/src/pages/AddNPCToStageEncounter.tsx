import { useContext, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import { Button, Form } from "react-bootstrap";
import { useNavigate, useParams } from "react-router-dom";
import UseLocation from "../context/UseLocation";
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";



const AddNPCToStageEncounter = () => {
    const { Logoff } = useContext(AuthContext);
    const { id, story, encounterid, storyteller_id } = useParams();

    const navigate = useNavigate();

    const [name, setName] = useState("");

    const cancelButton = () => {
        navigate(`/stages/${id}/story/${story}/encounter/${encounterid}`);
      };

    async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
        
        const apiURL = UseLocation();
        const urlAPI = new URL("api/v1/stage/npc", apiURL);
        const response = await fetch(urlAPI, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "X-Username": GetUsername(),
            "X-Access-Token": GetToken(),
          },
          body: JSON.stringify({
            storyteller_id: Number(storyteller_id),
            stage_id: Number(id),
            encounter_id: Number(encounterid),
            name: name,
          }),
        });
        if (response.ok) {
          alert("NPC generated! Have a great session with your friends!");
          navigate(`/stages/${id}/story/${story}/encounter/${encounterid}`);
        } else {
          alert("Something goes wrong. NPC was not generated.");
        }
    };

    return (
        <>
        <div className="container mt-3" key="1">
            <Layout Logoff={Logoff} />
            <h2>Add NPC to Encounter</h2>
        </div>
        <div className="container mt-3" key="2">
            <Form onSubmit={handleSubmit}>
                <Form.Group className="mb-3" controlId="formNPC">
                    <Form.Label>NPC Name</Form.Label>
                    <Form.Control
                      type="text"
                      placeholder="name"
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                    />
                    <Form.Text className="text-muted">
                      Name of the NPC
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
    
export default AddNPCToStageEncounter;