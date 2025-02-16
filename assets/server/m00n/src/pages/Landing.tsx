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
                <button
                  className="py-2 px-4 rounded-3xl text-base cursor-pointer transition-all duration-[0.3s] ease-in-out border-black border-solid border bg-transparent hover:bg-black hover:text-white"
                  onClick={login}
                >
                  Login
                </button>
                <button
                  className="py-2 px-4 rounded-3xl text-base cursor-pointer transition-all duration-[0.3s] ease-in-out border-black text-white bg-black hover:bg-[#333]"
                  onClick={register}
                >
                  Register
                </button>
              </>
            ) : (
              <>
                <Link
                  className="py-2 px-4 rounded-3xl text-base cursor-pointer transition-all duration-[0.3s] ease-in-out border-black border-solid border bg-transparent hover:bg-black hover:text-white"
                  to="/dashboard"
                >
                  Dashboard
                </Link>
                <button
                  className="py-2 px-4 rounded-3xl text-base cursor-pointer transition-all duration-[0.3s] ease-in-out border-black text-white bg-black hover:bg-[#333]"
                  onClick={logout}
                >
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
          <p className="text-xl text-[#666] mb-12 max-w-[600px] mx-auto">
            Moon creates secure tunnels from your machine to the internet,
            making your local services globally accessible in seconds.
          </p>

          <div className="max-w-[640px] my-0 mx-auto bg-[#f5f5f5] rounded-xl p-4 shadow-lg animate-[float_6s_ease-in-out_infinite]">
            <div className="flex gap-2 mb-4">
              <div className="w-3 h-3 rounded-4xl bg-[#ff5f56]"></div>
              <div className="w-3 h-3 rounded-4xl bg-[#ffbd2e]"></div>
              <div className="w-3 h-3 rounded-4xl bg-[#27c93f]"></div>
            </div>
            <div className="font-mono text-sm text-left">
              <p className="animate-delay-0 my-2 mx-0 border-r-[2px] border-transparent text-[#2ea043] animate-typing whitespace-nowrap overflow-hidden w-0">
                $ moon start http://localhost:3000
              </p>
              <p className="animate-delay-2 my-2 mx-0 border-r-[2px] border-transparent text-[#666] animate-typing-slow whitespace-nowrap overflow-hidden w-0">
                🌒 Establishing tunnel...
              </p>
              <p className="animate-delay-3-5 my-2 mx-0 border-r-[2px] border-transparent animate-typing-slow whitespace-nowrap overflow-hidden w-0">
                🌓 Tunnel created successfully!
              </p>
              <p className="animate-delay-5 my-2 mx-0 border-r-[2px] border-transparent animate-typing-slow whitespace-nowrap overflow-hidden w-0">
                🌕 Your service is now available at:
              </p>
              <p className="animate-delay-6-5 my-2 mx-0 border-transparent text-black animate-url whitespace-nowrap overflow-hidden w-0 border-r-[4px]">
                https://my-app.m00n.fr
              </p>
            </div>
          </div>
        </div>
      </main>

      <footer className="border-t-[#eee] py-8 px-0">
        <div className="max-w-[1200px] my-0 mx-auto py-0 px-6 flex justify-between items-center text-[#666]">
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
