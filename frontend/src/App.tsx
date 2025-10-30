import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Index from "./pages/HomePage";
import Checker from "./pages/CheckRepo";
import ResultDetail from "./pages/ResultDetail";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Index />} />
        <Route path="/checker" element={<Checker />} />
        <Route path="/result" element={<ResultDetail />} />
      </Routes>
    </BrowserRouter>
  </QueryClientProvider>
);

export default App;
