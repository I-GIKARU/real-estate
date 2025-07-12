import { useEffect } from 'react'

const Services = () => {
  useEffect(() => {
    window.scrollTo(0, 0)
  }, [])

  const services = [
    {
      icon: 'fas fa-home',
      title: 'Property Sales',
      description: 'Expert guidance for buying and selling residential and commercial properties with maximum value.',
      features: ['Market Analysis', 'Property Valuation', 'Negotiation Support', 'Legal Assistance']
    },
    {
      icon: 'fas fa-key',
      title: 'Rental Services',
      description: 'Comprehensive rental management services for landlords and tenants.',
      features: ['Tenant Screening', 'Property Management', 'Maintenance Coordination', 'Rent Collection']
    },
    {
      icon: 'fas fa-chart-line',
      title: 'Investment Advisory',
      description: 'Strategic real estate investment advice to maximize your portfolio returns.',
      features: ['ROI Analysis', 'Market Research', 'Risk Assessment', 'Portfolio Management']
    },
    {
      icon: 'fas fa-search',
      title: 'Property Consultation',
      description: 'Personalized consultation services to help you make informed property decisions.',
      features: ['Site Inspection', 'Due Diligence', 'Market Insights', 'Custom Solutions']
    },
    {
      icon: 'fas fa-hammer',
      title: 'Property Development',
      description: 'End-to-end property development services from planning to completion.',
      features: ['Project Planning', 'Construction Management', 'Quality Assurance', 'Timeline Management']
    },
    {
      icon: 'fas fa-handshake',
      title: 'Mortgage Assistance',
      description: 'Professional mortgage and financing assistance to secure the best deals.',
      features: ['Loan Processing', 'Rate Comparison', 'Documentation Support', 'Bank Liaison']
    }
  ]

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-50 via-blue-50 to-pink-50">
      {/* Hero Section */}
      <section 
        className="h-[40vh] bg-gradient-to-r from-purple-600/90 to-blue-600/90 bg-cover bg-center flex flex-col justify-center relative overflow-hidden"
        style={{ backgroundImage: "url('https://images.unsplash.com/photo-1486406146926-c627a92ad1ab?ixlib=rb-1.2.1&auto=format&fit=crop&w=1350&q=80')" }}
      >
        {/* Floating shapes */}
        <div className="absolute inset-0 overflow-hidden">
          <div className="absolute top-20 left-10 w-20 h-20 bg-white/10 rounded-full blur-xl animate-pulse"></div>
          <div className="absolute bottom-20 right-10 w-32 h-32 bg-pink-400/20 rounded-full blur-xl animate-pulse delay-1000"></div>
          <div className="absolute top-1/2 left-1/2 w-24 h-24 bg-blue-400/20 rounded-full blur-xl animate-pulse delay-500"></div>
        </div>
        
        <div className="max-w-4xl mx-auto text-center text-white relative z-10 px-5 mt-20">
          <h1 className="text-5xl md:text-6xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-white to-blue-100">
            Our Services
          </h1>
          <p className="text-xl md:text-2xl text-blue-100">
            Comprehensive real estate solutions tailored to your needs
          </p>
        </div>
      </section>

      {/* Main Services Section */}
      <section className="py-20 px-5">
        <div className="max-w-7xl mx-auto">
          {/* Intro */}
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold mb-6 bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent">
              What We Offer
            </h2>
            <p className="text-lg text-gray-700 max-w-3xl mx-auto">
              From buying your first home to building a real estate empire, we provide end-to-end services 
              that cover every aspect of the property journey.
            </p>
          </div>

          {/* Services Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 mb-20">
            {services.map((service, index) => (
              <div key={index} className="bg-white/70 backdrop-blur-sm rounded-2xl p-8 shadow-xl border border-white/20 hover:shadow-2xl transition-all hover:scale-105 group">
                <div className="w-16 h-16 bg-gradient-to-r from-purple-500 to-blue-500 rounded-full flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
                  <i className={`${service.icon} text-white text-xl`}></i>
                </div>
                <h3 className="text-xl font-bold mb-4 text-gray-800 group-hover:text-purple-600 transition-colors">
                  {service.title}
                </h3>
                <p className="text-gray-600 mb-6 leading-relaxed">
                  {service.description}
                </p>
                <ul className="space-y-2">
                  {service.features.map((feature, idx) => (
                    <li key={idx} className="flex items-center text-sm text-gray-600">
                      <i className="fas fa-check text-green-500 mr-2"></i>
                      {feature}
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>

          {/* Process Section */}
          <div className="bg-white/60 backdrop-blur-sm rounded-2xl p-8 md:p-12 shadow-xl border border-white/20 mb-20">
            <div className="text-center mb-12">
              <h2 className="text-4xl font-bold mb-6 bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent">
                Our Process
              </h2>
              <p className="text-lg text-gray-700 max-w-2xl mx-auto">
                We follow a systematic approach to ensure smooth and successful property transactions.
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
              <div className="text-center">
                <div className="w-16 h-16 bg-gradient-to-r from-purple-500 to-pink-500 rounded-full flex items-center justify-center mx-auto mb-4">
                  <span className="text-white font-bold text-xl">1</span>
                </div>
                <h3 className="font-bold text-lg mb-2">Consultation</h3>
                <p className="text-gray-600 text-sm">Initial discussion to understand your needs and goals</p>
              </div>
              
              <div className="text-center">
                <div className="w-16 h-16 bg-gradient-to-r from-pink-500 to-blue-500 rounded-full flex items-center justify-center mx-auto mb-4">
                  <span className="text-white font-bold text-xl">2</span>
                </div>
                <h3 className="font-bold text-lg mb-2">Analysis</h3>
                <p className="text-gray-600 text-sm">Comprehensive market analysis and property evaluation</p>
              </div>
              
              <div className="text-center">
                <div className="w-16 h-16 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full flex items-center justify-center mx-auto mb-4">
                  <span className="text-white font-bold text-xl">3</span>
                </div>
                <h3 className="font-bold text-lg mb-2">Execution</h3>
                <p className="text-gray-600 text-sm">Implementation of strategy with continuous support</p>
              </div>
              
              <div className="text-center">
                <div className="w-16 h-16 bg-gradient-to-r from-purple-500 to-blue-500 rounded-full flex items-center justify-center mx-auto mb-4">
                  <span className="text-white font-bold text-xl">4</span>
                </div>
                <h3 className="font-bold text-lg mb-2">Completion</h3>
                <p className="text-gray-600 text-sm">Successful completion with ongoing relationship</p>
              </div>
            </div>
          </div>

          {/* Call to Action */}
          <div className="text-center bg-gradient-to-r from-purple-600 to-blue-600 rounded-2xl p-12 text-white">
            <h2 className="text-4xl font-bold mb-6">Ready to Get Started?</h2>
            <p className="text-xl mb-8 opacity-90">
              Let's discuss how we can help you achieve your real estate goals.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <a
                href="https://wa.me/254757577018"
                className="bg-white text-purple-600 px-8 py-3 rounded-full font-semibold hover:bg-gray-100 transition-all hover:scale-105 flex items-center justify-center gap-2"
                target="_blank"
                rel="noopener noreferrer"
              >
                <i className="fab fa-whatsapp"></i>
                Chat on WhatsApp
              </a>
              <a
                href="tel:+254757577018"
                className="bg-transparent border-2 border-white text-white px-8 py-3 rounded-full font-semibold hover:bg-white hover:text-purple-600 transition-all hover:scale-105 flex items-center justify-center gap-2"
              >
                <i className="fas fa-phone"></i>
                Call Now
              </a>
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

export default Services
