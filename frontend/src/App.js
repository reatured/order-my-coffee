import React, { useEffect, useState } from 'react';
import { HashRouter as Router, Routes, Route, useNavigate, useParams } from 'react-router-dom';
import './App.css';

const API_BASE = "https://order-coffee-production.up.railway.app";

function CoffeeCard({ coffee, onClick, quantity, setQuantity }) {
  return (
    <div className="coffee-card" onClick={onClick}>
      <img src={`${process.env.PUBLIC_URL}/images/${coffee.image}`} alt={coffee.name} />
      <h2>{coffee.name}</h2>
      <div className="quantity-controls">
        <button type="button" onClick={e => { e.stopPropagation(); setQuantity(Math.max(1, quantity - 1)); }}>-</button>
        <span>{quantity}</span>
        <button type="button" onClick={e => { e.stopPropagation(); setQuantity(quantity + 1); }}>+</button>
      </div>
    </div>
  );
}

function CoffeeList() {
  const [coffees, setCoffees] = useState([]);
  const [quantities, setQuantities] = useState({});
  const navigate = useNavigate();

  useEffect(() => {
    fetch(`${API_BASE}/coffees`)
      .then(res => res.json())
      .then(data => {
        setCoffees(data);
        // Set default quantity to 1 for each coffee
        const q = {};
        data.forEach(c => { q[c.id] = 1; });
        setQuantities(q);
      });
  }, []);

  return (
    <div className="coffee-list-container">
      <h1>Order Coffee</h1>
      <div className="coffee-cards">
        {coffees.map(coffee => (
          <CoffeeCard
            key={coffee.id}
            coffee={coffee}
            quantity={quantities[coffee.id] || 1}
            setQuantity={q => setQuantities({ ...quantities, [coffee.id]: q })}
            onClick={() => navigate(`/order/${coffee.id}?quantity=${quantities[coffee.id] || 1}`)}
          />
        ))}
      </div>
    </div>
  );
}

function OrderPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [coffee, setCoffee] = useState(null);
  const [quantity, setQuantity] = useState(1);
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [notes, setNotes] = useState('');
  const [message, setMessage] = useState('');

  useEffect(() => {
    fetch(`${API_BASE}/coffees`)
      .then(res => res.json())
      .then(data => {
        const c = data.find(c => String(c.id) === id);
        setCoffee(c);
      });
    // Get quantity from query param
    const params = new URLSearchParams(window.location.search);
    setQuantity(parseInt(params.get('quantity')) || 1);
  }, [id]);

  const handleSubmit = (e) => {
    e.preventDefault();
    fetch(`${API_BASE}/order`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name,
        coffeeId: parseInt(id),
        notes,
        email
      })
    })
      .then(res => res.json())
      .then(data => {
        setMessage('Order submitted! Redirecting...');
        setTimeout(() => navigate('/'), 1500);
      });
  };

  if (!coffee) return <div>Loading...</div>;

  return (
    <div className="order-form-container">
      <img src={`${process.env.PUBLIC_URL}/images/${coffee.image}`} alt={coffee.name} />
      <h2>{coffee.name}</h2>
      <p>Quantity: {quantity}</p>
      <form onSubmit={handleSubmit}>
        <input
          placeholder="Your name"
          value={name}
          onChange={e => setName(e.target.value)}
          required
        />
        <input
          placeholder="Your email"
          value={email}
          onChange={e => setEmail(e.target.value)}
          required
          type="email"
        />
        <input
          placeholder="Notes"
          value={notes}
          onChange={e => setNotes(e.target.value)}
        />
        <button type="submit">Send Order</button>
      </form>
      {message && <p className="success-message">{message}</p>}
    </div>
  );
}

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<CoffeeList />} />
        <Route path="/order/:id" element={<OrderPage />} />
      </Routes>
    </Router>
  );
}

export default App;
