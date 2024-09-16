// import { useState } from "react";
// import { RouterProvider, createBrowserRouter } from "react-router-dom";
// import UsersPage from "./pages/Users";
// import PlayersPage from "./pages/Players";
// import StoriesPage from "./pages/Stories";
// import HomePublicPage from "./pages/Home";
// import Layout from "./components/Layout";
// import SaveToken from "./SaveToken";
// import CheckSession from "./CheckSession";
// import CleanSession from "./CleanSession";
// import StoryDetail from "./pages/StoryDetail";

// function createRouter() {
//   let iam = CheckSession();

//   const [isLoggedIn, setLogIn] = useState<boolean>(iam);

//   const router = createBrowserRouter([
//     {
//       path: "/",
//       element: (
//         <Layout
//           loggedIn={isLoggedIn}
//           setLoggedIn={setLogIn}
//           setToken={SaveToken}
//           logoff={CleanSession}
//         />
//       ),
//       children: [
//         {
//           path: "/",
//           element: <HomePublicPage />,
//         },
//         {
//           path: "/users",
//           element: <UsersPage loggedIn={isLoggedIn} />,
//         },
//         {
//           path: "/stories",
//           element: <StoriesPage loggedIn={isLoggedIn} />,
//         },
//         {
//           path: "/stories/:id",
//           element: <StoryDetail loggedIn={isLoggedIn} />,
//         },
//         {
//           path: "/stories/players/:id",
//           element: <PlayersPage loggedIn={isLoggedIn} />,
//         },
//       ],
//     },
//   ]);
//   return router;
// }

import { BrowserRouter } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import Routes from "./Routes";

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes />
      </AuthProvider>
    </BrowserRouter>
  );
}
