import axios from 'axios';

// Configuration
const BASE_URL = import.meta.env.VITE_API_BASE_URL;

// Agent registration data
const agentData = {
  email: 'ict622j1a@gmail.com',
  password: 'SecurePassword123!', // You should change this to a secure password
  name: 'Real Estate Agent',
  phone: '+254700000000', // Update with actual phone number
  role: 'agent'
};

// Function to register the agent
export const registerAgent = async () => {
  try {
    console.log('Registering agent with email:', agentData.email);
    
    const response = await axios.post(`${BASE_URL}/register`, agentData);
    
    console.log('Agent registered successfully:', response.data);
    
    // If registration returns a token, save it
    if (response.data.token) {
      localStorage.setItem('authToken', response.data.token);
      console.log('Auth token saved to localStorage');
    }
    
    return response.data;
  } catch (error) {
    console.error('Error registering agent:', error.response?.data || error.message);
    throw error;
  }
};

// Function to login the agent (if registration fails because user already exists)
export const loginAgent = async () => {
  try {
    console.log('Logging in agent with email:', agentData.email);
    
    const response = await axios.post(`${BASE_URL}/login`, {
      email: agentData.email,
      password: agentData.password
    });
    
    console.log('Agent logged in successfully:', response.data);
    
    // Save the token
    if (response.data.token) {
      localStorage.setItem('authToken', response.data.token);
      console.log('Auth token saved to localStorage');
    }
    
    return response.data;
  } catch (error) {
    console.error('Error logging in agent:', error.response?.data || error.message);
    throw error;
  }
};

// Function to get or create agent (register if new, login if exists)
export const getOrCreateAgent = async () => {
  try {
    // First try to register
    return await registerAgent();
  } catch (error) {
    // If registration fails (user might already exist), try to login
    if (error.response?.status === 409 || error.response?.status === 400) {
      console.log('User might already exist, trying to login...');
      return await loginAgent();
    }
    throw error;
  }
};

// If running this script directly
if (require.main === module) {
  getOrCreateAgent()
    .then(result => {
      console.log('Agent setup completed:', result);
      process.exit(0);
    })
    .catch(error => {
      console.error('Agent setup failed:', error);
      process.exit(1);
    });
}
