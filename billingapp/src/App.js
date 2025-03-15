import "./App.css";
import LoginPage from "./pages/LoginPage.js";
import Pizzatype from "./components/Pizzatype.js";
import Topping from "./components/Topping.js";
import Beverage from "./components/Beverage.js";

import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import Dashboard from "./layouts/Dashboard.js";

function App() {
  return (
    <>
      <Router>
        <Routes>
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/pizzatypes" element={<Pizzatype />} />
          <Route path="/toppings" element={<Topping />} />
          <Route path="/beverages" element={<Beverage />} />
        </Routes>
      </Router>
    </>
  );
}

export default App;
