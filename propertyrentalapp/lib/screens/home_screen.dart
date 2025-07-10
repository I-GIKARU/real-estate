import 'package:flutter/material.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import '../widgets/creative_bottom_nav.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _selectedIndex = 0;
  BottomNavStyle _currentStyle = BottomNavStyle.curved;
  
  final List<Widget> _screens = [
    const PropertyListScreen(),
    const SavedPropertiesScreen(),
    const BookingsScreen(),
    const ProfileScreen(),
  ];
  
  final List<BottomNavItem> _navItems = [
    const BottomNavItem(
      icon: Icons.explore,
      activeIcon: Icons.explore,
      label: 'Explore',
    ),
    const BottomNavItem(
      icon: Icons.favorite_border,
      activeIcon: Icons.favorite,
      label: 'Saved',
    ),
    const BottomNavItem(
      icon: Icons.calendar_today_outlined,
      activeIcon: Icons.calendar_today,
      label: 'Bookings',
    ),
    const BottomNavItem(
      icon: Icons.person_outline,
      activeIcon: Icons.person,
      label: 'Profile',
    ),
  ];
  
  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  void onItemTapped(int index) {
    _onItemTapped(index);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _screens[_selectedIndex],
      bottomNavigationBar: CreativeBottomNav(
        currentIndex: _selectedIndex,
        onTap: _onItemTapped,
        items: _navItems,
        style: _currentStyle,
      ),
      // Alternative: Use extendBody for floating effect
      extendBody: true,
      // Add floating action button to change styles
      floatingActionButton: FloatingActionButton(
        mini: true,
        onPressed: _showStylePicker,
        child: const Icon(Icons.style),
      ),
      floatingActionButtonLocation: FloatingActionButtonLocation.miniEndTop,
    );
  }
  
  void _showStylePicker() {
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        margin: const EdgeInsets.all(16),
        padding: const EdgeInsets.all(20),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(20),
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text(
              'Choose Navigation Style',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 20),
            Wrap(
              spacing: 12,
              runSpacing: 12,
              children: BottomNavStyle.values.map((style) {
                return GestureDetector(
                  onTap: () {
                    setState(() {
                      _currentStyle = style;
                    });
                    Navigator.pop(context);
                  },
                  child: Container(
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                    decoration: BoxDecoration(
                      color: _currentStyle == style 
                          ? Theme.of(context).colorScheme.primary
                          : Colors.grey.shade200,
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: Text(
                      style.name.toUpperCase(),
                      style: TextStyle(
                        color: _currentStyle == style ? Colors.white : Colors.black,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                );
              }).toList(),
            ),
          ],
        ),
      ),
    );
  }

}

class PropertyListScreen extends StatelessWidget {
  const PropertyListScreen({super.key});

  @override
  Widget build(BuildContext context) {
    // Sample property data
    final properties = [
      {
        'name': 'Modern Apartment in Downtown',
        'address': '123 Main St, Downtown',
        'price': '\$1,200/month',
        'bedrooms': 2,
        'bathrooms': 1,
        'area': '850 sq ft',
        'rating': 4.5,
        'image': '',
        'isFeatured': true,
      },
      {
        'name': 'Luxury Villa with Pool',
        'address': '456 Park Ave, Westside',
        'price': '\$3,500/month',
        'bedrooms': 4,
        'bathrooms': 3,
        'area': '2,200 sq ft',
        'rating': 4.8,
        'image': '',
        'isFeatured': true,
      },
      {
        'name': 'Cozy Studio near University',
        'address': '789 College Blvd, Eastside',
        'price': '\$800/month',
        'bedrooms': 1,
        'bathrooms': 1,
        'area': '450 sq ft',
        'rating': 4.2,
        'image': '',
        'isFeatured': false,
      },
      {
        'name': 'Family Home with Garden',
        'address': '101 Suburban Dr, Northside',
        'price': '\$2,100/month',
        'bedrooms': 3,
        'bathrooms': 2,
        'area': '1,500 sq ft',
        'rating': 4.6,
        'image': '',
        'isFeatured': false,
      },
      {
        'name': 'Penthouse with City View',
        'address': '202 Skyline Ave, Downtown',
        'price': '\$4,000/month',
        'bedrooms': 3,
        'bathrooms': 2,
        'area': '1,800 sq ft',
        'rating': 4.9,
        'image': '',
        'isFeatured': true,
      },
    ];

    return Scaffold(
      appBar: AppBar(
        title: const Text('HomeRental'),
        actions: [
          IconButton(
            icon: const Icon(Icons.notifications),
            onPressed: () {
              // Navigate to notifications
            },
          ),
        ],
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Search bar
              TextField(
                decoration: InputDecoration(
                  hintText: 'Search for a location',
                  prefixIcon: const Icon(Icons.search),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  contentPadding: const EdgeInsets.symmetric(vertical: 12),
                ),
              ),
              const SizedBox(height: 24),
              
              // Filter options
              SingleChildScrollView(
                scrollDirection: Axis.horizontal,
                child: Row(
                  children: [
                    _buildFilterChip(label: 'All', isSelected: true),
                    _buildFilterChip(label: 'Apartment'),
                    _buildFilterChip(label: 'House'),
                    _buildFilterChip(label: 'Villa'),
                    _buildFilterChip(label: 'Studio'),
                    _buildFilterChip(label: 'Penthouse'),
                  ],
                ),
              ),
              const SizedBox(height: 24),
              
              // Featured properties
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  const Text(
                    'Featured Properties',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  TextButton(
                    onPressed: () {
                      // View all featured properties
                    },
                    child: const Text('View All'),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              
              // Featured properties list
              SizedBox(
                height: 280,
                child: ListView.builder(
                  scrollDirection: Axis.horizontal,
                  itemCount: properties.where((p) => p['isFeatured'] == true).length,
                  itemBuilder: (context, index) {
                    final featuredProperties = properties.where((p) => p['isFeatured'] == true).toList();
                    return _buildFeaturedPropertyCard(context, featuredProperties[index]);
                  },
                ),
              ),
              const SizedBox(height: 24),
              
              // Nearby properties
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  const Text(
                    'Nearby Properties',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  TextButton(
                    onPressed: () {
                      // View all nearby properties
                    },
                    child: const Text('View All'),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              
              // Nearby properties list
              ListView.builder(
                physics: const NeverScrollableScrollPhysics(),
                shrinkWrap: true,
                itemCount: properties.length,
                itemBuilder: (context, index) {
                  return _buildPropertyListItem(context, properties[index]);
                },
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildFilterChip({required String label, bool isSelected = false}) {
    return Container(
      margin: const EdgeInsets.only(right: 8),
      child: FilterChip(
        label: Text(label),
        selected: isSelected,
        onSelected: (selected) {
          // Handle filter selection
        },
      ),
    );
  }

  Widget _buildFeaturedPropertyCard(BuildContext context, Map<String, dynamic> property) {
    return GestureDetector(
      onTap: () {
        // Navigate to property details
        Navigator.pushNamed(context, '/property_details', arguments: property);
      },
      child: Container(
        width: 220,
        margin: const EdgeInsets.only(right: 16),
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: Colors.grey.shade200),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Property image
            Container(
              height: 140,
              decoration: BoxDecoration(
                color: Colors.grey.shade300,
                borderRadius: const BorderRadius.only(
                  topLeft: Radius.circular(12),
                  topRight: Radius.circular(12),
                ),
              ),
              child: const Center(
                child: Icon(
                  Icons.home,
                  size: 50,
                  color: Colors.white,
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(12.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Property price
                  Text(
                    property['price'],
                    style: TextStyle(
                      color: Theme.of(context).colorScheme.primary,
                      fontWeight: FontWeight.bold,
                      fontSize: 16,
                    ),
                  ),
                  const SizedBox(height: 4),
                  // Property name
                  Text(
                    property['name'],
                    style: const TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: 14,
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 4),
                  // Property address
                  Text(
                    property['address'],
                    style: TextStyle(
                      color: Colors.grey.shade600,
                      fontSize: 12,
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 8),
                  // Property features
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      _buildFeatureItem(Icons.bed, '${property['bedrooms']}'),
                      _buildFeatureItem(Icons.bathtub, '${property['bathrooms']}'),
                      _buildFeatureItem(Icons.square_foot, property['area']),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildPropertyListItem(BuildContext context, Map<String, dynamic> property) {
    return GestureDetector(
      onTap: () {
        // Navigate to property details
        Navigator.pushNamed(context, '/property_details', arguments: property);
      },
      child: Container(
        margin: const EdgeInsets.only(bottom: 16),
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: Colors.grey.shade200),
        ),
        child: Row(
          children: [
            // Property image
            Container(
              width: 120,
              height: 120,
              decoration: BoxDecoration(
                color: Colors.grey.shade300,
                borderRadius: const BorderRadius.only(
                  topLeft: Radius.circular(12),
                  bottomLeft: Radius.circular(12),
                ),
              ),
              child: const Center(
                child: Icon(
                  Icons.home,
                  size: 40,
                  color: Colors.white,
                ),
              ),
            ),
            // Property details
            Expanded(
              child: Padding(
                padding: const EdgeInsets.all(12.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Property price
                    Text(
                      property['price'],
                      style: TextStyle(
                        color: Theme.of(context).colorScheme.primary,
                        fontWeight: FontWeight.bold,
                        fontSize: 16,
                      ),
                    ),
                    const SizedBox(height: 4),
                    // Property name
                    Text(
                      property['name'],
                      style: const TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 14,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 4),
                    // Property address
                    Text(
                      property['address'],
                      style: TextStyle(
                        color: Colors.grey.shade600,
                        fontSize: 12,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 8),
                    // Property features
                    Row(
                      children: [
                        _buildFeatureItem(Icons.bed, '${property['bedrooms']}'),
                        const SizedBox(width: 16),
                        _buildFeatureItem(Icons.bathtub, '${property['bathrooms']}'),
                        const SizedBox(width: 16),
                        _buildFeatureItem(Icons.star, '${property['rating']}'),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFeatureItem(IconData icon, String text) {
    return Row(
      children: [
        Icon(
          icon,
          size: 16,
          color: Colors.grey.shade600,
        ),
        const SizedBox(width: 4),
        Text(
          text,
          style: TextStyle(
            color: Colors.grey.shade600,
            fontSize: 12,
          ),
        ),
      ],
    );
  }
}

class SavedPropertiesScreen extends StatefulWidget {
  const SavedPropertiesScreen({super.key});

  @override
  State<SavedPropertiesScreen> createState() => _SavedPropertiesScreenState();
}

class _SavedPropertiesScreenState extends State<SavedPropertiesScreen> {
  // Sample saved properties data
  final List<Map<String, dynamic>> savedProperties = [
    {
      'id': '1',
      'name': 'Modern Apartment in Downtown',
      'address': '123 Main St, Downtown',
      'price': '\$1,200/month',
      'bedrooms': 2,
      'bathrooms': 1,
      'area': '850 sq ft',
      'rating': 4.5,
      'image': '',
      'savedDate': '2024-01-15',
    },
    {
      'id': '2',
      'name': 'Luxury Villa with Pool',
      'address': '456 Park Ave, Westside',
      'price': '\$3,500/month',
      'bedrooms': 4,
      'bathrooms': 3,
      'area': '2,200 sq ft',
      'rating': 4.8,
      'image': '',
      'savedDate': '2024-01-12',
    },
    {
      'id': '3',
      'name': 'Cozy Studio near University',
      'address': '789 College Blvd, Eastside',
      'price': '\$800/month',
      'bedrooms': 1,
      'bathrooms': 1,
      'area': '450 sq ft',
      'rating': 4.2,
      'image': '',
      'savedDate': '2024-01-10',
    },
  ];

  void _removeFromSaved(String propertyId) {
    setState(() {
      savedProperties.removeWhere((property) => property['id'] == propertyId);
    });
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Property removed from saved list'),
        duration: Duration(seconds: 2),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Saved Properties'),
        automaticallyImplyLeading: false,
        actions: [
          IconButton(
            icon: const Icon(Icons.clear_all),
            onPressed: savedProperties.isEmpty ? null : () {
              showDialog(
                context: context,
                builder: (BuildContext context) {
                  return AlertDialog(
                    title: const Text('Clear All'),
                    content: const Text('Are you sure you want to remove all saved properties?'),
                    actions: [
                      TextButton(
                        onPressed: () => Navigator.of(context).pop(),
                        child: const Text('Cancel'),
                      ),
                      TextButton(
                        onPressed: () {
                          setState(() {
                            savedProperties.clear();
                          });
                          Navigator.of(context).pop();
                        },
                        child: const Text('Clear All'),
                      ),
                    ],
                  );
                },
              );
            },
          ),
        ],
      ),
      body: savedProperties.isEmpty
          ? _buildEmptyState()
          : ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: savedProperties.length,
              itemBuilder: (context, index) {
                final property = savedProperties[index];
                return _buildSavedPropertyCard(property);
              },
            ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.favorite_border,
            size: 100,
            color: Colors.grey.shade300,
          ),
          const SizedBox(height: 20),
          Text(
            'No Saved Properties',
            style: TextStyle(
              fontSize: 24,
              fontWeight: FontWeight.bold,
              color: Colors.grey.shade600,
            ),
          ),
          const SizedBox(height: 10),
          Text(
            'Start exploring and save your favorite properties',
            style: TextStyle(
              fontSize: 16,
              color: Colors.grey.shade500,
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 30),
          ElevatedButton(
            onPressed: () {
              // Switch to explore tab
              if (context.findAncestorStateOfType<_HomeScreenState>() != null) {
                context.findAncestorStateOfType<_HomeScreenState>()!.onItemTapped(0);
              }
            },
            child: const Text('Explore Properties'),
          ),
        ],
      ),
    );
  }

  Widget _buildSavedPropertyCard(Map<String, dynamic> property) {
    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
      ),
      elevation: 2,
      child: InkWell(
        onTap: () {
          Navigator.pushNamed(context, '/property_details', arguments: property);
        },
        borderRadius: BorderRadius.circular(12),
        child: Column(
          children: [
            Row(
              children: [
                // Property image
                Container(
                  width: 120,
                  height: 120,
                  decoration: BoxDecoration(
                    color: Colors.grey.shade300,
                    borderRadius: const BorderRadius.only(
                      topLeft: Radius.circular(12),
                      bottomLeft: Radius.circular(12),
                    ),
                  ),
                  child: const Center(
                    child: Icon(
                      Icons.home,
                      size: 40,
                      color: Colors.white,
                    ),
                  ),
                ),
                // Property details
                Expanded(
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              property['price'],
                              style: TextStyle(
                                color: Theme.of(context).colorScheme.primary,
                                fontWeight: FontWeight.bold,
                                fontSize: 16,
                              ),
                            ),
                            IconButton(
                              icon: const Icon(Icons.favorite, color: Colors.red),
                              onPressed: () => _removeFromSaved(property['id']),
                            ),
                          ],
                        ),
                        Text(
                          property['name'],
                          style: const TextStyle(
                            fontWeight: FontWeight.bold,
                            fontSize: 16,
                          ),
                          maxLines: 2,
                          overflow: TextOverflow.ellipsis,
                        ),
                        const SizedBox(height: 4),
                        Text(
                          property['address'],
                          style: TextStyle(
                            color: Colors.grey.shade600,
                            fontSize: 14,
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                        const SizedBox(height: 12),
                        Row(
                          children: [
                            _buildPropertyFeature(Icons.bed, '${property['bedrooms']}'),
                            const SizedBox(width: 16),
                            _buildPropertyFeature(Icons.bathtub, '${property['bathrooms']}'),
                            const SizedBox(width: 16),
                            _buildPropertyFeature(Icons.star, '${property['rating']}'),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            ),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              decoration: BoxDecoration(
                color: Colors.grey.shade50,
                borderRadius: const BorderRadius.only(
                  bottomLeft: Radius.circular(12),
                  bottomRight: Radius.circular(12),
                ),
              ),
              child: Row(
                children: [
                  Icon(
                    Icons.access_time,
                    size: 16,
                    color: Colors.grey.shade600,
                  ),
                  const SizedBox(width: 4),
                  Text(
                    'Saved on ${property['savedDate']}',
                    style: TextStyle(
                      color: Colors.grey.shade600,
                      fontSize: 12,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildPropertyFeature(IconData icon, String text) {
    return Row(
      children: [
        Icon(
          icon,
          size: 16,
          color: Colors.grey.shade600,
        ),
        const SizedBox(width: 4),
        Text(
          text,
          style: TextStyle(
            color: Colors.grey.shade600,
            fontSize: 12,
          ),
        ),
      ],
    );
  }
}

class BookingsScreen extends StatefulWidget {
  const BookingsScreen({super.key});

  @override
  State<BookingsScreen> createState() => _BookingsScreenState();
}

class _BookingsScreenState extends State<BookingsScreen> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  
  // Sample bookings data
  final List<Map<String, dynamic>> activeBookings = [
    {
      'id': '1',
      'propertyName': 'Modern Apartment in Downtown',
      'address': '123 Main St, Downtown',
      'price': '\$1,200/month',
      'checkIn': '2024-02-01',
      'checkOut': '2024-02-28',
      'status': 'confirmed',
      'landlord': 'John Doe',
      'landlordPhone': '+254 700 123456',
      'image': '',
    },
    {
      'id': '2',
      'propertyName': 'Cozy Studio near University',
      'address': '789 College Blvd, Eastside',
      'price': '\$800/month',
      'checkIn': '2024-02-15',
      'checkOut': '2024-05-15',
      'status': 'pending',
      'landlord': 'Jane Smith',
      'landlordPhone': '+254 700 987654',
      'image': '',
    },
  ];
  
  final List<Map<String, dynamic>> pastBookings = [
    {
      'id': '3',
      'propertyName': 'Family Home with Garden',
      'address': '101 Suburban Dr, Northside',
      'price': '\$2,100/month',
      'checkIn': '2023-09-01',
      'checkOut': '2023-12-31',
      'status': 'completed',
      'landlord': 'Mike Johnson',
      'landlordPhone': '+254 700 456789',
      'image': '',
      'rating': 4.5,
    },
  ];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('My Bookings'),
        automaticallyImplyLeading: false,
        bottom: TabBar(
          controller: _tabController,
          labelColor: Theme.of(context).colorScheme.primary,
          unselectedLabelColor: Colors.grey,
          indicatorColor: Theme.of(context).colorScheme.primary,
          tabs: [
            Tab(
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Icon(Icons.schedule),
                  const SizedBox(width: 8),
                  Text('Active (${activeBookings.length})'),
                ],
              ),
            ),
            Tab(
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Icon(Icons.history),
                  const SizedBox(width: 8),
                  Text('Past (${pastBookings.length})'),
                ],
              ),
            ),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildActiveBookings(),
          _buildPastBookings(),
        ],
      ),
    );
  }

  Widget _buildActiveBookings() {
    if (activeBookings.isEmpty) {
      return _buildEmptyState(
        icon: Icons.calendar_today,
        title: 'No Active Bookings',
        subtitle: 'You don\'t have any active bookings yet',
        actionText: 'Explore Properties',
        onActionPressed: () {
          if (context.findAncestorStateOfType<_HomeScreenState>() != null) {
            context.findAncestorStateOfType<_HomeScreenState>()!.onItemTapped(0);
          }
        },
      );
    }
    
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: activeBookings.length,
      itemBuilder: (context, index) {
        final booking = activeBookings[index];
        return _buildBookingCard(booking, isActive: true);
      },
    );
  }

  Widget _buildPastBookings() {
    if (pastBookings.isEmpty) {
      return _buildEmptyState(
        icon: Icons.history,
        title: 'No Past Bookings',
        subtitle: 'Your booking history will appear here',
        actionText: 'Explore Properties',
        onActionPressed: () {
          if (context.findAncestorStateOfType<_HomeScreenState>() != null) {
            context.findAncestorStateOfType<_HomeScreenState>()!.onItemTapped(0);
          }
        },
      );
    }
    
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: pastBookings.length,
      itemBuilder: (context, index) {
        final booking = pastBookings[index];
        return _buildBookingCard(booking, isActive: false);
      },
    );
  }

  Widget _buildEmptyState({
    required IconData icon,
    required String title,
    required String subtitle,
    required String actionText,
    required VoidCallback onActionPressed,
  }) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            icon,
            size: 100,
            color: Colors.grey.shade300,
          ),
          const SizedBox(height: 20),
          Text(
            title,
            style: TextStyle(
              fontSize: 24,
              fontWeight: FontWeight.bold,
              color: Colors.grey.shade600,
            ),
          ),
          const SizedBox(height: 10),
          Text(
            subtitle,
            style: TextStyle(
              fontSize: 16,
              color: Colors.grey.shade500,
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 30),
          ElevatedButton(
            onPressed: onActionPressed,
            child: Text(actionText),
          ),
        ],
      ),
    );
  }

  Widget _buildBookingCard(Map<String, dynamic> booking, {required bool isActive}) {
    Color statusColor;
    String statusText;
    
    switch (booking['status']) {
      case 'confirmed':
        statusColor = Colors.green;
        statusText = 'Confirmed';
        break;
      case 'pending':
        statusColor = Colors.orange;
        statusText = 'Pending';
        break;
      case 'completed':
        statusColor = Colors.blue;
        statusText = 'Completed';
        break;
      default:
        statusColor = Colors.grey;
        statusText = 'Unknown';
    }
    
    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
      ),
      elevation: 2,
      child: InkWell(
        onTap: () {
          _showBookingDetails(booking);
        },
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Container(
                    width: 80,
                    height: 80,
                    decoration: BoxDecoration(
                      color: Colors.grey.shade300,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: const Center(
                      child: Icon(
                        Icons.home,
                        size: 30,
                        color: Colors.white,
                      ),
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          booking['propertyName'],
                          style: const TextStyle(
                            fontWeight: FontWeight.bold,
                            fontSize: 16,
                          ),
                          maxLines: 2,
                          overflow: TextOverflow.ellipsis,
                        ),
                        const SizedBox(height: 4),
                        Text(
                          booking['address'],
                          style: TextStyle(
                            color: Colors.grey.shade600,
                            fontSize: 14,
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                        const SizedBox(height: 8),
                        Container(
                          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                          decoration: BoxDecoration(
                            color: statusColor.withOpacity(0.1),
                            borderRadius: BorderRadius.circular(4),
                          ),
                          child: Text(
                            statusText,
                            style: TextStyle(
                              color: statusColor,
                              fontWeight: FontWeight.bold,
                              fontSize: 12,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Colors.grey.shade50,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Column(
                  children: [
                    Row(
                      children: [
                        Icon(
                          Icons.calendar_today,
                          size: 16,
                          color: Colors.grey.shade600,
                        ),
                        const SizedBox(width: 8),
                        Text(
                          'Check-in: ${booking['checkIn']}',
                          style: TextStyle(
                            color: Colors.grey.shade600,
                            fontSize: 14,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    Row(
                      children: [
                        Icon(
                          Icons.calendar_today,
                          size: 16,
                          color: Colors.grey.shade600,
                        ),
                        const SizedBox(width: 8),
                        Text(
                          'Check-out: ${booking['checkOut']}',
                          style: TextStyle(
                            color: Colors.grey.shade600,
                            fontSize: 14,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    Row(
                      children: [
                        Icon(
                          Icons.attach_money,
                          size: 16,
                          color: Colors.grey.shade600,
                        ),
                        const SizedBox(width: 8),
                        Text(
                          booking['price'],
                          style: TextStyle(
                            color: Theme.of(context).colorScheme.primary,
                            fontWeight: FontWeight.bold,
                            fontSize: 14,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
              if (isActive) ...
              [
                const SizedBox(height: 16),
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton.icon(
                        onPressed: () => _contactLandlord(booking),
                        icon: const Icon(Icons.phone, size: 16),
                        label: const Text('Contact'),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: ElevatedButton.icon(
                        onPressed: () => _viewBookingDetails(booking),
                        icon: const Icon(Icons.visibility, size: 16),
                        label: const Text('View Details'),
                      ),
                    ),
                  ],
                ),
              ],
              if (!isActive && booking['rating'] != null) ...
              [
                const SizedBox(height: 16),
                Row(
                  children: [
                    Icon(
                      Icons.star,
                      size: 16,
                      color: Colors.amber,
                    ),
                    const SizedBox(width: 4),
                    Text(
                      'You rated: ${booking['rating']}',
                      style: TextStyle(
                        color: Colors.grey.shade600,
                        fontSize: 14,
                      ),
                    ),
                  ],
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }

  void _showBookingDetails(Map<String, dynamic> booking) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => DraggableScrollableSheet(
        initialChildSize: 0.6,
        maxChildSize: 0.9,
        minChildSize: 0.3,
        expand: false,
        builder: (context, scrollController) {
          return Container(
            padding: const EdgeInsets.all(20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Center(
                  child: Container(
                    width: 40,
                    height: 4,
                    decoration: BoxDecoration(
                      color: Colors.grey.shade300,
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                ),
                const SizedBox(height: 20),
                Text(
                  'Booking Details',
                  style: const TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 20),
                Text(
                  booking['propertyName'],
                  style: const TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  booking['address'],
                  style: TextStyle(
                    color: Colors.grey.shade600,
                    fontSize: 16,
                  ),
                ),
                const SizedBox(height: 20),
                _buildDetailRow('Check-in Date', booking['checkIn']),
                _buildDetailRow('Check-out Date', booking['checkOut']),
                _buildDetailRow('Price', booking['price']),
                _buildDetailRow('Landlord', booking['landlord']),
                _buildDetailRow('Phone', booking['landlordPhone']),
                _buildDetailRow('Status', booking['status']),
                const SizedBox(height: 20),
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: () => Navigator.of(context).pop(),
                    child: const Text('Close'),
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildDetailRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 120,
            child: Text(
              label,
              style: TextStyle(
                fontWeight: FontWeight.w500,
                color: Colors.grey.shade600,
              ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
        ],
      ),
    );
  }

  void _contactLandlord(Map<String, dynamic> booking) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Contact Landlord'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Landlord: ${booking['landlord']}'),
            const SizedBox(height: 8),
            Text('Phone: ${booking['landlordPhone']}'),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Close'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.of(context).pop();
              // Here you would typically launch the phone dialer
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(content: Text('Calling ${booking['landlord']}...')),
              );
            },
            child: const Text('Call'),
          ),
        ],
      ),
    );
  }

  void _viewBookingDetails(Map<String, dynamic> booking) {
    Navigator.pushNamed(context, '/booking', arguments: booking);
  }
}

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({super.key});

  @override
  State<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> {
  // Sample user data
  final Map<String, dynamic> userData = {
    'name': 'John Doe',
    'email': 'john.doe@example.com',
    'phone': '+254 700 123456',
    'location': 'Nairobi, Kenya',
    'joinDate': '2023-01-15',
    'totalBookings': 5,
    'savedProperties': 12,
    'isVerified': true,
    'profileImage': '',
  };

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('My Profile'),
        automaticallyImplyLeading: false,
        actions: [
          IconButton(
            icon: const Icon(Icons.settings),
            onPressed: () => _showSettings(),
          ),
        ],
      ),
      body: SingleChildScrollView(
        child: Column(
          children: [
            _buildProfileHeader(),
            const SizedBox(height: 20),
            _buildProfileStats(),
            const SizedBox(height: 20),
            _buildProfileOptions(),
            const SizedBox(height: 20),
            _buildAppOptions(),
          ],
        ),
      ),
    );
  }

  Widget _buildProfileHeader() {
    return Container(
      padding: const EdgeInsets.all(20),
      child: Column(
        children: [
          Stack(
            children: [
              CircleAvatar(
                radius: 50,
                backgroundColor: Colors.grey.shade300,
                child: userData['profileImage'].isEmpty
                    ? Text(
                        userData['name'].toString().split(' ').map((e) => e[0]).join(''),
                        style: const TextStyle(
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      )
                    : null,
              ),
              Positioned(
                bottom: 0,
                right: 0,
                child: Container(
                  decoration: BoxDecoration(
                    color: Theme.of(context).colorScheme.primary,
                    borderRadius: BorderRadius.circular(20),
                  ),
                  child: IconButton(
                    icon: const Icon(Icons.camera_alt, color: Colors.white, size: 20),
                    onPressed: () => _changeProfilePicture(),
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(
                userData['name'],
                style: const TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
              ),
              if (userData['isVerified'])
                const Padding(
                  padding: EdgeInsets.only(left: 8),
                  child: Icon(
                    Icons.verified,
                    color: Colors.blue,
                    size: 20,
                  ),
                ),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            userData['email'],
            style: TextStyle(
              fontSize: 16,
              color: Colors.grey.shade600,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            userData['location'],
            style: TextStyle(
              fontSize: 14,
              color: Colors.grey.shade500,
            ),
          ),
          const SizedBox(height: 16),
          ElevatedButton(
            onPressed: () => _editProfile(),
            child: const Text('Edit Profile'),
          ),
        ],
      ),
    );
  }

  Widget _buildProfileStats() {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 20),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.grey.shade50,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          _buildStatItem('Total Bookings', userData['totalBookings'].toString()),
          _buildStatItem('Saved Properties', userData['savedProperties'].toString()),
          _buildStatItem('Member Since', userData['joinDate'].split('-')[0]),
        ],
      ),
    );
  }

  Widget _buildStatItem(String label, String value) {
    return Column(
      children: [
        Text(
          value,
          style: TextStyle(
            fontSize: 20,
            fontWeight: FontWeight.bold,
            color: Theme.of(context).colorScheme.primary,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: TextStyle(
            fontSize: 12,
            color: Colors.grey.shade600,
          ),
          textAlign: TextAlign.center,
        ),
      ],
    );
  }

  Widget _buildProfileOptions() {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 20),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.1),
            spreadRadius: 1,
            blurRadius: 5,
            offset: const Offset(0, 3),
          ),
        ],
      ),
      child: Column(
        children: [
          _buildProfileOption(
            icon: Icons.person,
            title: 'Personal Information',
            onTap: () => _showPersonalInfo(),
          ),
          _buildProfileOption(
            icon: Icons.security,
            title: 'Security',
            onTap: () => _showSecuritySettings(),
          ),
          _buildProfileOption(
            icon: Icons.payment,
            title: 'Payment Methods',
            onTap: () => _showPaymentMethods(),
          ),
          _buildProfileOption(
            icon: Icons.notifications,
            title: 'Notifications',
            onTap: () => _showNotificationSettings(),
          ),
          _buildProfileOption(
            icon: Icons.help,
            title: 'Help & Support',
            onTap: () => _showHelpSupport(),
          ),
        ],
      ),
    );
  }

  Widget _buildAppOptions() {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 20),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.1),
            spreadRadius: 1,
            blurRadius: 5,
            offset: const Offset(0, 3),
          ),
        ],
      ),
      child: Column(
        children: [
          _buildProfileOption(
            icon: Icons.language,
            title: 'Language',
            subtitle: 'English',
            onTap: () => _showLanguageSettings(),
          ),
          _buildProfileOption(
            icon: Icons.star_rate,
            title: 'Rate the App',
            onTap: () => _rateApp(),
          ),
          _buildProfileOption(
            icon: Icons.share,
            title: 'Share App',
            onTap: () => _shareApp(),
          ),
          _buildProfileOption(
            icon: Icons.info,
            title: 'About',
            onTap: () => _showAbout(),
          ),
          _buildProfileOption(
            icon: Icons.logout,
            title: 'Logout',
            titleColor: Colors.red,
            onTap: () => _logout(),
            showDivider: false,
          ),
        ],
      ),
    );
  }

  Widget _buildProfileOption({
    required IconData icon,
    required String title,
    String? subtitle,
    Color? titleColor,
    required VoidCallback onTap,
    bool showDivider = true,
  }) {
    return Column(
      children: [
        ListTile(
          leading: Icon(
            icon,
            color: titleColor ?? Colors.grey.shade600,
          ),
          title: Text(
            title,
            style: TextStyle(
              fontWeight: FontWeight.w500,
              color: titleColor,
            ),
          ),
          subtitle: subtitle != null
              ? Text(
                  subtitle,
                  style: TextStyle(
                    color: Colors.grey.shade500,
                    fontSize: 12,
                  ),
                )
              : null,
          trailing: const Icon(Icons.arrow_forward_ios, size: 16),
          onTap: onTap,
        ),
        if (showDivider)
          Divider(
            height: 1,
            color: Colors.grey.shade200,
          ),
      ],
    );
  }

  void _changeProfilePicture() {
    showModalBottomSheet(
      context: context,
      builder: (context) => Container(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text(
              'Change Profile Picture',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 20),
            ListTile(
              leading: const Icon(Icons.photo_camera),
              title: const Text('Take Photo'),
              onTap: () {
                Navigator.pop(context);
                // Implement camera functionality
              },
            ),
            ListTile(
              leading: const Icon(Icons.photo_library),
              title: const Text('Choose from Gallery'),
              onTap: () {
                Navigator.pop(context);
                // Implement gallery functionality
              },
            ),
            ListTile(
              leading: const Icon(Icons.delete),
              title: const Text('Remove Photo'),
              onTap: () {
                Navigator.pop(context);
                // Implement remove photo functionality
              },
            ),
          ],
        ),
      ),
    );
  }

  void _editProfile() {
    // Navigate to edit profile screen
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Edit Profile'),
        content: const Text('This will open the profile editing screen.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }

  void _showPersonalInfo() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => DraggableScrollableSheet(
        initialChildSize: 0.6,
        maxChildSize: 0.9,
        minChildSize: 0.3,
        expand: false,
        builder: (context, scrollController) {
          return Container(
            padding: const EdgeInsets.all(20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Center(
                  child: Container(
                    width: 40,
                    height: 4,
                    decoration: BoxDecoration(
                      color: Colors.grey.shade300,
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                ),
                const SizedBox(height: 20),
                const Text(
                  'Personal Information',
                  style: TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 20),
                _buildInfoRow('Full Name', userData['name']),
                _buildInfoRow('Email', userData['email']),
                _buildInfoRow('Phone', userData['phone']),
                _buildInfoRow('Location', userData['location']),
                _buildInfoRow('Join Date', userData['joinDate']),
                _buildInfoRow('Account Status', userData['isVerified'] ? 'Verified' : 'Not Verified'),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 120,
            child: Text(
              label,
              style: TextStyle(
                fontWeight: FontWeight.w500,
                color: Colors.grey.shade600,
              ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
        ],
      ),
    );
  }

  void _showSettings() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Settings'),
        content: const Text('This will open the settings screen.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }

  void _showSecuritySettings() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Security Settings'),
        content: const Text('This will open the security settings screen.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }

  void _showPaymentMethods() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Payment Methods'),
        content: const Text('This will open the payment methods screen.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }

  void _showNotificationSettings() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Notification Settings'),
        content: const Text('This will open the notification settings screen.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }

  void _showHelpSupport() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Help & Support'),
        content: const Text('This will open the help and support screen.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }

  void _showLanguageSettings() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Language Settings'),
        content: const Text('This will open the language settings screen.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }

  void _rateApp() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Rate the App'),
        content: const Text('This will open the app store rating page.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }

  void _shareApp() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Share App'),
        content: const Text('This will open the sharing options.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }

  void _showAbout() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('About'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('App Name: ${dotenv.env['APP_NAME'] ?? 'Kenya Property Rental'}'),
            Text('Version: ${dotenv.env['APP_VERSION'] ?? '1.0.0'}'),
            const SizedBox(height: 10),
            const Text('A modern property rental app for Kenya.'),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }

  void _logout() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Logout'),
        content: const Text('Are you sure you want to logout?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              Navigator.pushReplacementNamed(context, '/login');
            },
            child: const Text('Logout'),
          ),
        ],
      ),
    );
  }
}
