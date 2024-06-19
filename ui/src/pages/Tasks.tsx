import { useContext } from "react";
import Layout from "../components/Layout";
import { AuthContext } from "../context/AuthContext";
import TasksList from "../components/TasksList";

const TasksPage = () => {
  const { Logoff } = useContext(AuthContext);
  return (
    <>
      <>
        <div className="container mt-3" key="1">
          <Layout Logoff={Logoff} />
          <h2>Tasks List</h2>
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
