// import { useState } from "react";

import { useEffect, useState } from "react";
import NavigateButton from "./Button/NavigateButton";
import AutoPlay from "../types/AutoPlay";
import { FetchAutoPlayByID } from "../functions/AutoPlay";


interface props {
    id: string;
    storyID: string;
}

const AutoPlayDetailHeader = ({ id, storyID }: props) => {
  const [autoPlay, setAutoPLay] = useState<AutoPlay>();

  useEffect(() => {
    FetchAutoPlayByID(id, setAutoPLay);
  }, []);
    return (
      <div
        className="p-4 p-md-5 mb-4 rounded text-body-emphasis bg-body-secondary"
        key="1"
      >
        <div className="col-lg-6 px-0" key="1">
          <h1 className="display-4 fst-italic">
            { autoPlay?.text || "Title not found"}
          </h1>
          <p className="lead my-3">
            ID: { autoPlay?.id || "AutoPlayID not found"} 
          </p>
          <p className="lead mb-0">Story ID: { autoPlay?.story_id || "Story not found"}</p>
          <br />
          <p className="lead mb-0">
          {
            autoPlay?.solo ? (
              <p>This is a Solo Adventure</p>
            ) : (
              <p>This is Didatic Adventure</p>
            )
          }
          </p>
          <br />
         
          <NavigateButton link="/autoplay" variant="secondary">
            Back to Auto Play
          </NavigateButton>{" "}
          <NavigateButton link={`/autoplay/${id}/story/${storyID}/next`} variant="primary">
            Add Next Encounter
          </NavigateButton>{" "}
          
        </div>
      </div>
    );
  };
  
  export default AutoPlayDetailHeader;