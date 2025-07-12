import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

const Header = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const { isAuthenticated, user, logout, isAgent } = useAuth()
  const navigate = useNavigate()

  const toggleMenu = () => {
    setIsMenuOpen(!isMenuOpen)
  }

  const handleLogout = () => {
    logout()
    navigate('/')
    setIsMenuOpen(false)
  }

  return (
      <header className="fixed w-full top-0 z-50 bg-white bg-opacity-90 shadow-md">
        <nav className="flex justify-between items-center px-5 py-5 max-w-7xl mx-auto">
          <div className="flex items-center gap-2">
            <h1 className="text-xl font-bold">Realtor's</h1>
            <img src="/logo2.png" alt="Logo" className="w-12" />
            <h1 className="text-xl font-bold">Space</h1>
          </div>

          <ul className={`md:flex gap-8 ${isMenuOpen ? 'flex flex-col fixed top-20 left-0 w-full h-[calc(100vh-80px)] bg-white items-center justify-center' : 'hidden'}`}>
            <li><Link to="/" className="font-medium hover:text-blue-600 transition-colors" onClick={() => setIsMenuOpen(false)}>Home</Link></li>
            <li><Link to="/listings" className="font-medium hover:text-blue-600 transition-colors" onClick={() => setIsMenuOpen(false)}>Listings</Link></li>
            {isAgent() && (
              <li><Link to="/agent/properties" className="font-medium hover:text-blue-600 transition-colors" onClick={() => setIsMenuOpen(false)}>My Properties</Link></li>
            )}
            <li><Link to="/about" className="font-medium hover:text-blue-600 transition-colors" onClick={() => setIsMenuOpen(false)}>About</Link></li>
            <li><Link to="/services" className="font-medium hover:text-blue-600 transition-colors" onClick={() => setIsMenuOpen(false)}>Services</Link></li>
            <li><Link to="/contact" className="font-medium hover:text-blue-600 transition-colors" onClick={() => setIsMenuOpen(false)}>Contact</Link></li>
            
            {!isAuthenticated ? (
              <>
                <li><Link to="/login" className="font-medium hover:text-blue-600 transition-colors" onClick={() => setIsMenuOpen(false)}>Login</Link></li>
                <li><Link to="/register" className="font-medium hover:text-blue-600 transition-colors" onClick={() => setIsMenuOpen(false)}>Register</Link></li>
              </>
            ) : (
              <>
                <li className="flex items-center gap-2">
                  <span className="text-sm text-gray-600">Welcome, {user?.first_name}</span>
                  <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded-full">{user?.user_type}</span>
                </li>
                <li>
                  <button 
                    onClick={handleLogout}
                    className="font-medium hover:text-red-600 transition-colors bg-red-50 hover:bg-red-100 px-3 py-1 rounded-md"
                  >
                    Logout
                  </button>
                </li>
              </>
            )}
          </ul>

          <div
              className="md:hidden cursor-pointer"
              onClick={toggleMenu}
          >
            <div className={`w-6 h-0.5 bg-black my-1.5 transition-all ${isMenuOpen ? 'transform rotate-45 translate-y-2' : ''}`}></div>
            <div className={`w-6 h-0.5 bg-black my-1.5 transition-all ${isMenuOpen ? 'opacity-0' : ''}`}></div>
            <div className={`w-6 h-0.5 bg-black my-1.5 transition-all ${isMenuOpen ? 'transform -rotate-45 -translate-y-2' : ''}`}></div>
          </div>
        </nav>
      </header>
  )
}

export default Header