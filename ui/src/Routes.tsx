import { useContext } from "react";
import {
  Routes as Router,
  Route,
  Navigate,
  Outlet,
  useParams,
} from "react-router-dom";
import { AuthContext } from "./context/AuthContext";
import Home from "./pages/Home";
import Login from "./pages/Login";
import StoriesPage from "./pages/Stories";
import UsersPage from "./pages/Users";
import StoryPlayers from "./pages/StoryPlayers";
import StoryDetail from "./pages/StoryDetail";
import NewStory from "./pages/NewStory";
import NewEncounter from "./pages/NewEncounter";

const PrivateRoutes = () => {
  const { authenticated } = useContext(AuthContext);

  // if (!authenticated) return <Navigate to="/login" replace />;

  // return <Outlet />;
  return authenticated ? <Outlet /> : <Navigate to="/login" replace />;
};

function NoMatch() {
  return (
    <div className="container mt-3">
      <h2>404: Page Not Found</h2>
      <p>Try again!</p>
    </div>
  );
}

function SlugTest() {
  const { id } = useParams();
  return (
    <div className="container mt-3">
      <h2>Slug Page</h2>
      <p>value {id}</p>
    </div>
  );
}

const Routes = () => {
  return (
    <Router>
      <Route path="/login" element={<Login />} />
      <Route path="/slug/:id" element={<SlugTest />} />

      <Route element={<PrivateRoutes />}>
        <Route path="/" element={<Home />} />
        <Route path="/stories">
          <Route path="/stories" element={<StoriesPage />} />
          <Route path="/stories/new" element={<NewStory />} />
          <Route path="/stories/:id" element={<StoryDetail />} />
          <Route path="/stories/players/:id" element={<StoryPlayers />} />
          <Route path="/stories/encounter/new/:id" element={<NewEncounter />} />
        </Route>
        <Route path="/users" element={<UsersPage />} />
        <Route path="*" element={<NoMatch />} />
      </Route>
    </Router>
  );
};

export default Routes;
