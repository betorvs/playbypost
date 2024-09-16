import { AuthContext } from "../context/AuthContext";
import { useContext } from "react";
import Layout from "../components/Layout";

const HomePublicPage = () => {
  const { Logoff } = useContext(AuthContext);

  return (
    <div className="container mt-3" key="1">
      <Layout Logoff={Logoff} />
    </div>
  );
};

export default HomePublicPage;
