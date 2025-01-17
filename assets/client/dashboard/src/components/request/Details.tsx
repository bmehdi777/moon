import { HttpMessage } from "@/types/request.types";

interface DetailsProps {
  resetSelectedLine: () => void;
  selectedHttpMessage: HttpMessage;
}

function Details(props: DetailsProps) {
  const { selectedHttpMessage, resetSelectedLine } = props;

  return (
    <div className="card">
      <div className="details">
        <div className="details-header">
          <div>
            <span className="verb verb-get">{selectedHttpMessage.request.method}</span>
            <span className="endpoint">{selectedHttpMessage.request.path}</span>
          </div>
          <button className="close-button" onClick={resetSelectedLine}>
            ✕
          </button>
        </div>

        <div>
          <h4>Headers</h4>
          <pre>{JSON.stringify(selectedHttpMessage.request.headers, null, 2)}</pre>
        </div>

        <div>
          <h4>Response</h4>
          <pre>{JSON.stringify(selectedHttpMessage.response.headers, null, 2)}</pre>
        </div>
      </div>
    </div>
  );
}

export default Details;
