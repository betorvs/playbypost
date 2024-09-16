import { useParams } from "react-router-dom";
import Layout from "../components/Layout";
import { AuthContext } from "../context/AuthContext";
import { useContext } from "react";
import StoryDetailHeader from "../components/StoryDetailHeader";
import PlayersList from "../components/PlayersList";

const StoryPlayers = () => {
  const { id } = useParams();
  const { Logoff } = useContext(AuthContext);

  const safeID: string = id ?? "";

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        {<StoryDetailHeader detail={false} id={safeID} />}
        <div className="row mb-2" key="2">
          <PlayersList id={safeID} />
        </div>
      </div>
    </>
  );
};

export default StoryPlayers;
