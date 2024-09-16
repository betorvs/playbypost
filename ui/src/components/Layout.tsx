import { Outlet, NavLink } from "react-router-dom";
import "bootstrap/dist/css/bootstrap.css";

interface Props {
  Logoff: (newState: boolean) => void;
}

const Layout = ({ Logoff }: Props) => {
  const handleLogoff = () => {
    Logoff(false);
  };
  return (
    <>
      <nav className="navbar bg-body-tertiary">
        <ul className="nav nav-pills nav-fill">
          <li className="nav-item">
            <NavLink className="nav-link" aria-current="page" to="/">
              Home
            </NavLink>{" "}
          </li>
          <li className="nav-item">
            <NavLink className="nav-link" to="/stories">
              Stories
            </NavLink>{" "}
          </li>
          <li className="nav-item">
            <NavLink className="nav-link" to="/tasks">
              Tasks
            </NavLink>{" "}
          </li>
          <li className="nav-item">
            <NavLink className="nav-link" to="/users">
              Users
            </NavLink>{" "}
          </li>
          <li className="nav-item">
            <NavLink className="nav-link" to="/stages">
              Stages
            </NavLink>{" "}
          </li>
          <li className="nav-item">
            <NavLink className="nav-link" to="/autoplay">
              Auto Play
            </NavLink>{" "}
          </li>
          <li></li>
        </ul>

        <div className="d-grid gap-2 d-md-flex justify-content-md-end">
          <span>
            <button className="btn btn-secondary" onClick={handleLogoff}>
              Logoff
            </button>{" "}
          </span>
        </div>
      </nav>
      <Outlet />
    </>
  );
};

export default Layout;
