import { useContext } from "react";
import Layout from "../components/Layout";
import { AuthContext } from "../context/AuthContext";
import PlayersList from "../components/PlayersList";
import { useParams } from "react-router-dom";
import NavigateButton from "../components/Button/NavigateButton";

const PlayersPage = () => {
  const { id, story } = useParams();
  const { Logoff } = useContext(AuthContext);

  const safeID: string = id ?? "";
  const storySafeID: string = story ?? "";

  return (
    <>
      <>
        <div className="container mt-3" key="1">
          <Layout Logoff={Logoff} />
          <h2>Players List</h2>
          <NavigateButton link={`/stages/${safeID}/story/${storySafeID}`} variant="secondary">
            Back to Stage
          </NavigateButton>{" "}
          <hr />
        </div>
        <div className="container mt-3" key="2">
          {<PlayersList id={safeID} />}
        </div>
      </>
    </>
  );
};

export default PlayersPage;
