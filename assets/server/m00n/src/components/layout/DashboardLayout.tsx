import { Outlet } from "react-router";
import DashboardMenu from "./DashboardMenu";
import "@/assets/dashboard.css";

function DashboardLayout() {
  return (
    <div id="dashboard">
      <DashboardMenu />
      <main >
        <Outlet />
      </main>
    </div>
  );
}

export default DashboardLayout;
