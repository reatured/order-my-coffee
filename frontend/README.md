# Order Your Coffee - Full Stack Coffee Ordering Application

A modern, responsive web application for ordering coffee drinks with email confirmation and admin notifications.

## 🌟 Features

### Frontend (React)
- **Beautiful UI**: Modern, coffee-themed design with a clean, white palette and Poppins font
- **Responsive Design**: Works perfectly on desktop (2x2 grid) and mobile (single column)
- **Interactive Cards**: Browse coffee drinks with quantity controls and hover effects
- **Order Form**: Easy-to-use form with name, email, and notes fields
- **Guest Ordering**: Anyone can order by just entering their name and email (no login required)
- **Optional Login/Register**: Users can create an account for a personalized experience
- **Auto-filled Order Form**: Logged-in users have their name and email auto-filled and email field hidden
- **Clickable Header Title**: Click the "Order Your Coffee" title to return to the home page
- **GitHub Pages Deployment**: Live at https://reatured.github.io/order-my-coffee

### Backend Integration
- **Order Processing**: Sends orders to backend API with all necessary details
- **Email Confirmation**: Customers receive confirmation emails (if email provided)
- **Admin Notifications**: Backend notifies admin of new orders
- **API Endpoints**: 
  - `GET /coffees` - Fetch available coffee drinks
  - `POST /order` - Submit new order with customer details
  - `POST /login` / `POST /register` / `POST /logout` / `GET /me` - User authentication (optional)

## 🛠️ Tech Stack

### Frontend
- **React 19** with Create React App
- **React Router** (HashRouter for GitHub Pages compatibility)
- **Custom CSS** with modern design principles
- **Responsive Grid Layout** using CSS Grid and Flexbox
- **GitHub Pages** for static hosting

### Backend
- **API Server** at https://order-coffee-production.up.railway.app
- **Email Service** for order confirmations and admin notifications
- **Order Management** with customer details and notes

## 📱 User Experience

1. **Browse**: View coffee drinks in an attractive card layout
2. **Select**: Choose quantity using +/- buttons on each card
3. **Order**: Click card to proceed to order form
4. **Details**: Enter name, email, and any special notes (or auto-filled if logged in)
5. **Confirm**: Submit order and receive email confirmation
6. **Admin**: Backend automatically notifies admin of new orders
7. **Login/Register (Optional)**: Users can log in for a personalized experience

## 🎨 Design Features

- **Clean White Palette**: Modern, minimal, and professional look
- **Card Hover Effects**: Subtle animations and shadow changes
- **1:1 Image Aspect Ratio**: Perfect square coffee images
- **Mobile-First**: Responsive design that works on all devices
- **Modern Typography**: Poppins font for all text

## 🚀 Getting Started

### Frontend Development
```bash
cd frontend
npm install
npm start
```

### Production Build
```bash
npm run build
npm run deploy  # Deploys to GitHub Pages
```

### Backend API
The backend is hosted on Railway and handles:
- Coffee menu data
- Order processing
- Email notifications
- Admin alerts

## 📁 Project Structure

```
order-my-coffee/
├── frontend/                 # React application
│   ├── public/
│   │   ├── images/          # Coffee drink images
│   │   └── index.html       # Main HTML template
│   ├── src/
│   │   ├── App.js           # Main React component
│   │   ├── App.css          # Styling
│   │   └── index.js         # App entry point
│   └── package.json         # Dependencies and scripts
└── README.md                # This file
```

## 🔧 Configuration

### Environment Variables
- `API_BASE`: Backend API URL (currently Railway deployment)
- `PUBLIC_URL`: GitHub Pages base path for assets

### Dependencies
- `react-router-dom`: Client-side routing
- `gh-pages`: GitHub Pages deployment

## 🌐 Live Demo

**Frontend**: https://reatured.github.io/order-my-coffee

## 📧 Order Flow

1. Customer selects coffee and quantity
2. Fills out order form with personal details (or auto-filled if logged in)
3. Frontend sends POST request to `/order` endpoint
4. Backend processes order and sends emails:
   - Confirmation to customer (if email provided)
   - Notification to admin
5. Customer receives success message and is redirected

## 🎯 Future Enhancements

- [ ] Add more coffee varieties
- [ ] Implement order history
- [ ] Add payment processing
- [ ] Real-time order tracking
- [ ] Customer account system
- [ ] Admin dashboard

---

**Order Your Coffee** - A complete full-stack coffee ordering experience! ☕
