import { useEffect } from 'react';

const Contact = () => {
  useEffect(() => {
    window.scrollTo(0, 0);
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-br from-green-50 via-yellow-50 to-blue-50">
      {/* Hero Section */}
      <section 
        className="h-[40vh] bg-gradient-to-r from-green-500/90 to-blue-600/90 bg-cover bg-center flex flex-col justify-center relative overflow-hidden"
        style={{ backgroundImage: "url('https://images.unsplash.com/photo-1528909514045-2fa4ac7a08ba?ixlib=rb-1.2.1&auto=format&fit=crop&w=1350&q=80')" }}
      >
        <div className="max-w-4xl mx-auto text-center text-white relative z-10 px-5 mt-20">
          <h1 className="text-5xl md:text-6xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-white to-green-100">
            Contact Us
          </h1>
          <p className="text-xl md:text-2xl text-green-100">
            We’re here to help with any questions or concerns
          </p>
        </div>
      </section>

      {/* Contact Form Section */}
      <section className="py-20 px-5">
        <div className="max-w-4xl mx-auto">
          <div className="bg-white/70 backdrop-blur-sm rounded-2xl p-8 shadow-xl border border-white/20">
            <h2 className="text-3xl font-bold mb-6 text-gray-800">
              Get In Touch
            </h2>
            <form className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-gray-700">Name</label>
                <input
                  type="text"
                  className="mt-1 block w-full p-3 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Your Name"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">Email</label>
                <input
                  type="email"
                  className="mt-1 block w-full p-3 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                  placeholder="you@example.com"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">Message</label>
                <textarea
                  rows="4"
                  className="mt-1 block w-full p-3 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Your message..."
                ></textarea>
              </div>
              <div>
                <button
                  type="submit"
                  className="w-full py-3 px-6 bg-blue-600 text-white font-semibold rounded-md shadow-md hover:bg-blue-700 transition-all"
                >
                  Send Message
                </button>
              </div>
            </form>
          </div>
        </div>
      </section>
    </div>
  );
};

export default Contact;

