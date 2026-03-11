import { BrowserRouter, Routes, Route } from 'react-router-dom';
import SolutionsPage from './pages/SolutionsPage';
import ArchitecturePage from './pages/ArchitecturePage';
import Layout from './components/Layout';

export default function App() {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/architecture" element={<ArchitecturePage />} />
          <Route path="/" element={<SolutionsPage />} />
  
        </Routes>
      </Layout>
    </BrowserRouter>
  );
}
