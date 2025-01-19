import { Outlet } from "react-router";
import Menu from "@/components/layout/Menu";

import "@/assets/main.css";
import { useEffect, useState } from "react";
import { HttpCallsContext } from "@/context/HttpCallContext";
import { HttpMessage } from "@/types/request.types";

function Layout() {
  const [statusOnline, setStatusOnline] = useState<boolean>(false);
  const [httpCalls, setHttpCalls] = useState<HttpMessage[]>([]);

  useEffect(() => {
    const eventSource: EventSource = new EventSource("/api/tunnels/status");

    eventSource.onopen = () => {
      setStatusOnline(true);
      console.log("Connection open");
    };
    eventSource.onerror = (err) => {
      setStatusOnline(false);
      console.log("Error : ", err);
    };
    eventSource.onmessage = (msg) => {
      setHttpCalls(JSON.parse(msg.data));
    };

    return () => eventSource.close();
  }, []);

  return (
    <HttpCallsContext.Provider
      value={{ httpCalls, setHttpCalls, statusOnline, setStatusOnline }}
    >
      <div className="main-container">
        <Menu />
        <main>
          <Outlet />
        </main>
      </div>
    </HttpCallsContext.Provider>
  );
}

export default Layout;
