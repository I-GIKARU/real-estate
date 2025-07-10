import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';

// Import new services and models
import 'services/auth_service.dart';
import 'services/property_service.dart';
import 'provider/user_provider.dart';
import 'screens/splash_screen.dart';
import 'screens/login_screen.dart';
import 'screens/register_screen.dart';
import 'screens/home_screen.dart';
import 'screens/property_details_screen.dart';
import 'screens/booking_screen.dart';


void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await dotenv.load(fileName: ".env");
  runApp(PropertyRentalApp());
}

class PropertyRentalApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => UserProvider()),
        Provider(create: (_) => AuthService()),
        Provider(create: (_) => PropertyService()),
      ],
      child: MaterialApp(
        title: 'Kenya Property Rental',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(
          primarySwatch: Colors.green,
          primaryColor: Colors.green[700],
          visualDensity: VisualDensity.adaptivePlatformDensity,
          appBarTheme: AppBarTheme(
            backgroundColor: Colors.green[700],
            foregroundColor: Colors.white,
            elevation: 0,
          ),
          elevatedButtonTheme: ElevatedButtonThemeData(
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.green[700],
              foregroundColor: Colors.white,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
            ),
          ),
        ),
home: SplashScreen(),
        routes: {
          '/login': (context) => LoginScreen(),
          '/register': (context) => RegisterScreen(),
          '/home': (context) => HomeScreen(),
          '/tenant-home': (context) => HomeScreen(),
          '/landlord-home': (context) => HomeScreen(),
'/properties': (context) => PropertyListScreen(),
          '/property-detail': (context) => PropertyDetailsScreen(),
          '/property_details': (context) => PropertyDetailsScreen(),
          '/booking': (context) => BookingScreen(),
        },
      ),
    );
  }
}
