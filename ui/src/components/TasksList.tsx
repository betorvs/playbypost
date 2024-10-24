import { useEffect, useState } from "react";
import FetchTasks from "../functions/Tasks";
import Task from "../types/Task";
import TaskCard from "./Cards/Task";
import { useTranslation } from "react-i18next";

const TasksList = () => {
  const [tasks, setTask] = useState<Task[]>([]);
  const { t } = useTranslation(['home', 'main']);

  useEffect(() => {
    FetchTasks(setTask);
  }, []);
  return (
    <div className="row mb-2" key="1">
      {tasks.length !== 0 ? (
        tasks.map((task, index) => <TaskCard task={task} key={index} />)
      ) : (
        <p>{t("task.not-found", {ns: ['main', 'home']})}</p>
      )}
    </div>
  );
};

export default TasksList;
