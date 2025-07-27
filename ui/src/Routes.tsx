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
import SessionMonitorPage from "./pages/SessionMonitorPage";
import StoryPlayers from "./pages/StoryPlayers";
import StoryDetail from "./pages/StoryDetail";
import NewStory from "./pages/NewStory";
import NewEncounter from "./pages/NewEncounter";
import TasksPage from "./pages/Tasks";
import StagesPage from "./pages/Stages";
import StageDetail from "./pages/StageDetail";
import NewTask from "./pages/NewTask";
import UserAsStoryteller from "./pages/UserAsStoryteller";
import UserAsPlayer from "./pages/UserAsPlayer";
import PlayersPage from "./pages/Players";
import StageStart from "./pages/StageStart";
import EncounterToStage from "./pages/EncounterToStage";
import StageEncounterDetail from "./pages/StageEncounterDetail";
import AddPlayerToStageEncounter from "./pages/AddPlayersToStageEncounter";
import TaskToEncounter from "./pages/TaskToEncounter";
import NextEncounter from "./pages/NextEncounter";
import AddNPCToStageEncounter from "./pages/AddNPCToStageEncounter";
import AutoPlayPage from "./pages/AutoPlay";
import AutoPlayAdd from "./pages/AutoPlayAdd";
import AutoPlayDetail from "./pages/AutoPlayDetail";
import AutoPlayNext from "./pages/AutoPlayNext";
import StageManageNextEncounter from "./pages/StageManageNextEncounter";
import StageNextEncounter from "./pages/StageNextEncounter";
import EditStory from "./pages/EditStory";
import EditEncounter from "./pages/EditEncounter";
import EditTask from "./pages/EditTask";
import CharacterListPage from "./pages/CharacterListPage";
import CharacterEditPage from "./pages/CharacterEditPage";

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
          <Route path="/stories/:id/edit" element={<EditStory />} />
          <Route path="/stories/:id" element={<StoryDetail />} />
          <Route path="/stories/players/:id" element={<StoryPlayers />} />
          <Route path="/stories/encounter/new/:id" element={<NewEncounter />} />
          <Route path="/stories/:story_id/encounter/edit/:enc_id" element={<EditEncounter />} />
          <Route path="/stories/:story_id/encounter/:enc_id" element={<EncounterToStage />} />
        </Route>
        <Route path="/tasks">
          <Route path="/tasks" element={<TasksPage />} />
          <Route path="/tasks/new" element={<NewTask />} />
          <Route path="/tasks/:id/edit" element={<EditTask />} />
        </Route>
        <Route path="/stages"  >
          <Route path="/stages" element={<StagesPage />} />
          <Route path="/stages/:id/story/:story" element={<StageDetail />} />
          <Route path="/stages/:id/story/:story/players" element={<PlayersPage />} />
          <Route path="/stages/:id/story/:story/next" element={<StageManageNextEncounter />} />
          <Route path="/stages/:id/story/:story/addnext" element={<StageNextEncounter />} />
          <Route path="/stages/:id/story/:story/encounter/:encounterid" element={<StageEncounterDetail />} />
          <Route path="/stages/:id/story/:story/encounter/:encounterid/players" element={<AddPlayerToStageEncounter />} />
          <Route path="/stages/:id/story/:story/encounter/:encounterid/npc/:storyteller_id" element={<AddNPCToStageEncounter />} />
          <Route path="/stages/:id/story/:story/encounter/:encounterid/task/:storyteller_id" element={<TaskToEncounter />} />
          <Route path="/stages/:id/story/:story/encounter/:encounterid/encounter" element={<NextEncounter />} />
          <Route path="/stages/start/:id" element={<StageStart />} />
        </Route>
        <Route path="/autoplay" >
          <Route path="/autoplay" element={<AutoPlayPage />} />
          <Route path="/autoplay/new" element={<AutoPlayAdd />} />
          <Route path="/autoplay/:id/story/:story" element={<AutoPlayDetail />} />
          <Route path="/autoplay/:id/story/:story/next" element={<AutoPlayNext />} />
        </Route>
        <Route path="/users"  >
          <Route path="/users" element={<UsersPage />} />
          <Route path="/users/:id" element={<UserAsStoryteller />} />
          <Route path="/users/player/:id" element={<UserAsPlayer />} />
        </Route>
        <Route path="/admin/sessions" element={<SessionMonitorPage />} />
        <Route path="/characters"  >
          <Route path="/characters" element={<CharacterListPage />} />
          <Route path="/characters/:id/edit" element={<CharacterEditPage />} />
        </Route>
        <Route path="*" element={<NoMatch />} />
      </Route>
    </Router>
  );
};

export default Routes;
