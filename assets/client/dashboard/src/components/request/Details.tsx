import {
  HttpMessage,
  HttpMessageRequest,
  HttpMessageResponse,
} from "@/types/request.types";
import { useEffect, useState } from "react";
import DetailCode from "./DetailCode";

interface DetailsProps {
  resetSelectedLine: () => void;
  selectedHttpMessage: HttpMessage;
}

function Details(props: DetailsProps) {
  const { selectedHttpMessage, resetSelectedLine } = props;

  const [currentTab, setCurrentTab] = useState<"request" | "response">(
    "request",
  );
  const [currentMessage, setCurrentMessage] = useState<
    HttpMessageRequest | HttpMessageResponse
  >(selectedHttpMessage.request);

  useEffect(() => {
    setCurrentMessage(selectedHttpMessage[currentTab]);
  }, [currentTab]);

  return (
    <div className="card">
      <div className="details">
        <div className="details-header">
          <div>
            <span className="verb verb-get">
              {selectedHttpMessage.request.method}{" "}
            </span>
            <span className="endpoint">{selectedHttpMessage.request.path}</span>
          </div>
          <button className="close-button" onClick={resetSelectedLine}>
            ✕
          </button>
        </div>

        <div className="tabs">
          <button
            className={`tab ${currentTab === "request" ? "active" : ""}`}
            onClick={() => setCurrentTab("request")}
          >
            Request
          </button>
          <button
            className={`tab ${currentTab === "response" ? "active" : ""}`}
            onClick={() => setCurrentTab("response")}
          >
            Response
          </button>
        </div>

        <DetailCode title="Headers" content={currentMessage.headers} />
        <DetailCode
          title="Body"
          content={currentMessage.body}
          emptyMessage="No data in body."
        />
      </div>
    </div>
  );
}

export default Details;
