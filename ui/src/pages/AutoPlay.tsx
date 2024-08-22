import { useContext } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import NavigateButton from "../components/Button/NavigateButton";
import AutoPlayList from "../components/AutoPlayList";

const AutoPlayPage = () => {
  const { Logoff } = useContext(AuthContext);

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>Auto Play Stories</h2>
        <NavigateButton link="/autoplay/new" variant="primary">
          Add Story to Auto Play
        </NavigateButton>{" "}
        <hr />
      </div>
      {<AutoPlayList />}
    </>
  );
};

export default AutoPlayPage;