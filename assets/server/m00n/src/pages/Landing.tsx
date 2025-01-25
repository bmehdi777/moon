import "@/assets/landing.css";
import { useKeycloak } from "@/contexts/KeycloakContext";

function Landing() {
  const { authenticated, login, logout, register } =
    useKeycloak();

  const currentYear: number = new Date().getFullYear();

  return (
    <>
      <nav>
        <div className="container nav-content">
          <div className="logo-section">
            <div className="moon-logo"></div>
            <div className="moon-name">Moon</div>
          </div>
          <div className="auth-buttons">
            {!authenticated ? (
              <>
                <button className="button primary-btn" onClick={login}>
                  Login
                </button>
                <button className="button secondary-btn" onClick={register}>
                  Register
                </button>
              </>
            ) : (
              <button className="button secondary-btn" onClick={logout}>
                Logout
              </button>
            )}
          </div>
        </div>
      </nav>

      <main>
        <div className="container main-content">
          <h1>
            Your Local Services,
            <br />
            Available Worldwide
          </h1>
          <p className="main-description">
            Moon creates secure tunnels from your machine to the internet,
            making your local services globally accessible in seconds.
          </p>

          <div className="terminal">
            <div className="terminal-controls">
              <div className="control-dot red"></div>
              <div className="control-dot yellow"></div>
              <div className="control-dot green"></div>
            </div>
            <div className="terminal-content">
              <p className="command">$ moon start http://localhost:3000</p>
              <p className="output">🌒 Establishing tunnel...</p>
              <p className="output">🌓 Tunnel created successfully!</p>
              <p className="output">🌕 Your service is now available at:</p>
              <p className="url">https://my-app.m00n.fr</p>
            </div>
          </div>
        </div>
      </main>

      <footer>
        <div className="container footer-content">
          <div>© {currentYear} Moon. All rights reserved.</div>
          <div>
            <a
              target="_blank"
              href="https://github.com/bmehdi777/moon"
              className="github-link"
            >
              GitHub
            </a>
          </div>
        </div>
      </footer>
    </>
  );
}

export default Landing;
