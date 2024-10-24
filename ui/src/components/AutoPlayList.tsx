import { useState, useEffect } from "react";
import AutoPlay from "../types/AutoPlay";
import FetchAutoPlay from "../functions/AutoPlay";
import AutoPlayCards from "./Cards/AutoPlay";
import { useTranslation } from "react-i18next";

const AutoPlayList = () => {
  const [auto, setAutoPlay] = useState<AutoPlay[]>([]);
  const { t } = useTranslation(['home', 'main']);

  useEffect(() => {
    FetchAutoPlay(setAutoPlay);
  }, []);
  return (
    <div className="container mt-3" key="2">
      {auto.length !== 0 ? (
        auto.map((autoplay) => (
          <AutoPlayCards
            key={autoplay.id}
            ID={autoplay.id}
            autoPlay={autoplay}
          />
        ))
      ) : (
        <p>{t("auto-play.not-found", {ns: ['main', 'home']})}</p>
      )}
    </div>
  );
};

export default AutoPlayList;