import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Login from './pages/Login';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/dashboard" element={<div className="p-10">Bienvenido al dashboard 🔐</div>} />
      </Routes>
    </Router>
  );
}

export default App;
