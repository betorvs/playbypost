import { useContext } from "react";
import Layout from "../components/Layout";
import { AuthContext } from "../context/AuthContext";
import TasksList from "../components/TasksList";
import NavigateButton from "../components/Button/NavigateButton";
import { useTranslation } from "react-i18next";

const TasksPage = () => {
  const { Logoff } = useContext(AuthContext);
  const { t } = useTranslation(['home', 'main']);
  return (
    <>
      <>
        <div className="container mt-3" key="1">
          <Layout Logoff={Logoff} />
          <h2>{t("task.header", {ns: ['main', 'home']})}</h2>
          <NavigateButton link="/tasks/new" variant="primary">
           {t("task.button", {ns: ['main', 'home']})}
          </NavigateButton>{" "}
          <hr />
        </div>
        <div className="container mt-3" key="2">
          {<TasksList />}
        </div>
      </>
    </>
  );
};

export default TasksPage;
