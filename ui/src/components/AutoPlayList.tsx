import { useState, useEffect } from "react";
import AutoPlay from "../types/AutoPlay";
import FetchAutoPlay from "../functions/AutoPlay";
import AutoPlayCards from "./Cards/AutoPlay";

const AutoPlayList = () => {
  const [auto, setAutoPlay] = useState<AutoPlay[]>([]);

  useEffect(() => {
    FetchAutoPlay(setAutoPlay);
  }, []);
  return (
    <div className="container mt-3" key="2">
      {auto != null ? (
        auto.map((autoplay) => (
          <AutoPlayCards
            key={autoplay.id}
            ID={autoplay.id}
            autoPlay={autoplay}
          />
        ))
      ) : (
        <p>no auto play stories for you</p>
      )}
    </div>
  );
};

export default AutoPlayList;