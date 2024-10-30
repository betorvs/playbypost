import { AuthContext } from "../context/AuthContext";
import { useContext } from "react";
import Layout from "../components/Layout";
import { useTranslation } from "react-i18next";

const HomePublicPage = () => {
  const { Logoff } = useContext(AuthContext);
  const { t } = useTranslation(['home', 'main']);

  return (
    <div className="container mt-3" key="1">
      <Layout Logoff={Logoff} />
      <h2>{t("home.title", {ns: ['main', 'home']})}</h2>
      <hr />
      <p>{t("home.description", {ns: ['main', 'home']})}</p>
      <hr />
      <h4>{t("home.sub-header", {ns: ['main', 'home']})}</h4>
      <ul>
        <li>{t("home.list-item-1", {ns: ['main', 'home']})}: {t("home.item-1-description", {ns: ['main', 'home']})}</li>
        <li>{t("home.list-item-2", {ns: ['main', 'home']})}: {t("home.item-2-description", {ns: ['main', 'home']})}</li>
        <li>{t("home.list-item-3", {ns: ['main', 'home']})}: {t("home.item-3-description", {ns: ['main', 'home']})}</li>
        <li>{t("home.list-item-4", {ns: ['main', 'home']})}: {t("home.item-4-description", {ns: ['main', 'home']})}</li>
        <li>{t("home.list-item-5", {ns: ['main', 'home']})}: {t("home.item-5-description", {ns: ['main', 'home']})}</li>
        <li>{t("home.list-item-6", {ns: ['main', 'home']})}: {t("home.item-6-description", {ns: ['main', 'home']})}</li>
        <li>{t("home.list-item-7", {ns: ['main', 'home']})}: {t("home.item-7-description", {ns: ['main', 'home']})}</li>
        <li>{t("home.list-item-8", {ns: ['main', 'home']})}: {t("home.item-8-description", {ns: ['main', 'home']})}</li>
        <li>{t("home.list-item-9", {ns: ['main', 'home']})}: {t("home.item-9-description", {ns: ['main', 'home']})}</li>
      </ul>
      <hr />
      <h4>{t("home.sub-header-2", {ns: ['main', 'home']})}</h4>
      <ol>
        <li>{t("home.list-item-10", {ns: ['main', 'home']})}</li>
        <li>{t("home.list-item-11", {ns: ['main', 'home']})}</li>
        <li>{t("home.list-item-12", {ns: ['main', 'home']})}</li>
          <ul>
            <li>{t("home.list-item-13", {ns: ['main', 'home']})}</li>
            <li>{t("home.list-item-14", {ns: ['main', 'home']})}</li>
          </ul>
        <li>{t("home.list-item-15", {ns: ['main', 'home']})}</li>
      </ol>
      <p>{t("home.notes-1", {ns: ['main', 'home']})}</p>
      <h6>{t("home.sub-header-3", {ns: ['main', 'home']})}</h6>
      <ol>
        <li>{t("home.list-item-16", {ns: ['main', 'home']})}</li>
        <li>{t("home.list-item-17", {ns: ['main', 'home']})}</li>
        <li>{t("home.list-item-18", {ns: ['main', 'home']})}</li>
        <li>{t("home.list-item-19", {ns: ['main', 'home']})}</li>
        <li>{t("home.list-item-20", {ns: ['main', 'home']})}</li>
        <li>{t("home.list-item-21", {ns: ['main', 'home']})}</li>
      </ol>
      <hr />
      <p>{t("home.sub-header-4", {ns: ['main', 'home']})}</p>
    </div>
  );
};

export default HomePublicPage;
