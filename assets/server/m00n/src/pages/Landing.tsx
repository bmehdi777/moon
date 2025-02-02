import "@/assets/landing.css";
import { useKeycloak } from "@/contexts/KeycloakContext";
import { Link } from "react-router";
import "@/assets/main.css";

function Landing() {
  const { authenticated, login, logout, register } = useKeycloak();

  const currentYear: number = new Date().getFullYear();

  return (
    <div id="landing" className="flex flex-col bg-white text-black min-h-lvh">
      <nav className="py-16 px-0">
        <div className="flex my-0 mx-auto py-0 px-24 max-w-[1200px] justify-between items-center">
          <div className="logo-section">
            <div className="moon-logo"></div>
            <div className="moon-name">Moon</div>
          </div>
          <div className="flex gap-6">
            {!authenticated ? (
              <>
                <button className="py-2 px-4 rounded-3xl text-base cursor-pointer transition-all duration-[0.3s] ease-in-out border-black border-solid border bg-transparent hover:bg-black hover:text-white" onClick={login}>
                  Login
                </button>
                <button className="py-2 px-4 rounded-3xl text-base cursor-pointer transition-all duration-[0.3s] ease-in-out border-black text-white bg-black hover:bg-[#333]" onClick={register}>
                  Register
                </button>
              </>
            ) : (
              <>
                <Link className="py-2 px-4 rounded-3xl text-base cursor-pointer transition-all duration-[0.3s] ease-in-out border-black border-solid border bg-transparent hover:bg-black hover:text-white" to="/dashboard">
                  Dashboard
                </Link>
                <button className="py-2 px-4 rounded-3xl text-base cursor-pointer transition-all duration-[0.3s] ease-in-out border-black text-white bg-black hover:bg-[#333]" onClick={logout}>
                  Logout
                </button>
              </>
            )}
          </div>
        </div>
      </nav>

      <main className="py-20 px-0 text-center flex-1">
        <div className="py-0 px-24 max-w-[800px] my-0 mx-auto">
          <h1>
            Your local services,
            <br />
            Available everywhere
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
    </div>
  );
}

export default Landing;
