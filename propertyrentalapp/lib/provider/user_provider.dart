import 'package:flutter/foundation.dart';
import '../models/user_model.dart';
import '../services/auth_service.dart';

class UserProvider with ChangeNotifier {
  User? _user;
  String? _token;
  bool _isLoading = false;
  String? _errorMessage;

  User? get user => _user;
  String? get token => _token;
  bool get isLoading => _isLoading;
  String? get errorMessage => _errorMessage;
  bool get isAuthenticated => _user != null;
  bool get isTenant => _user?.userType == 'tenant';
  bool get isAgent => _user?.userType == 'agent';
  bool get isAdmin => _user?.userType == 'admin';
  
  // Agent approval status
  bool get isAgentApproved => isAgent && (_user?.isApproved ?? false);
  bool get canManageProperties => isAgent && isAgentApproved;

  final AuthService _authService = AuthService();

  // Set loading state
  void _setLoading(bool loading) {
    _isLoading = loading;
    notifyListeners();
  }

  // Set error message
  void _setError(String? error) {
    _errorMessage = error;
    notifyListeners();
  }

  // Clear error
  void clearError() {
    _errorMessage = null;
    notifyListeners();
  }

  // Set user
  void setUser(User user) {
    _user = user;
    notifyListeners();
  }

  // Set token
  void setToken(String token) {
    _token = token;
    notifyListeners();
  }

  // Initialize user from stored data
  Future<void> initializeUser() async {
    _setLoading(true);
    try {
      final storedUser = await _authService.getStoredUser();
      if (storedUser != null) {
        _user = storedUser;
        
        // Try to get fresh user data
        final freshUser = await _authService.getProfile();
        if (freshUser != null) {
          _user = freshUser;
        }
      }
    } catch (e) {
      _setError('Failed to load user data');
    } finally {
      _setLoading(false);
    }
  }

  // Login user
  Future<bool> login(String email, String password) async {
    _setLoading(true);
    _setError(null);
    
    try {
      final authResponse = await _authService.login(email, password);
      if (authResponse != null) {
        _user = authResponse.user;
        _setLoading(false);
        return true;
      } else {
        _setError('Invalid email or password');
        _setLoading(false);
        return false;
      }
    } catch (e) {
      _setError(e.toString());
      _setLoading(false);
      return false;
    }
  }

  // Register new user
  Future<bool> register({
    required String email,
    required String password,
    required String firstName,
    required String lastName,
    required String phoneNumber,
    required String userType,
  }) async {
    _setLoading(true);
    _setError(null);
    
    try {
      final authResponse = await _authService.register(
        email: email,
        password: password,
        firstName: firstName,
        lastName: lastName,
        phoneNumber: phoneNumber,
        userType: userType,
      );
      
      if (authResponse != null) {
        _user = authResponse.user;
        _setLoading(false);
        return true;
      } else {
        _setError('Registration failed');
        _setLoading(false);
        return false;
      }
    } catch (e) {
      _setError(e.toString());
      _setLoading(false);
      return false;
    }
  }

  // Update user profile
  Future<bool> updateProfile({
    String? firstName,
    String? lastName,
    String? phoneNumber,
  }) async {
    _setLoading(true);
    _setError(null);
    
    try {
      final updatedUser = await _authService.updateProfile(
        firstName: firstName,
        lastName: lastName,
        phoneNumber: phoneNumber,
      );
      
      if (updatedUser != null) {
        _user = updatedUser;
        _setLoading(false);
        return true;
      } else {
        _setError('Failed to update profile');
        _setLoading(false);
        return false;
      }
    } catch (e) {
      _setError(e.toString());
      _setLoading(false);
      return false;
    }
  }

  // Send email verification
  Future<bool> sendVerificationEmail() async {
    try {
      return await _authService.sendVerificationEmail();
    } catch (e) {
      _setError(e.toString());
      return false;
    }
  }

  // Check verification status
  Future<bool> checkVerificationStatus() async {
    try {
      final isVerified = await _authService.getVerificationStatus();
      if (isVerified && _user != null) {
        // Update local user data if verified
        final freshUser = await _authService.getProfile();
        if (freshUser != null) {
          _user = freshUser;
          notifyListeners();
        }
      }
      return isVerified;
    } catch (e) {
      _setError(e.toString());
      return false;
    }
  }

  // Logout user
  Future<void> logout() async {
    await _authService.signOut();
    _user = null;
    _errorMessage = null;
    notifyListeners();
  }

  // Refresh user data
  Future<void> refreshUser() async {
    if (_user != null) {
      try {
        final freshUser = await _authService.getProfile();
        if (freshUser != null) {
          _user = freshUser;
          notifyListeners();
        }
      } catch (e) {
        _setError('Failed to refresh user data');
      }
    }
  }
}
