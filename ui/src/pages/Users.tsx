import { useContext } from "react";
import Layout from "../components/Layout";
import { AuthContext } from "../context/AuthContext";
import UsersList from "../components/UsersList";
import { useTranslation } from "react-i18next";

const UsersPage = () => {
  const { Logoff } = useContext(AuthContext);
  const { t } = useTranslation(['home', 'main']);
  return (
    <>
      <>
        <div className="container mt-3" key="1">
          <Layout Logoff={Logoff} />
          <h2>{t("user.user-list", {ns: ['main', 'home']})}</h2>
          <hr />
        </div>
        <div className="container mt-3" key="2">
          {<UsersList />}
        </div>
      </>
    </>
  );
};

export default UsersPage;
