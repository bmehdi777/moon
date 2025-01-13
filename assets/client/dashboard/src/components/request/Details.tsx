interface DetailsProps {
  resetSelectedLine: () => void;
}

function Details(props: DetailsProps) {
  return (
    <div className="card">
      <div className="details">
        <div className="details-header">
          <div>
            <span className="verb verb-get">GET</span>
            <span className="endpoint">/api/users</span>
          </div>
          <button className="close-button" onClick={props.resetSelectedLine}>
            ✕
          </button>
        </div>

        <div>
          <h4>Headers</h4>
          <pre>
            {JSON.stringify({
              "Content-Type": "application/json",
              Authorization: "Bearer xxx",
            })}
          </pre>
        </div>

        <div>
          <h4>Response</h4>
          <pre>
            {JSON.stringify({
              status: 200,
              body: {
                success: true,
                data: [],
              },
            })}
          </pre>
        </div>
      </div>
    </div>
  );
}

export default Details;
