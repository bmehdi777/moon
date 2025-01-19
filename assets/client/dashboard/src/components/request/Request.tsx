import "@/assets/request.css";
import { useEffect, useMemo, useState } from "react";
import Details from "./Details";
import { HttpMessage } from "@/types/request.types";

interface RequestLineProps {
  method: "GET" | "PUT" | "POST" | "DELETE" | string;
  endpoint: string;
  duration: string;
  status: number;
  datetime: string;

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
      <td className={`verb verb-${props.method.toLowerCase()}`}>
        {props.method}
      </td>
      <td className="endpoint">{props.endpoint}</td>
      <td className="duration">{props.duration}</td>
      <td className={`status ${statusClassName}`}>{props.status}</td>
      <td className="timestamp">{props.datetime}</td>
    </tr>
  );
}

function Request() {
  const [activeLineId, setActiveLineId] = useState<number>(-1);
  const [httpCalls, setHttpCalls] = useState<HttpMessage[]>([]);

  const [filter, setFilter] = useState<string>("");

  const filteredHttpCalls: HttpMessage[] = useMemo(() => {
    return httpCalls.filter((call) =>
      `${call.request.method} ${call.request.path} ${call.response.duration} ${call.response.status} ${call.request.datetime}`
        .toLowerCase()
        .includes(filter),
    );
  }, [filter, httpCalls]);

  useEffect(() => {
    const eventSource: EventSource = new EventSource("/api/tunnels/status");

    eventSource.onopen = () => console.log("Connection open");
    eventSource.onerror = (err) => console.log("Error : ", err);
    eventSource.onmessage = (msg) => {
      setHttpCalls(JSON.parse(msg.data));
    };

    return () => eventSource.close();
  }, []);

  return (
    <div className={`dashboard ${activeLineId !== -1 ? "selected" : ""}`}>
      <div className="card card-req">
        <div className="search-container">
          <input
            type="text"
            className="search-input"
            placeholder="Search by method, endpoint, or status..."
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            id="searchInput"
          />
        </div>
        <div className="request-table-container">
          <table>
            <thead>
              <tr>
                <th>Method</th>
                <th>Endpoint</th>
                <th>Duration</th>
                <th>Status</th>
                <th>Datetime</th>
              </tr>
            </thead>
            <tbody>
              {httpCalls.length > 0 ? (
                <>
                  {filteredHttpCalls.length > 0 ? (
                    filteredHttpCalls.map((element, index) => (
                      <RequestLine
                        lineId={index}
                        currentActiveLineId={activeLineId}
                        setCurrentActiveLineId={setActiveLineId}
                        method={element.request.method}
                        endpoint={element.request.path}
                        duration={element.response.duration}
                        status={element.response.status}
                        datetime={element.request.datetime}
                      />
                    ))
                  ) : (
                    <tr className="no-results">
                      <td colSpan={5}>No matching requests found.</td>
                    </tr>
                  )}
                </>
              ) : (
                <tr className="no-results">
                  <td colSpan={5}>No request here.</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      {activeLineId != -1 && (
        <Details
          selectedHttpMessage={httpCalls[activeLineId]}
          resetSelectedLine={() => setActiveLineId(-1)}
        />
      )}
    </div>
  );
}

export default Request;
