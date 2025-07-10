import 'package:flutter/material.dart';

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({super.key});

  @override
  State<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> {
  // Sample user data
  final Map<String, dynamic> _userData = {
    'name': 'John Doe',
    'email': 'john.doe@example.com',
    'phone': '+1 (555) 123-4567',
    'profileImage': '',
    'memberSince': 'May 2023',
    'bookings': 12,
    'savedProperties': 8,
    'reviews': 5,
  };

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('My Profile'),
        actions: [
          IconButton(
            icon: const Icon(Icons.settings),
            onPressed: () {
              // Navigate to settings
            },
          ),
        ],
      ),
      body: SingleChildScrollView(
        child: Column(
          children: [
            // Profile header
            Container(
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: Theme.of(context).colorScheme.primary.withOpacity(0.1),
              ),
              child: Column(
                children: [
                  // Profile image
                  CircleAvatar(
                    radius: 50,
                    backgroundColor: Theme.of(context).colorScheme.primary,
                    child: const Icon(
                      Icons.person,
                      size: 50,
                      color: Colors.white,
                    ),
                  ),
                  const SizedBox(height: 16),
                  // User name
                  Text(
                    _userData['name'],
                    style: const TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 4),
                  // Member since
                  Text(
                    'Member since ${_userData['memberSince']}',
                    style: TextStyle(
                      color: Colors.grey.shade600,
                    ),
                  ),
                  const SizedBox(height: 16),
                  // Edit profile button
                  OutlinedButton.icon(
                    onPressed: () {
                      // Navigate to edit profile
                      Navigator.pushNamed(context, '/profile_edit');
                    },
                    icon: const Icon(Icons.edit),
                    label: const Text('Edit Profile'),
                  ),
                ],
              ),
            ),
            
            // Profile stats
            Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                  _buildStatItem(
                    context,
                    Icons.calendar_today,
                    '${_userData['bookings']}',
                    'Bookings',
                  ),
                  _buildStatItem(
                    context,
                    Icons.favorite,
                    '${_userData['savedProperties']}',
                    'Saved',
                  ),
                  _buildStatItem(
                    context,
                    Icons.star,
                    '${_userData['reviews']}',
                    'Reviews',
                  ),
                ],
              ),
            ),
            
            const Divider(),
            
            // Profile sections
            _buildProfileSection(
              context,
              'Personal Information',
              Icons.person,
              [
                _buildInfoItem('Email', _userData['email']),
                _buildInfoItem('Phone', _userData['phone']),
              ],
            ),
            
            _buildProfileSection(
              context,
              'Payment Methods',
              Icons.payment,
              [
                _buildPaymentMethodItem(
                  'Visa ending in 4242',
                  'Expires 12/25',
                  Icons.credit_card,
                ),
                _buildPaymentMethodItem(
                  'PayPal',
                  _userData['email'],
                  Icons.account_balance_wallet,
                ),
              ],
              actionText: 'Add New',
              onActionPressed: () {
                // Navigate to add payment method
              },
            ),
            
            _buildProfileSection(
              context,
              'Notifications',
              Icons.notifications,
              [
                _buildSwitchItem(
                  'Booking Updates',
                  'Get notified about your booking status',
                  true,
                  (value) {
                    // Update notification preference
                  },
                ),
                _buildSwitchItem(
                  'New Properties',
                  'Get notified about new properties matching your preferences',
                  false,
                  (value) {
                    // Update notification preference
                  },
                ),
                _buildSwitchItem(
                  'Promotions',
                  'Receive promotional offers and discounts',
                  true,
                  (value) {
                    // Update notification preference
                  },
                ),
              ],
            ),
            
            _buildProfileSection(
              context,
              'Help & Support',
              Icons.help,
              [
                _buildActionItem(
                  'Frequently Asked Questions',
                  Icons.question_answer,
                  () {
                    // Navigate to FAQs
                  },
                ),
                _buildActionItem(
                  'Contact Support',
                  Icons.support_agent,
                  () {
                    // Navigate to support
                  },
                ),
                _buildActionItem(
                  'Terms & Conditions',
                  Icons.description,
                  () {
                    // Navigate to terms
                  },
                ),
                _buildActionItem(
                  'Privacy Policy',
                  Icons.privacy_tip,
                  () {
                    // Navigate to privacy policy
                  },
                ),
              ],
            ),
            
            // Logout button
            Padding(
              padding: const EdgeInsets.all(24),
              child: SizedBox(
                width: double.infinity,
                child: OutlinedButton.icon(
                  onPressed: () {
                    // Show logout confirmation
                    showDialog(
                      context: context,
                      builder: (context) => AlertDialog(
                        title: const Text('Logout'),
                        content: const Text('Are you sure you want to logout?'),
                        actions: [
                          TextButton(
                            onPressed: () {
                              Navigator.pop(context);
                            },
                            child: const Text('Cancel'),
                          ),
                          TextButton(
                            onPressed: () {
                              Navigator.pop(context);
                              // Perform logout
                              Navigator.pushReplacementNamed(context, '/login');
                            },
                            style: TextButton.styleFrom(
                              foregroundColor: Colors.red,
                            ),
                            child: const Text('Logout'),
                          ),
                        ],
                      ),
                    );
                  },
                  icon: const Icon(Icons.logout, color: Colors.red),
                  label: const Text(
                    'Logout',
                    style: TextStyle(color: Colors.red),
                  ),
                  style: OutlinedButton.styleFrom(
                    side: const BorderSide(color: Colors.red),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatItem(
    BuildContext context,
    IconData icon,
    String value,
    String label,
  ) {
    return Column(
      children: [
        Container(
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: Theme.of(context).colorScheme.primary.withOpacity(0.1),
            shape: BoxShape.circle,
          ),
          child: Icon(
            icon,
            color: Theme.of(context).colorScheme.primary,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          value,
          style: const TextStyle(
            fontWeight: FontWeight.bold,
            fontSize: 16,
          ),
        ),
        Text(
          label,
          style: TextStyle(
            color: Colors.grey.shade600,
            fontSize: 12,
          ),
        ),
      ],
    );
  }

  Widget _buildProfileSection(
    BuildContext context,
    String title,
    IconData icon,
    List<Widget> items, {
    String? actionText,
    VoidCallback? onActionPressed,
  }) {
    return Container(
      margin: const EdgeInsets.symmetric(vertical: 8),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Row(
                  children: [
                    Icon(
                      icon,
                      color: Theme.of(context).colorScheme.primary,
                    ),
                    const SizedBox(width: 8),
                    Text(
                      title,
                      style: const TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                if (actionText != null && onActionPressed != null)
                  TextButton(
                    onPressed: onActionPressed,
                    child: Text(actionText),
                  ),
              ],
            ),
          ),
          ...items,
          const Divider(),
        ],
      ),
    );
  }

  Widget _buildInfoItem(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: TextStyle(
              color: Colors.grey.shade600,
            ),
          ),
          Text(
            value,
            style: const TextStyle(
              fontWeight: FontWeight.bold,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildPaymentMethodItem(
    String title,
    String subtitle,
    IconData icon,
  ) {
    return ListTile(
      leading: Icon(icon),
      title: Text(title),
      subtitle: Text(subtitle),
      trailing: IconButton(
        icon: const Icon(Icons.more_vert),
        onPressed: () {
          // Show payment method options
        },
      ),
    );
  }

  Widget _buildSwitchItem(
    String title,
    String subtitle,
    bool value,
    Function(bool) onChanged,
  ) {
    return ListTile(
      title: Text(title),
      subtitle: Text(subtitle),
      trailing: Switch(
        value: value,
        onChanged: onChanged,
      ),
    );
  }

  Widget _buildActionItem(
    String title,
    IconData icon,
    VoidCallback onTap,
  ) {
    return ListTile(
      leading: Icon(icon),
      title: Text(title),
      trailing: const Icon(Icons.chevron_right),
      onTap: onTap,
    );
  }
}
