import { NavLink } from "react-router";

function DashboardMenu() {
  return (
    <nav>
      <ul className="sidebar-items">
        <NavLink
          className={(data) => (data.isActive ? "active" : "")}
          to="/dashboard"
					end
        >
          <li>Overview</li>
        </NavLink>
        <NavLink
          className={(data) => (data.isActive ? "active" : "")}
          to="/dashboard/certificates"
					end
        >
          <li>Certificates</li>
        </NavLink>
        <NavLink
          className={(data) => (data.isActive ? "active" : "")}
          to="/dashboard/endpoints"
					end
        >
          <li>Endpoints</li>
        </NavLink>
      </ul>

      <div className="sidebar-footer">User info</div>
    </nav>
  );
}

export default DashboardMenu;
