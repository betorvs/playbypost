import { createContext, ReactNode, useState } from "react";
import CleanSession from "./CleanSession";
import CheckSession from "./CheckSession";

type Props = {
  children?: ReactNode;
};

interface AuthContextData {
  authenticated: boolean;
  setAuthenticated: (newState: boolean) => void;
  Logoff(): void;
}

const initialValue = {
  authenticated: false,
  setAuthenticated: () => {},
  Logoff: () => {},
};

const AuthContext = createContext<AuthContextData>(initialValue);

const AuthProvider = ({ children }: Props) => {
  const amIAuth = CheckSession();

  //Initializing an auth state with false value (unauthenticated)
  const [authenticated, setAuthenticated] = useState(amIAuth);

  function Logoff() {
    setAuthenticated(false);
    CleanSession();
  }

  return (
    <AuthContext.Provider value={{ authenticated, setAuthenticated, Logoff }}>
      {children}
    </AuthContext.Provider>
  );
};

export { AuthContext, AuthProvider };
