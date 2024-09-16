import { useEffect, useState } from "react";
import FetchTasks from "../functions/Tasks";
import Task from "../types/Task";
import TaskCard from "./Cards/Task";

const TasksList = () => {
  const [tasks, setTask] = useState<Task[]>([]);

  useEffect(() => {
    FetchTasks(setTask);
  }, []);
  return (
    <div className="row mb-2" key="1">
      {tasks != null ? (
        tasks.map((task, index) => <TaskCard task={task} key={index} />)
      ) : (
        <p>no tasks found</p>
      )}
    </div>
  );
};

export default TasksList;
