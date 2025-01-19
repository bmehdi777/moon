export interface HttpMessageRequest {
  method: string;
  path: string;
  headers: Record<string, string>;
  body: string;
	datetime: string;
}
export interface HttpMessageResponse {
  status: number;
  headers: Record<string, string>;
  body: string;
  duration: string;
}
export interface HttpMessage {
  request: HttpMessageRequest;
  response: HttpMessageResponse;
}
