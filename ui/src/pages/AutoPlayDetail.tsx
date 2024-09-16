import { useContext, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import AutoPlayDetailHeader from "../components/AutoPlayDetailHeader";
import { FetchEncountersAutoPlay } from "../functions/AutoPlay";
import { AutoPlayEncounterList } from "../types/AutoPlay";

const AutoPlayDetail = () => {
    const { id, story } = useParams();
    const { Logoff } = useContext(AuthContext);
  
    const safeID: string = id ?? "";
  
    const storySafeID: string = story ?? "";
  
    const [encountersList, setEncountersList] = useState<AutoPlayEncounterList>();

  useEffect(() => {
    FetchEncountersAutoPlay(storySafeID, setEncountersList);
  }, []);
  
    return (
      <>
        <div className="container mt-3" key="1">
          <Layout Logoff={Logoff} />
          {<AutoPlayDetailHeader id={safeID} storyID={storySafeID} />}
        </div>
        <div className="container mt-3" key="3">
          <div className="card" >
            <div className="card-header">
              Encounter list
            </div>
            <ul className="list-group list-group-flush">
              {
                encountersList?.encounter_list != null ? (
                  encountersList?.encounter_list.map((encounter) => (
                    <li className="list-group-item" key={encounter.id} >{encounter.name}
                    {
                      encountersList.link != null ? (
                        encountersList.link.filter((next) => next.encounter === encounter.name)
                        .map((next) => (
                          <ul>
                            <li>Linked to: {next.next_encounter}</li>
                          </ul>
                        ))
                      ) : (
                        <p>Links not found</p>
                      )
                    }
                    </li>
                  ))) : (
                  <li className="list-group-item">Encounter not found</li>
                  )
              }
            </ul>
          </div>
        </div>
      </>
    );
  };
  
  export default AutoPlayDetail;
  