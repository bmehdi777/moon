import { Outlet } from "react-router";
import Menu from "@/components/layout/Menu";

import "@/assets/main.css";

function Layout() {
  return (
    <div className="main-container">
      <Menu />
      <main>
        <Outlet />
      </main>
    </div>
  );
}

export default Layout;
