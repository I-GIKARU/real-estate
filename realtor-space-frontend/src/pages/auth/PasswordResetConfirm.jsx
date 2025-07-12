import { useState, useEffect } from 'react'
import { useNavigate, useSearchParams, Link } from 'react-router-dom'
import { confirmPasswordReset } from '../../services/api'

function PasswordResetConfirm() {
    const [formData, setFormData] = useState({
        token: '',
        newPassword: '',
        confirmPassword: ''
    })
    const [errors, setErrors] = useState({})
    const [isLoading, setIsLoading] = useState(false)
    const [message, setMessage] = useState('')
    const [searchParams] = useSearchParams()
    const navigate = useNavigate()

    // Extract token from URL parameters
    useEffect(() => {
        const urlToken = searchParams.get('token')
        if (urlToken) {
            setFormData(prev => ({ ...prev, token: urlToken }))
        }
    }, [searchParams])

    const validateForm = () => {
        const newErrors = {}

        if (!formData.token.trim()) {
            newErrors.token = 'Reset token is required'
        }

        if (!formData.newPassword) {
            newErrors.newPassword = 'New password is required'
        } else if (formData.newPassword.length < 8) {
            newErrors.newPassword = 'Password must be at least 8 characters'
        } else if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(formData.newPassword)) {
            newErrors.newPassword = 'Password must contain at least one uppercase letter, one lowercase letter, and one number'
        }

        if (!formData.confirmPassword) {
            newErrors.confirmPassword = 'Please confirm your password'
        } else if (formData.newPassword !== formData.confirmPassword) {
            newErrors.confirmPassword = 'Passwords do not match'
        }

        setErrors(newErrors)
        return Object.keys(newErrors).length === 0
    }

    const handleChange = (e) => {
        const { name, value } = e.target
        setFormData(prev => ({
            ...prev,
            [name]: value
        }))
        
        // Clear field error when user starts typing
        if (errors[name]) {
            setErrors(prev => ({
                ...prev,
                [name]: ''
            }))
        }
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        
        if (!validateForm()) {
            return
        }

        setIsLoading(true)
        setMessage('')

        try {
            await confirmPasswordReset(formData.token, formData.newPassword)
            
            setMessage('Password reset successful! You can now login with your new password.')
            
            // Redirect to login after 3 seconds
            setTimeout(() => {
                navigate('/login')
            }, 3000)
            
        } catch (error) {
            setErrors({
                submit: error.response?.data?.message || 'Failed to reset password. Please try again.'
            })
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-900 via-purple-900 to-indigo-900 py-12 px-4 sm:px-6 lg:px-8">
            <div className="max-w-md w-full space-y-8">
                {/* Glassmorphic Container */}
                <div className="bg-white/10 backdrop-blur-xl rounded-3xl p-8 shadow-2xl border border-white/20">
                    <div>
                        <div className="mx-auto h-16 w-16 bg-gradient-to-r from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center mb-6">
                            <svg className="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                            </svg>
                        </div>
                        <h2 className="text-center text-3xl font-bold text-white mb-2">
                            Reset Your Password
                        </h2>
                        <p className="text-center text-gray-300 text-sm mb-8">
                            Enter your new password below
                        </p>
                    </div>

                    {message && (
                        <div className="mb-6 p-4 bg-green-500/20 border border-green-500/30 rounded-xl">
                            <p className="text-green-300 text-sm text-center">{message}</p>
                        </div>
                    )}

                    {errors.submit && (
                        <div className="mb-6 p-4 bg-red-500/20 border border-red-500/30 rounded-xl">
                            <p className="text-red-300 text-sm text-center">{errors.submit}</p>
                        </div>
                    )}

                    <form className="space-y-6" onSubmit={handleSubmit}>
                        <div>
                            <label htmlFor="token" className="block text-sm font-medium text-gray-300 mb-2">
                                Reset Token
                            </label>
                            <input
                                id="token"
                                name="token"
                                type="text"
                                value={formData.token}
                                onChange={handleChange}
                                className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-xl text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
                                placeholder="Enter reset token"
                            />
                            {errors.token && <p className="mt-2 text-sm text-red-400">{errors.token}</p>}
                        </div>

                        <div>
                            <label htmlFor="newPassword" className="block text-sm font-medium text-gray-300 mb-2">
                                New Password
                            </label>
                            <input
                                id="newPassword"
                                name="newPassword"
                                type="password"
                                value={formData.newPassword}
                                onChange={handleChange}
                                className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-xl text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
                                placeholder="Enter new password"
                            />
                            {errors.newPassword && <p className="mt-2 text-sm text-red-400">{errors.newPassword}</p>}
                        </div>

                        <div>
                            <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-300 mb-2">
                                Confirm New Password
                            </label>
                            <input
                                id="confirmPassword"
                                name="confirmPassword"
                                type="password"
                                value={formData.confirmPassword}
                                onChange={handleChange}
                                className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-xl text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
                                placeholder="Confirm new password"
                            />
                            {errors.confirmPassword && <p className="mt-2 text-sm text-red-400">{errors.confirmPassword}</p>}
                        </div>

                        <div>
                            <button
                                type="submit"
                                disabled={isLoading}
                                className="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-xl text-white bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 transform hover:scale-105"
                            >
                                {isLoading ? (
                                    <div className="flex items-center">
                                        <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                        </svg>
                                        Resetting Password...
                                    </div>
                                ) : (
                                    'Reset Password'
                                )}
                            </button>
                        </div>

                        <div className="text-center">
                            <Link
                                to="/login"
                                className="text-blue-400 hover:text-blue-300 text-sm transition-colors duration-200"
                            >
                                Back to Login
                            </Link>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    )
}

export default PasswordResetConfirm
