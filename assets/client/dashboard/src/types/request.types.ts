export interface HttpMessage {
  request: {
    method: string;
    path: string;
    headers: Record<string, string>;
    body: string;
  };
  response: {
    status: number;
    headers: Record<string, string>;
    body: string;
  };
}
