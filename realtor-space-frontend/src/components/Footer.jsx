const Footer = () => {
  return (
      <footer className="bg-gray-900 text-white pt-16 pb-8 px-5">
        <div className="max-w-7xl mx-auto">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-10 mb-10">
            {/* Logo and Info */}
            <div>
              <div className="flex items-center gap-3 mb-5">
                <img src="/logo2.png" alt="Logo" className="w-12" />
                <h1 className="text-xl font-bold">Realtor's Space</h1>
              </div>
              <p className="text-gray-400 mb-3">Premium property listings with luxury experience</p>
              <p className="text-gray-400">Nyeri, Kenya</p>
            </div>

            {/* Quick Links */}
            <div>
              <h3 className="text-lg font-semibold mb-5">Quick Links</h3>
              <ul className="space-y-2">
                <li><a href="/" className="text-gray-400 hover:text-white transition-colors">Home</a></li>
                <li><a href="/listings" className="text-gray-400 hover:text-white transition-colors">Listings</a></li>
                <li><a href="#about" className="text-gray-400 hover:text-white transition-colors">About Us</a></li>
                <li><a href="#contact" className="text-gray-400 hover:text-white transition-colors">Contact</a></li>
              </ul>
            </div>

            {/* Contact Info */}
            <div>
              <h3 className="text-lg font-semibold mb-5">Contact Us</h3>
              <ul className="space-y-3 text-gray-400">
                <li className="flex items-center gap-3">
                  <i className="fas fa-phone text-blue-500"></i>
                  <span>+254 757 577 018</span>
                </li>
                <li className="flex items-center gap-3">
                  <i className="fas fa-envelope text-blue-500"></i>
                  <span>realtorspace04@gmail.com</span>
                </li>
                <li className="flex items-center gap-3">
                  <i className="fas fa-map-marker-alt text-blue-500"></i>
                  <span>Nyeri Town, Kenya</span>
                </li>
              </ul>
            </div>
          </div>

          <div className="pt-8 border-t border-gray-800 text-center text-gray-500 text-sm">
            <p>&copy; 2025 Realtor's Space. All rights reserved.</p>
          </div>
        </div>
      </footer>
  )
}

export default Footer