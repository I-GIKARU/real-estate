import { BrowserRouter as Router, Routes, Route, useLocation } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import Header from './components/Header'
import Home from './pages/Home'
import Listings from './pages/Listings'
import PropertyDetails from './pages/PropertyDetails'
import About from './pages/About'
import Services from './pages/Services'
import Contact from './pages/Contact'
import Login from './pages/Login'
import Register from './pages/Register'
import PasswordResetRequest from './pages/auth/PasswordResetRequest'
import PasswordResetConfirm from './pages/auth/PasswordResetConfirm'
import AgentProperties from './pages/agent/AgentProperties'
import PropertyForm from './pages/agent/PropertyForm'
import AdminDashboard from './pages/AdminDashboard'
import Footer from './components/Footer'
import './index.css'

function AppContent() {
    const location = useLocation()
    const isAdminRoute = location.pathname.startsWith('/admin')

    return (
        <div className="flex flex-col min-h-screen">
            {!isAdminRoute && <Header />}
            <main className="flex-grow">
                <Routes>
                    <Route path="/" element={<Home />} />
                    <Route path="/listings" element={<Listings />} />
                    <Route path="/property/:id" element={<PropertyDetails />} />
                    <Route path="/about" element={<About />} />
                    <Route path="/services" element={<Services />} />
                    <Route path="/contact" element={<Contact />} />
                    <Route path="/login" element={<Login />} />
                    <Route path="/register" element={<Register />} />
                    <Route path="/reset-password" element={<PasswordResetRequest />} />
                    <Route path="/reset-password/confirm" element={<PasswordResetConfirm />} />
                    <Route path="/agent/properties" element={<AgentProperties />} />
                    <Route path="/agent/properties/new" element={<PropertyForm />} />
                    <Route path="/agent/properties/:id/edit" element={<PropertyForm />} />
                    <Route path="/admin/dashboard" element={<AdminDashboard />} />
                </Routes>
            </main>
            {!isAdminRoute && <Footer />}
        </div>
    )
}

function App() {
    return (
        <AuthProvider>
            <Router>
                <AppContent />
            </Router>
        </AuthProvider>
    )
}

export default App