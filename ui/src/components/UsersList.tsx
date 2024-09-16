import { useEffect, useState } from "react";
import UserCards from "./Cards/User";
import UsersCard from "../types/UserCard";
import FetchUsers from "../functions/Users";

const UsersList = () => {
  const [users, setUser] = useState<UsersCard[]>([]);

  useEffect(() => {
    FetchUsers(setUser);
  }, []);
  return (
    <div className="row mb-2" key="1">
      {users != null ? (
        users.map((user, index) => <UserCards user={user} key={index} />)
      ) : (
        <p>no users found</p>
      )}
    </div>
  );
};

export default UsersList;
