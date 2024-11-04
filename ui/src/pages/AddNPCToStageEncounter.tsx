import { useContext, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import { Button, Form } from "react-bootstrap";
import { useNavigate, useParams } from "react-router-dom";
import UseLocation from "../context/UseLocation";
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import { useTranslation } from "react-i18next";


const AddNPCToStageEncounter = () => {
    const { Logoff } = useContext(AuthContext);
    const { id, story, encounterid, storyteller_id } = useParams();

    const navigate = useNavigate();

    const { t } = useTranslation(['home', 'main']);

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
          alert(t("alert.npc", {ns: ['main', 'home']}));
          navigate(`/stages/${id}/story/${story}/encounter/${encounterid}`);
        } else {
          const data = await response.text();
          let error = JSON.parse(data);
          console.log(error);
          alert(t("alert.npc-wrong", {ns: ['main', 'home']}) + "\n" + error.msg);
        }
    };

    return (
        <>
        <div className="container mt-3" key="1">
            <Layout Logoff={Logoff} />
            <h2>{t("player.add-npc", {ns: ['main', 'home']})}</h2>
        </div>
        <div className="container mt-3" key="2">
            <Form onSubmit={handleSubmit}>
                <Form.Group className="mb-3" controlId="formNPC">
                    <Form.Label>{t("player.npc-name", {ns: ['main', 'home']})}</Form.Label>
                    <Form.Control
                      type="text"
                      placeholder="name"
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                    />
                    <Form.Text className="text-muted">
                    {t("player.npc-name-text", {ns: ['main', 'home']})}
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
    
export default AddNPCToStageEncounter;