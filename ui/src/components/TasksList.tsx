import { useEffect, useState } from "react";
import UserCards from "./Cards/User";
import UsersCard from "../types/UserCard";
import FetchTasks from "../functions/Tasks";

const TasksList = () => {
  const [users, setUser] = useState<UsersCard[]>([]);

  useEffect(() => {
    FetchTasks(setUser);
  }, []);
  return (
    <div className="row mb-2" key="1">
      {users != null ? (
        users.map((user, index) => <UserCards user={user} key={index} />)
      ) : (
        <p>no tasks found</p>
      )}
    </div>
  );
};

export default TasksList;
