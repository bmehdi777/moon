import Keycloak from "keycloak-js";
import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useState,
} from "react";

// TODO: put it in a public route on the server
const keycloakClient = new Keycloak({
  url: "http://localhost:8081",
  realm: "moon",
  clientId: "moon-public",
});

interface KeycloakContextType {
  client?: Keycloak;
  authenticated: boolean;
  login: () => void;
  logout: () => void;
  register: () => void;
}

const KeycloakContext = createContext<KeycloakContextType>({
  client: undefined,
  authenticated: false,
  login: () => {},
  logout: () => {},
  register: () => {},
});

export const useKeycloak = function (): Omit<KeycloakContextType, "client"> {
  const context = useContext(KeycloakContext);

  if (!context.client) {
    throw new Error(
      "Keycloak client has not been assigned to KeycloakProvider",
    );
  }

  const { authenticated, login, logout, register } = context;
  return {
    authenticated,
    login,
    logout,
    register,
  };
};

interface KeycloakProviderProps {
  children: ReactNode;
}

function KeycloakProvider(props: KeycloakProviderProps) {
  const [client] = useState<Keycloak>(keycloakClient);
  const [authenticated, setAuthenticated] = useState<boolean>(false);

  useEffect(() => {
    (async function () {
      if (!client.didInitialize) {
        const auth = await client.init({
          onLoad: "check-sso",
        });
        setAuthenticated(auth);
      }
    })();
  }, []);

  const login = async () => {
    await client.login();
  };

  const logout = async () => {
    await client.logout();
  };

  const register = async () => {
    await client.register();
  };

  return (
    <KeycloakContext.Provider
      value={{ client, authenticated, login, logout, register }}
    >
      {props.children}
    </KeycloakContext.Provider>
  );
}

export default KeycloakProvider;
