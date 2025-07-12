import { useEffect } from 'react'

const About = () => {
  useEffect(() => {
    window.scrollTo(0, 0)
  }, [])

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-50 via-blue-50 to-pink-50">
      {/* Hero Section */}
      <section 
        className="h-[40vh] bg-gradient-to-r from-purple-600/90 to-blue-600/90 bg-cover bg-center flex flex-col justify-center relative overflow-hidden"
        style={{ backgroundImage: "url('https://images.unsplash.com/photo-1560518883-ce09059eeffa?ixlib=rb-1.2.1&auto=format&fit=crop&w=1350&q=80')" }}
      >
        {/* Floating shapes */}
        <div className="absolute inset-0 overflow-hidden">
          <div className="absolute top-20 left-10 w-20 h-20 bg-white/10 rounded-full blur-xl animate-pulse"></div>
          <div className="absolute bottom-20 right-10 w-32 h-32 bg-pink-400/20 rounded-full blur-xl animate-pulse delay-1000"></div>
          <div className="absolute top-1/2 left-1/2 w-24 h-24 bg-blue-400/20 rounded-full blur-xl animate-pulse delay-500"></div>
        </div>
        
        <div className="max-w-4xl mx-auto text-center text-white relative z-10 px-5 mt-20">
          <h1 className="text-5xl md:text-6xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-white to-blue-100">
            About Realtor's Space
          </h1>
          <p className="text-xl md:text-2xl text-blue-100">
            Your trusted partner in premium real estate
          </p>
        </div>
      </section>

      {/* Main Content */}
      <section className="py-20 px-5">
        <div className="max-w-7xl mx-auto">
          {/* Story Section */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center mb-20">
            <div className="order-2 lg:order-1">
              <div className="bg-white/70 backdrop-blur-sm rounded-2xl p-8 shadow-xl border border-white/20">
                <h2 className="text-4xl font-bold mb-6 bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent">
                  Our Story
                </h2>
                <p className="text-lg text-gray-700 leading-relaxed mb-6">
                  Based in the heart of Nyeri, Kenya, Realtor's Space was born from a vision to revolutionize 
                  the real estate experience. We believe that finding your perfect home should be an exciting 
                  journey, not a stressful ordeal.
                </p>
                <p className="text-lg text-gray-700 leading-relaxed mb-6">
                  Our platform combines cutting-edge technology with local expertise to bring you the finest 
                  property listings in Kenya. We're not just about selling properties; we're about helping 
                  you find your place in the world.
                </p>
                <div className="flex items-center gap-4 mt-8">
                  <div className="w-16 h-16 bg-gradient-to-r from-purple-500 to-blue-500 rounded-full flex items-center justify-center">
                    <i className="fas fa-heart text-white text-xl"></i>
                  </div>
                  <div>
                    <h4 className="font-semibold text-lg">Our Mission</h4>
                    <p className="text-gray-600">Making premium real estate accessible to everyone</p>
                  </div>
                </div>
              </div>
            </div>
            
            <div className="order-1 lg:order-2">
              <div className="relative">
                <img
                  src="https://images.unsplash.com/photo-1560518883-ce09059eeffa?ixlib=rb-1.2.1&auto=format&fit=crop&w=500&q=60"
                  alt="Luxury building"
                  className="w-full rounded-2xl shadow-2xl"
                />
                <div className="absolute -bottom-4 -right-4 w-full h-full bg-gradient-to-r from-purple-400/20 to-blue-400/20 rounded-2xl -z-10"></div>
              </div>
            </div>
          </div>

          {/* Stats Section */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-20">
            <div className="text-center bg-white/60 backdrop-blur-sm rounded-2xl p-8 shadow-lg border border-white/20 hover:shadow-xl transition-all hover:scale-105">
              <div className="w-16 h-16 bg-gradient-to-r from-purple-500 to-blue-500 rounded-full flex items-center justify-center mx-auto mb-4">
                <i className="fas fa-home text-white text-xl"></i>
              </div>
              <h3 className="text-4xl font-bold bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent mb-2">40+</h3>
              <p className="text-gray-700 font-medium">Properties Listed</p>
            </div>
            
            <div className="text-center bg-white/60 backdrop-blur-sm rounded-2xl p-8 shadow-lg border border-white/20 hover:shadow-xl transition-all hover:scale-105">
              <div className="w-16 h-16 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full flex items-center justify-center mx-auto mb-4">
                <i className="fas fa-smile text-white text-xl"></i>
              </div>
              <h3 className="text-4xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent mb-2">98%</h3>
              <p className="text-gray-700 font-medium">Client Satisfaction</p>
            </div>
            
            <div className="text-center bg-white/60 backdrop-blur-sm rounded-2xl p-8 shadow-lg border border-white/20 hover:shadow-xl transition-all hover:scale-105">
              <div className="w-16 h-16 bg-gradient-to-r from-pink-500 to-purple-500 rounded-full flex items-center justify-center mx-auto mb-4">
                <i className="fas fa-calendar-alt text-white text-xl"></i>
              </div>
              <h3 className="text-4xl font-bold bg-gradient-to-r from-pink-600 to-purple-600 bg-clip-text text-transparent mb-2">2</h3>
              <p className="text-gray-700 font-medium">Years Experience</p>
            </div>
          </div>

          {/* Team Section */}
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold mb-6 bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent">
              Meet Our Team
            </h2>
            <p className="text-lg text-gray-700 max-w-2xl mx-auto">
              Our dedicated team of real estate professionals is here to guide you through every step of your property journey.
            </p>
          </div>

          {/* Values Section */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            <div className="bg-white/60 backdrop-blur-sm rounded-2xl p-8 shadow-lg border border-white/20 hover:shadow-xl transition-all hover:scale-105">
              <div className="w-16 h-16 bg-gradient-to-r from-purple-500 to-pink-500 rounded-full flex items-center justify-center mb-6">
                <i className="fas fa-shield-alt text-white text-xl"></i>
              </div>
              <h3 className="text-xl font-bold mb-4 text-gray-800">Trust & Transparency</h3>
              <p className="text-gray-600">
                We believe in complete transparency in all our dealings. No hidden fees, no surprises - just honest, straightforward service.
              </p>
            </div>
            
            <div className="bg-white/60 backdrop-blur-sm rounded-2xl p-8 shadow-lg border border-white/20 hover:shadow-xl transition-all hover:scale-105">
              <div className="w-16 h-16 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full flex items-center justify-center mb-6">
                <i className="fas fa-rocket text-white text-xl"></i>
              </div>
              <h3 className="text-xl font-bold mb-4 text-gray-800">Innovation</h3>
              <p className="text-gray-600">
                We leverage cutting-edge technology to provide you with the most advanced property search and viewing experience.
              </p>
            </div>
            
            <div className="bg-white/60 backdrop-blur-sm rounded-2xl p-8 shadow-lg border border-white/20 hover:shadow-xl transition-all hover:scale-105">
              <div className="w-16 h-16 bg-gradient-to-r from-pink-500 to-blue-500 rounded-full flex items-center justify-center mb-6">
                <i className="fas fa-users text-white text-xl"></i>
              </div>
              <h3 className="text-xl font-bold mb-4 text-gray-800">Customer First</h3>
              <p className="text-gray-600">
                Your satisfaction is our priority. We go above and beyond to ensure you find the perfect property for your needs.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* WhatsApp Float Button */}
      <a
        href="https://wa.me/254757577018"
        className="fixed bottom-6 right-6 bg-green-500 text-white w-16 h-16 rounded-full flex items-center justify-center text-2xl shadow-lg hover:bg-green-600 transition-all hover:scale-110 z-50"
        target="_blank"
        rel="noopener noreferrer"
      >
        <i className="fab fa-whatsapp"></i>
      </a>
    </div>
  )
}

export default About
