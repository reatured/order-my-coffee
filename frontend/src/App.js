import React, { useEffect, useState } from 'react';
import { HashRouter as Router, Routes, Route, useNavigate, useParams, useLocation } from 'react-router-dom';
import './App.css';

const API_BASE = "https://order-coffee-production.up.railway.app";

// Login Page Component
function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const navigate = useNavigate();

  const handleSubmit = (e) => {
    e.preventDefault();
    fetch(`${API_BASE}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ username, password })
    })
    .then(res => res.json())
    .then(data => {
      if (data.status === 'ok') {
        setMessage('Login successful!');
        setTimeout(() => navigate('/'), 1000);
      } else {
        setMessage(data.error || 'Login failed');
      }
    })
    .catch(err => setMessage('Login failed'));
  };

  return (
    <div className="order-form-container">
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        <input
          placeholder="Username"
          value={username}
          onChange={e => setUsername(e.target.value)}
          required
        />
        <input
          placeholder="Password"
          type="password"
          value={password}
          onChange={e => setPassword(e.target.value)}
          required
        />
        <button type="submit">Login</button>
      </form>
      {message && <p className="success-message">{message}</p>}
      <p>Don't have an account? <button onClick={() => navigate('/register')} style={{ background: 'none', border: 'none', color: '#c7a17a', textDecoration: 'underline', cursor: 'pointer' }}>Register</button></p>
    </div>
  );
}

// Register Page Component
function RegisterPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [message, setMessage] = useState('');
  const navigate = useNavigate();

  const handleSubmit = (e) => {
    e.preventDefault();
    if (password !== confirmPassword) {
      setMessage('Passwords do not match');
      return;
    }
    fetch(`${API_BASE}/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ username, password })
    })
    .then(res => res.json())
    .then(data => {
      if (data.status === 'ok') {
        setMessage('Registration successful!');
        setTimeout(() => navigate('/'), 1000);
      } else {
        setMessage(data.error || 'Registration failed');
      }
    })
    .catch(err => setMessage('Registration failed'));
  };

  return (
    <div className="order-form-container">
      <h2>Register</h2>
      <form onSubmit={handleSubmit}>
        <input
          placeholder="Username"
          value={username}
          onChange={e => setUsername(e.target.value)}
          required
        />
        <input
          placeholder="Password"
          type="password"
          value={password}
          onChange={e => setPassword(e.target.value)}
          required
        />
        <input
          placeholder="Confirm Password"
          type="password"
          value={confirmPassword}
          onChange={e => setConfirmPassword(e.target.value)}
          required
        />
        <button type="submit">Register</button>
      </form>
      {message && <p className="success-message">{message}</p>}
      <p>Already have an account? <button onClick={() => navigate('/login')} style={{ background: 'none', border: 'none', color: '#c7a17a', textDecoration: 'underline', cursor: 'pointer' }}>Login</button></p>
    </div>
  );
}

// Header Component with Login/Logout
function Header({ user, onLogout }) {
  const navigate = useNavigate();

  const handleLogout = () => {
    fetch(`${API_BASE}/logout`, {
      method: 'POST',
      credentials: 'include'
    })
    .then(() => {
      onLogout();
      navigate('/');
    });
  };

  return (
    <div style={{ background: '#ffffff', padding: '1rem', borderBottom: '1px solid #e0e0e0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
      <div style={{ flex: 1 }}></div>
      <h1
        style={{ margin: 0, color: '#333333', fontFamily: 'Poppins, sans-serif', fontWeight: '600', fontSize: '1.8rem', cursor: 'pointer' }}
        onClick={() => navigate('/')}
        title="Go to Home"
      >
        Order Your Coffee
      </h1>
      <div style={{ flex: 1, display: 'flex', justifyContent: 'flex-end' }}>
        {user ? (
          <span>
            Welcome, {user.username}! 
            <button onClick={handleLogout} style={{ marginLeft: '1rem', background: '#666666', color: 'white', border: 'none', padding: '0.5rem 1rem', borderRadius: '4px', cursor: 'pointer' }}>
              Logout
            </button>
          </span>
        ) : (
          <span>
            <button onClick={() => navigate('/login')} style={{ background: '#666666', color: 'white', border: 'none', padding: '0.5rem 1rem', borderRadius: '4px', cursor: 'pointer', marginRight: '0.5rem' }}>
              Login
            </button>
            <button onClick={() => navigate('/register')} style={{ background: '#f0f0f0', color: '#333333', border: '1px solid #ddd', padding: '0.5rem 1rem', borderRadius: '4px', cursor: 'pointer' }}>
              Register
            </button>
          </span>
        )}
      </div>
    </div>
  );
}

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

function CoffeeList({ user }) {
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

function OrderPage({ user }) {
  const { id } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const [coffee, setCoffee] = useState(null);
  const [quantity, setQuantity] = useState(1);
  const [name, setName] = useState(user ? user.username : '');
  const [email, setEmail] = useState(user ? user.email || '' : '');
  const [notes, setNotes] = useState('');
  const [message, setMessage] = useState('');

  useEffect(() => {
    fetch(`${API_BASE}/coffees`)
      .then(res => res.json())
      .then(data => {
        const c = data.find(c => String(c.id) === id);
        setCoffee(c);
      });
  }, [id]);

  useEffect(() => {
    // Get quantity from query param whenever location changes
    const params = new URLSearchParams(location.search);
    const quantityFromUrl = parseInt(params.get('quantity')) || 1;
    setQuantity(quantityFromUrl);
  }, [location.search]);

  useEffect(() => {
    // Autofill name and email if user is logged in
    if (user) {
      setName(user.username);
      setEmail(user.email || '');
    }
  }, [user]);

  const handleSubmit = (e) => {
    e.preventDefault();
    fetch(`${API_BASE}/order`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        name,
        coffeeId: parseInt(id),
        quantity,
        notes,
        email
      })
    })
      .then(res => res.json())
      .then(data => {
        setMessage('Order submitted! Back to home...');
        setTimeout(() => navigate('/'), 1500);
      });
  };

  if (!coffee) return <div>Loading...</div>;

  return (
    <div className="order-form-container">
      <img src={`${process.env.PUBLIC_URL}/images/${coffee.image}`} alt={coffee.name} />
      <h2>{coffee.name}</h2>
      <div className="quantity-controls" style={{ marginBottom: '16px' }}>
        <button type="button" onClick={() => setQuantity(Math.max(1, quantity - 1))}>-</button>
        <span>Quantity: {quantity}</span>
        <button type="button" onClick={() => setQuantity(quantity + 1)}>+</button>
      </div>
      <form onSubmit={handleSubmit}>
        <input
          placeholder="Your name"
          value={name}
          onChange={e => setName(e.target.value)}
          required
          disabled={!!user}
        />
        {!user && (
          <input
            placeholder="Your email"
            value={email}
            onChange={e => setEmail(e.target.value)}
            required
            type="email"
          />
        )}
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
  const [user, setUser] = useState(null);

  useEffect(() => {
    // Check if user is logged in on app load
    fetch(`${API_BASE}/me`, {
      credentials: 'include'
    })
    .then(res => res.json())
    .then(data => {
      if (data.status === 'ok') {
        setUser(data.user);
      }
    })
    .catch(err => console.log('Not logged in'));
  }, []);

  return (
    <Router>
      <Header user={user} onLogout={() => setUser(null)} />
      <Routes>
        <Route path="/" element={<CoffeeList user={user} />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/order/:id" element={<OrderPage user={user} />} />
      </Routes>
    </Router>
  );
}

export default App;
