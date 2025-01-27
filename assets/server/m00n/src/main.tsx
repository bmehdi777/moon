import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Route, Routes } from "react-router";
import Landing from "@/pages/Landing";
import KeycloakProvider from "@/contexts/KeycloakContext";
import NotFound from "@/pages/NotFound";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <KeycloakProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Landing />} />

          <Route path="*" element={<NotFound />} />
        </Routes>
      </BrowserRouter>
    </KeycloakProvider>
  </StrictMode>,
);
