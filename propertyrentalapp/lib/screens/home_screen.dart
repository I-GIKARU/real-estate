import 'package:flutter/material.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _selectedIndex = 0;
  
  final List<Widget> _screens = [
    const PropertyListScreen(),
    const SavedPropertiesScreen(),
    const BookingsScreen(),
    const ProfileScreen(),
  ];
  
  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _screens[_selectedIndex],
      bottomNavigationBar: BottomNavigationBar(
        type: BottomNavigationBarType.fixed,
        currentIndex: _selectedIndex,
        onTap: _onItemTapped,
        selectedItemColor: Theme.of(context).colorScheme.primary,
        unselectedItemColor: Colors.grey,
        items: const [
          BottomNavigationBarItem(
            icon: Icon(Icons.home),
            label: 'Explore',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.favorite),
            label: 'Saved',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.calendar_today),
            label: 'Bookings',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.person),
            label: 'Profile',
          ),
        ],
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

class SavedPropertiesScreen extends StatelessWidget {
  const SavedPropertiesScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Saved Properties'),
      ),
      body: const Center(
        child: Text('Saved Properties Screen'),
      ),
    );
  }
}

class BookingsScreen extends StatelessWidget {
  const BookingsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('My Bookings'),
      ),
      body: const Center(
        child: Text('Bookings Screen'),
      ),
    );
  }
}

class ProfileScreen extends StatelessWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('My Profile'),
      ),
      body: const Center(
        child: Text('Profile Screen'),
      ),
    );
  }
}
