import { useEffect, useState } from "react";
import UserCards from "./Cards/User";
import UsersCard from "../types/UserCard";
import FetchUsers from "../functions/Users";
import { useTranslation } from "react-i18next";

const UsersList = () => {
  const [users, setUser] = useState<UsersCard[]>([]);
  const { t } = useTranslation(['home', 'main']);

  useEffect(() => {
    FetchUsers(setUser);
  }, []);
  return (
    <div className="row mb-2" key="1">
      {users != null ? (
        users.map((user, index) => <UserCards user={user} key={index} />)
      ) : (
        <p>{t("user.not-found", {ns: ['main', 'home']})}</p>
      )}
    </div>
  );
};

export default UsersList;
