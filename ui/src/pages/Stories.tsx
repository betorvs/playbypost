import { useContext } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import StoryList from "../components/StoryList";
import NavigateButton from "../components/Button/NavigateButton";
import { useTranslation } from "react-i18next";

const StoriesPage = () => {
  const { Logoff } = useContext(AuthContext);
  const { t } = useTranslation(['home', 'main']);

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>{t("story.header", {ns: ['main', 'home']})}</h2>
        <NavigateButton link="/stories/new" variant="primary">
        {t("story.button", {ns: ['main', 'home']})}
        </NavigateButton>{" "}
        <hr />
      </div>
      {<StoryList />}
    </>
  );
};

export default StoriesPage;
