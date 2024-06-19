import { useContext } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import StoryList from "../components/StoryList";

const StagesPage = () => {
  const { Logoff } = useContext(AuthContext);

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>Stages</h2>
        <hr />
      </div>
      {<StoryList />}
    </>
  );
};

export default StagesPage;
