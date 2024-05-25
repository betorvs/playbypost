import { useContext } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import StoryList from "../components/StoryList";
import NavigateButton from "../components/Button/NavigateButton";

const StoriesPage = () => {
  const { Logoff } = useContext(AuthContext);

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>Stories</h2>
        <NavigateButton link="/stories/new" variant="primary">
          New Story
        </NavigateButton>{" "}
        <hr />
      </div>
      {<StoryList />}
    </>
  );
};

export default StoriesPage;
