import "@/assets/status.css";
import { HttpCallsContext } from "@/context/HttpCallContext";
import { useContext } from "react";

function Status() {
  const { statusOnline } = useContext(HttpCallsContext);
  return (
    <div className={`status-container ${statusOnline ? "" : "offline"}`}>
      <div className="status-indicator"></div>
      <div className="status-text">
        System {statusOnline ? "online" : "offline"}
      </div>
    </div>
  );
}

export default Status;
