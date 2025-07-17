import { Outlet, NavLink } from "react-router-dom";
import "bootstrap/dist/css/bootstrap.css";
import { useTranslation } from "react-i18next";
import SaveLanguage, { GetLanguage } from "../context/SaveLanguage";

interface Props {
  Logoff: (newState: boolean) => void;
}

const Layout = ({ Logoff }: Props) => {
  const { t, i18n } = useTranslation(['home', 'main']);

  const usedLanguage = GetLanguage();

  const onClickLanguageChange = (e: any) => {
      const language = e.target.value;
      SaveLanguage(language);  //set the language
      i18n.changeLanguage(language); //change the language
  }
  const handleLogoff = () => {
    Logoff(false);
  };
  return (
    <>
      <nav className="navbar bg-body-tertiary">
        <ul className="nav nav-pills nav-fill">
          <li className="nav-item">
            <NavLink className="nav-link" aria-current="page" to="/">
            {t("common.home", {ns: ['main','home']})}
            </NavLink>{" "}
          </li>
          <li className="nav-item">
            <NavLink className="nav-link" to="/stories">
            {t("common.story", {ns: ['main','home']})}
            </NavLink>{" "}
          </li>
          <li className="nav-item">
            <NavLink className="nav-link" to="/tasks">
            {t("common.task", {ns: ['main','home']})}
            </NavLink>{" "}
          </li>
          <li className="nav-item">
            <NavLink className="nav-link" to="/users">
            {t("common.user", {ns: ['main','home']})}
            </NavLink>{" "}
          </li>
          <li className="nav-item">
            <NavLink className="nav-link" to="/characters">
            {t("common.character", {ns: ['main','home']})}
            </NavLink>{" "}
          </li>
          <li className="nav-item">
            <NavLink className="nav-link" to="/stages">
            {t("common.stage", {ns: ['main','home']})}
            </NavLink>{" "}
          </li>
          <li className="nav-item">
            <NavLink className="nav-link" to="/autoplay">
            {t("common.auto-play", {ns: ['main','home']})}
            </NavLink>{" "}
          </li>
          <li>
              <select className="form-select" value={usedLanguage} onChange={onClickLanguageChange}>
                <option value="en" >English</option>
                <option value="pt" >Português Brasileiro</option>
              </select>
          </li>
        </ul>

        <div className="d-grid gap-2 d-md-flex justify-content-md-end">
          <span>
            <button className="btn btn-secondary" onClick={handleLogoff}>
            {t("common.logoff", {ns: ['main','home']})}
            </button>{" "}
          </span>
        </div>
      </nav>
      <Outlet />
    </>
  );
};

export default Layout;
