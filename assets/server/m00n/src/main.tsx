import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Route, Routes } from "react-router";
import Landing from "@/pages/Landing";
import KeycloakProvider from "@/contexts/KeycloakContext";
import NotFound from "@/pages/NotFound";
import Overview from "@/pages/Dashboard/Overview";
import Endpoints from "@/pages/Dashboard/Endpoints";
import Certificates from "@/pages/Dashboard/Certificates";
import DashboardLayout from "@/components/layout/DashboardLayout";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <KeycloakProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Landing />} />

          <Route element={<DashboardLayout />}>
            <Route exact path="/dashboard" element={<Overview />} />
            <Route exact path="/dashboard/certificates" element={<Certificates />} />
            <Route exact path="/dashboard/endpoints" element={<Endpoints />} />
          </Route>

          <Route path="*" element={<NotFound />} />
        </Routes>
      </BrowserRouter>
    </KeycloakProvider>
  </StrictMode>,
);
