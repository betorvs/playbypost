import { useContext } from "react";
import Layout from "../components/Layout";
import { AuthContext } from "../context/AuthContext";
import UsersList from "../components/UsersList";

const UsersPage = () => {
  const { Logoff } = useContext(AuthContext);
  return (
    <>
      <>
        <div className="container mt-3" key="1">
          <Layout Logoff={Logoff} />
          <h2>Users List</h2>
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
