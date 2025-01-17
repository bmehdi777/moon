import "@/assets/request.css";
import { useEffect, useMemo, useState } from "react";
import Details from "./Details";
import { HttpMessage } from "@/types/request.types";

interface RequestLineProps {
  method: "GET" | "PUT" | "POST" | "DELETE" | string;
  endpoint: string;
  duration: string;
  status: number;
  timestamp: string;

  currentActiveLineId: number;
  setCurrentActiveLineId: React.Dispatch<React.SetStateAction<number>>;

  lineId: number;
}


function RequestLine(props: RequestLineProps) {
  const isActive = useMemo(
    () => props.currentActiveLineId === props.lineId,
    [props.currentActiveLineId, props.lineId],
  );

  const statusClassName = useMemo(() => {
    if (props.status >= 200 && props.status < 300) return "status-2xx";
    else if (props.status >= 300 && props.status < 400) return "status-3xx";
    else if (props.status >= 400 && props.status < 500) return "status-4xx";
    else if (props.status >= 500) return "status-5xx";
    else return "";
  }, [props.status]);

  return (
    <tr
      className={isActive ? "selected" : ""}
      onClick={() => props.setCurrentActiveLineId(props.lineId)}
    >
      <th className={`verb verb-${props.method.toLowerCase()}`}>
        {props.method}
      </th>
      <th className="endpoint">{props.endpoint}</th>
      <th className="duration">{props.duration}</th>
      <th className={`status ${statusClassName}`}>{props.status}</th>
      <th className="timestamp">{props.timestamp}</th>
    </tr>
  );
}

function Request() {
  const [activeLineId, setActiveLineId] = useState<number>(-1);
  const [httpCalls, setHttpCalls] = useState<HttpMessage[]>([]);

  useEffect(() => {
    const eventSource: EventSource = new EventSource("/api/tunnels/status");

    eventSource.onopen = () => console.log("Connection open");
    eventSource.onerror = (err) => console.log("Error : ", err);
    eventSource.onmessage = (msg) => {
      setHttpCalls(JSON.parse(msg.data));
      console.log("msg : ", msg);
    };

    return () => eventSource.close();
  }, []);

  return (
    <div className={`dashboard ${activeLineId !== -1 ? "selected" : ""}`}>
      <div className="card card-req">
        <div className="request-table-container">
          <table>
            <thead>
              <tr>
                <th>Method</th>
                <th>Endpoint</th>
                <th>Duration</th>
                <th>Status</th>
                <th>Timestamp</th>
              </tr>
            </thead>
            <tbody>
              {httpCalls.map((element, index) => (
                <RequestLine
                  lineId={index}
                  currentActiveLineId={activeLineId}
                  setCurrentActiveLineId={setActiveLineId}
                  method={element.request.method}
                  endpoint={element.request.path}
                  duration="0.5ms"
                  status={element.response.status}
                  timestamp="10-01-2025"
                />
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {activeLineId != -1 && (
        <Details selectedHttpMessage={httpCalls[activeLineId]} resetSelectedLine={() => setActiveLineId(-1)} />
      )}
    </div>
  );
}

export default Request;
