import 'package:flutter/material.dart';
import 'dart:math' as math;

enum BottomNavStyle {
  curved,
  floating,
  wave,
  notched,
  glass,
  neon,
}

class CreativeBottomNav extends StatefulWidget {
  final int currentIndex;
  final ValueChanged<int> onTap;
  final List<BottomNavItem> items;
  final BottomNavStyle style;

  const CreativeBottomNav({
    Key? key,
    required this.currentIndex,
    required this.onTap,
    required this.items,
    this.style = BottomNavStyle.curved,
  }) : super(key: key);

  @override
  State<CreativeBottomNav> createState() => _CreativeBottomNavState();
}

class _CreativeBottomNavState extends State<CreativeBottomNav>
    with TickerProviderStateMixin {
  late AnimationController _animationController;
  late Animation<double> _scaleAnimation;

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      duration: const Duration(milliseconds: 300),
      vsync: this,
    );
    _scaleAnimation = Tween<double>(
      begin: 0.8,
      end: 1.0,
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: Curves.elasticOut,
    ));
  }

  @override
  void dispose() {
    _animationController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    switch (widget.style) {
      case BottomNavStyle.curved:
        return _buildCurvedNav();
      case BottomNavStyle.floating:
        return _buildFloatingNav();
      case BottomNavStyle.wave:
        return _buildWaveNav();
      case BottomNavStyle.notched:
        return _buildNotchedNav();
      case BottomNavStyle.glass:
        return _buildGlassNav();
      case BottomNavStyle.neon:
        return _buildNeonNav();
    }
  }

  Widget _buildCurvedNav() {
    return Container(
      height: 80,
      margin: const EdgeInsets.all(16),
      child: Stack(
        children: [
          // Main curved background
          Container(
            height: 80,
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(40),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withOpacity(0.15),
                  blurRadius: 20,
                  offset: const Offset(0, 5),
                ),
                BoxShadow(
                  color: Colors.black.withOpacity(0.08),
                  blurRadius: 40,
                  offset: const Offset(0, 10),
                ),
              ],
            ),
          ),
          // Navigation items
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: widget.items.asMap().entries.map((entry) {
              return _buildNavItem(entry.key, entry.value);
            }).toList(),
          ),
        ],
      ),
    );
  }

  Widget _buildFloatingNav() {
    return Container(
      height: 70,
      margin: const EdgeInsets.symmetric(horizontal: 50, vertical: 20),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.primary,
        borderRadius: BorderRadius.circular(35),
        boxShadow: [
          BoxShadow(
            color: Theme.of(context).colorScheme.primary.withOpacity(0.4),
            blurRadius: 20,
            offset: const Offset(0, 10),
          ),
        ],
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: widget.items.asMap().entries.map((entry) {
          final isSelected = widget.currentIndex == entry.key;
          return GestureDetector(
            onTap: () => widget.onTap(entry.key),
            child: AnimatedContainer(
              duration: const Duration(milliseconds: 300),
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
              decoration: BoxDecoration(
                color: isSelected ? Colors.white : Colors.transparent,
                borderRadius: BorderRadius.circular(25),
              ),
              child: Icon(
                isSelected ? entry.value.activeIcon : entry.value.icon,
                color: isSelected ? Theme.of(context).colorScheme.primary : Colors.white,
                size: 24,
              ),
            ),
          );
        }).toList(),
      ),
    );
  }

  Widget _buildWaveNav() {
    return Container(
      height: 80,
      child: Stack(
        children: [
          CustomPaint(
            size: Size(MediaQuery.of(context).size.width, 80),
            painter: WavePainter(),
          ),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: widget.items.asMap().entries.map((entry) {
              return _buildWaveNavItem(entry.key, entry.value);
            }).toList(),
          ),
        ],
      ),
    );
  }

  Widget _buildNotchedNav() {
    return Container(
      height: 80,
      child: Stack(
        children: [
          CustomPaint(
            size: Size(MediaQuery.of(context).size.width, 80),
            painter: NotchedPainter(widget.currentIndex, widget.items.length),
          ),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: widget.items.asMap().entries.map((entry) {
              return _buildNotchedNavItem(entry.key, entry.value);
            }).toList(),
          ),
        ],
      ),
    );
  }

  Widget _buildGlassNav() {
    return Container(
      height: 80,
      margin: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(25),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.1),
            blurRadius: 20,
            offset: const Offset(0, 10),
          ),
        ],
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(25),
        child: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: [
                Colors.white.withOpacity(0.8),
                Colors.white.withOpacity(0.6),
              ],
            ),
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: widget.items.asMap().entries.map((entry) {
              return _buildGlassNavItem(entry.key, entry.value);
            }).toList(),
          ),
        ),
      ),
    );
  }

  Widget _buildNeonNav() {
    return Container(
      height: 80,
      margin: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.black,
        borderRadius: BorderRadius.circular(25),
        border: Border.all(
          color: Theme.of(context).colorScheme.primary.withOpacity(0.3),
          width: 1,
        ),
        boxShadow: [
          BoxShadow(
            color: Theme.of(context).colorScheme.primary.withOpacity(0.2),
            blurRadius: 20,
            offset: const Offset(0, 10),
          ),
        ],
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: widget.items.asMap().entries.map((entry) {
          return _buildNeonNavItem(entry.key, entry.value);
        }).toList(),
      ),
    );
  }

  Widget _buildNavItem(int index, BottomNavItem item) {
    final isSelected = widget.currentIndex == index;
    final primaryColor = Theme.of(context).colorScheme.primary;
    
    return GestureDetector(
      onTap: () {
        widget.onTap(index);
        _animationController.reset();
        _animationController.forward();
      },
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 300),
        curve: Curves.easeInOut,
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        decoration: BoxDecoration(
          color: isSelected ? primaryColor.withOpacity(0.1) : Colors.transparent,
          borderRadius: BorderRadius.circular(25),
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            ScaleTransition(
              scale: isSelected ? _scaleAnimation : const AlwaysStoppedAnimation(1.0),
              child: AnimatedContainer(
                duration: const Duration(milliseconds: 300),
                curve: Curves.easeInOut,
                padding: const EdgeInsets.all(6),
                decoration: BoxDecoration(
                  color: isSelected ? primaryColor : Colors.transparent,
                  borderRadius: BorderRadius.circular(20),
                  boxShadow: isSelected
                      ? [
                          BoxShadow(
                            color: primaryColor.withOpacity(0.3),
                            blurRadius: 10,
                            offset: const Offset(0, 3),
                          ),
                        ]
                      : [],
                ),
                child: Icon(
                  isSelected ? item.activeIcon : item.icon,
                  color: isSelected ? Colors.white : Colors.grey.shade600,
                  size: isSelected ? 22 : 20,
                ),
              ),
            ),
            const SizedBox(height: 4),
            AnimatedDefaultTextStyle(
              duration: const Duration(milliseconds: 300),
              style: TextStyle(
                fontSize: isSelected ? 12 : 10,
                fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
                color: isSelected ? primaryColor : Colors.grey.shade600,
              ),
              child: Text(item.label),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildWaveNavItem(int index, BottomNavItem item) {
    final isSelected = widget.currentIndex == index;
    
    return GestureDetector(
      onTap: () => widget.onTap(index),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 300),
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              isSelected ? item.activeIcon : item.icon,
              color: isSelected ? Theme.of(context).colorScheme.primary : Colors.grey.shade600,
              size: isSelected ? 28 : 24,
            ),
            const SizedBox(height: 4),
            Text(
              item.label,
              style: TextStyle(
                fontSize: 10,
                fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
                color: isSelected ? Theme.of(context).colorScheme.primary : Colors.grey.shade600,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildNotchedNavItem(int index, BottomNavItem item) {
    final isSelected = widget.currentIndex == index;
    
    return GestureDetector(
      onTap: () => widget.onTap(index),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 300),
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        decoration: BoxDecoration(
          color: isSelected ? Theme.of(context).colorScheme.primary : Colors.transparent,
          borderRadius: BorderRadius.circular(30),
          boxShadow: isSelected
              ? [
                  BoxShadow(
                    color: Theme.of(context).colorScheme.primary.withOpacity(0.3),
                    blurRadius: 10,
                    offset: const Offset(0, 5),
                  ),
                ]
              : [],
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              isSelected ? item.activeIcon : item.icon,
              color: isSelected ? Colors.white : Colors.grey.shade600,
              size: isSelected ? 28 : 24,
            ),
            const SizedBox(height: 4),
            Text(
              item.label,
              style: TextStyle(
                fontSize: 10,
                fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
                color: isSelected ? Colors.white : Colors.grey.shade600,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildGlassNavItem(int index, BottomNavItem item) {
    final isSelected = widget.currentIndex == index;
    
    return GestureDetector(
      onTap: () => widget.onTap(index),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 300),
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        decoration: BoxDecoration(
          color: isSelected ? Colors.white.withOpacity(0.3) : Colors.transparent,
          borderRadius: BorderRadius.circular(20),
          border: isSelected ? Border.all(color: Colors.white.withOpacity(0.5)) : null,
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              isSelected ? item.activeIcon : item.icon,
              color: isSelected ? Theme.of(context).colorScheme.primary : Colors.grey.shade700,
              size: isSelected ? 28 : 24,
            ),
            const SizedBox(height: 4),
            Text(
              item.label,
              style: TextStyle(
                fontSize: 10,
                fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
                color: isSelected ? Theme.of(context).colorScheme.primary : Colors.grey.shade700,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildNeonNavItem(int index, BottomNavItem item) {
    final isSelected = widget.currentIndex == index;
    final primaryColor = Theme.of(context).colorScheme.primary;
    
    return GestureDetector(
      onTap: () => widget.onTap(index),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 300),
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(20),
          boxShadow: isSelected
              ? [
                  BoxShadow(
                    color: primaryColor.withOpacity(0.5),
                    blurRadius: 20,
                    spreadRadius: 2,
                  ),
                ]
              : [],
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              isSelected ? item.activeIcon : item.icon,
              color: isSelected ? primaryColor : Colors.grey.shade600,
              size: isSelected ? 28 : 24,
            ),
            const SizedBox(height: 4),
            Text(
              item.label,
              style: TextStyle(
                fontSize: 10,
                fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
                color: isSelected ? primaryColor : Colors.grey.shade600,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class BottomNavItem {
  final IconData icon;
  final IconData activeIcon;
  final String label;

  const BottomNavItem({
    required this.icon,
    required this.activeIcon,
    required this.label,
  });
}

class WavePainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = Colors.white
      ..style = PaintingStyle.fill;

    final path = Path();
    path.moveTo(0, 20);
    
    // Create wave effect
    for (double i = 0; i < size.width; i++) {
      path.lineTo(i, 20 + 10 * math.sin(i * 0.02));
    }
    
    path.lineTo(size.width, size.height);
    path.lineTo(0, size.height);
    path.close();

    canvas.drawPath(path, paint);
    
    // Add shadow
    canvas.drawShadow(path, Colors.black.withOpacity(0.2), 10, false);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}

class NotchedPainter extends CustomPainter {
  final int selectedIndex;
  final int itemCount;

  NotchedPainter(this.selectedIndex, this.itemCount);

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = Colors.white
      ..style = PaintingStyle.fill;

    final path = Path();
    final notchWidth = size.width / itemCount;
    final notchCenter = notchWidth * selectedIndex + notchWidth / 2;

    path.moveTo(0, 20);
    
    // Create notch for selected item
    path.lineTo(notchCenter - 30, 20);
    path.quadraticBezierTo(notchCenter - 15, 0, notchCenter, 0);
    path.quadraticBezierTo(notchCenter + 15, 0, notchCenter + 30, 20);
    
    path.lineTo(size.width, 20);
    path.lineTo(size.width, size.height);
    path.lineTo(0, size.height);
    path.close();

    canvas.drawPath(path, paint);
    
    // Add shadow
    canvas.drawShadow(path, Colors.black.withOpacity(0.2), 10, false);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => true;
}
