import { Link, NavLink } from "react-router";
import Status from "@/components/layout/Status";

function Menu() {
  return (
    <>
      <nav>
        <Link to="/" className="logo">
          <div className="logo-circle"></div>
          <div className="logo-text">MOON</div>
        </Link>

				<Status />

        <div className="menu">
          <NavLink to="/dashboard" className="menu-item">
            Dashboard
          </NavLink>
          <NavLink to="/settings" className="menu-item">
            Settings
          </NavLink>
        </div>
      </nav>
    </>
  );
}

export default Menu;
