import { Outlet } from "react-router";
import Menu from "@/components/layout/Menu";

import "@/assets/main.css";

function Layout() {
  return (
    <>
      <Menu />
      <main>
        <Outlet />
      </main>
    </>
  );
}

export default Layout;
