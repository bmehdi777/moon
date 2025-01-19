import { HttpMessage } from "@/types/request.types";
import { createContext } from "react";

interface HttpCallsContextType {
  httpCalls: HttpMessage[];
  statusOnline: boolean;
  setStatusOnline: (value: boolean) => void;
  setHttpCalls: (value: HttpMessage[]) => void;
}

export const HttpCallsContext = createContext<HttpCallsContextType>({
  httpCalls: [],
  statusOnline: false,
  setStatusOnline: () => {},
  setHttpCalls: () => {},
});
