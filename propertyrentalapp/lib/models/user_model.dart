class User {
  final String id;
  final String email;
  final String firstName;
  final String lastName;
  final String phoneNumber;
  final String userType;
  final bool isEmailVerified;
  final bool isApproved; // For agent approval by admin
  final DateTime? emailVerifiedAt;
  final DateTime? approvedAt;
  final String? profilePicture;
  final DateTime createdAt;
  final DateTime updatedAt;

  User({
    required this.id,
    required this.email,
    required this.firstName,
    required this.lastName,
    required this.phoneNumber,
    required this.userType,
    required this.isEmailVerified,
    this.isApproved = false,
    this.emailVerifiedAt,
    this.approvedAt,
    this.profilePicture,
    required this.createdAt,
    required this.updatedAt,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] as String,
      email: json['email'] as String,
      firstName: json['first_name'] as String,
      lastName: json['last_name'] as String,
      phoneNumber: json['phone_number'] as String,
      userType: json['user_type'] as String,
      // Backend uses 'is_verified' not 'is_email_verified'
      isEmailVerified: json['is_verified'] as bool? ?? false,
      isApproved: json['is_approved'] as bool? ?? false,
      emailVerifiedAt: json['email_verified_at'] != null 
          ? DateTime.parse(json['email_verified_at']) 
          : null,
      approvedAt: json['approved_at'] != null 
          ? DateTime.parse(json['approved_at']) 
          : null,
      // Backend uses 'profile_image_url' not 'profile_picture'
      profilePicture: json['profile_image_url'] as String?,
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email': email,
      'first_name': firstName,
      'last_name': lastName,
      'phone_number': phoneNumber,
      'user_type': userType,
      'is_email_verified': isEmailVerified,
      'is_approved': isApproved,
      'email_verified_at': emailVerifiedAt?.toIso8601String(),
      'approved_at': approvedAt?.toIso8601String(),
      'profile_picture': profilePicture,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  String get fullName => '$firstName $lastName';
}

class AuthResponse {
  final String token;
  final String refreshToken;
  final User user;

  AuthResponse({
    required this.token,
    this.refreshToken = '', // Made optional since backend doesn't provide it
    required this.user,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) {
    return AuthResponse(
      token: json['token'] as String,
      refreshToken: json['refresh_token'] as String? ?? '', // Handle missing refresh token
      user: User.fromJson(json['user']),
    );
  }
}
