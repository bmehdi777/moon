import { FormatToRaw } from "@/utils/format";
import { useState } from "react";

interface DetailCodeProps {
  title: string;
  content: Record<string, string> | string;
  disableToggleRaw?: boolean;
  emptyMessage?: string;
}

function DetailCode(props: DetailCodeProps) {
  const { title, content, disableToggleRaw = false, emptyMessage } = props;

  const [formatToggle, setFormatToggle] = useState<boolean>(false);

  return (
    <div className="detail-section">
      <h4>{title}</h4>
      {content === "" && emptyMessage ? (
        <span className="empty">{emptyMessage}</span>
      ) : (
        <>
          {!disableToggleRaw && (
            <div className="format-switch">
              <span className="format-label">JSON</span>
              <label className="switch">
                <input
                  type="checkbox"
                  className="format-toggle"
                  checked={formatToggle}
                  onChange={() => setFormatToggle(!formatToggle)}
                />
                <span className="slider"></span>
              </label>
              <span className="format-label">Raw</span>
            </div>
          )}
          <pre>
            {!formatToggle
              ? JSON.stringify(content, null, 2)
              : FormatToRaw(content as Record<string, string>)}
          </pre>
        </>
      )}
    </div>
  );
}

export default DetailCode;
