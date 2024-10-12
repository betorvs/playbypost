import { useNavigate, useParams } from "react-router-dom";
import Layout from "../components/Layout";
import { Button, Form } from "react-bootstrap";
import { useContext, useEffect, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import UseLocation from "../context/UseLocation";
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import Players from "../types/Players";
import FetchPlayers from "../functions/Players";
import { useTranslation } from "react-i18next";

const AddPlayerToStageEncounter = () => {
    const { Logoff } = useContext(AuthContext);
    const { id, story, encounterid } = useParams();

    const navigate = useNavigate();

    const { t } = useTranslation(['home', 'main']);

    const safeID: string = id ?? "";

    const [players, setPlayer] = useState<Players[]>();
    const [ids, setIDs] = useState<number[]>([]);

    useEffect(() => {
      FetchPlayers(safeID, setPlayer);
    }, []);


    const cancelButton = () => {
      navigate(`/stages/${id}/story/${story}/encounter/${encounterid}`);
    };
    
    async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
   
        const apiURL = UseLocation();
        const urlAPI = new URL("api/v1/stage/encounter/participants", apiURL);
        const response = await fetch(urlAPI, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "X-Username": GetUsername(),
            "X-Access-Token": GetToken(),
          },
          body: JSON.stringify({
            identifies: ids,
            encounter_id: Number(encounterid),
          }),
        });
        if (response.ok) {
          alert("Players added! Have a great session with your friends!");
          navigate(`/stages/${id}/story/${story}/encounter/${encounterid}`);
        } else {
          alert("Something goes wrong. Players was not added to encounter.");
        }
    }
    console.log("starting with values:", ids);

    return (
        <>
        <div className="container mt-3" key="1">
            <Layout Logoff={Logoff} />
            <h2>{t("player.add-player", {ns: ['main', 'home']})}</h2>
        </div>
        <div className="container mt-3" key="2">
            <Form onSubmit={handleSubmit}>
                <Form.Group className="mb-3" controlId="formPlayers">
                    <Form.Label>{t("player.this", {ns: ['main', 'home']})}</Form.Label>
                    {
                        players != null ? (
                            players.map((player) => (
                                <Form.Check type="checkbox" id={String(player.id)} value={player.id} label={player.name}
                                  onClick={() => {
                                      const numValue = Number(player.id);
                                      if (ids.includes(numValue)) {
                                          setIDs(ids.filter((id) => id !== numValue));
                                      } else {
                                          setIDs([...ids, numValue]);
                                      }
                                  }}
                                />
                            ))) : (
                                <Form.Check disabled type="checkbox" label={`disabled `} id={`disabled-default`} />
                            )
                    }
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

export default AddPlayerToStageEncounter;