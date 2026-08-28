import 'dart:io';
import 'package:purchases_flutter/purchases_flutter.dart';

class RevenueCatService {
  static final RevenueCatService _instance = RevenueCatService._internal();

  factory RevenueCatService() {
    return _instance;
  }

  RevenueCatService._internal();

  bool _isConfigured = false;

  Future<void> init() async {
    if (_isConfigured) return;

    // TODO: Replace with real API keys from dashboard
    final apiKey = Platform.isIOS ? 'mock_ios_key' : 'mock_android_key';
    
    await Purchases.setLogLevel(LogLevel.debug);
    
    PurchasesConfiguration configuration = PurchasesConfiguration(apiKey);
    await Purchases.configure(configuration);
    
    _isConfigured = true;
  }

  Future<List<Package>> getPackages() async {
    try {
      if (!_isConfigured) {
        await init();
      }
      final offerings = await Purchases.getOfferings();
      if (offerings.current != null && offerings.current!.availablePackages.isNotEmpty) {
        return offerings.current!.availablePackages;
      }
      return [];
    } catch (e) {
      print('Error fetching offerings: $e');
      return [];
    }
  }

  Future<bool> purchasePackage(Package package) async {
    try {
      CustomerInfo customerInfo = await Purchases.purchasePackage(package);
      // Backend webhook will handle actual coin increment. 
      // We just need to check if purchase succeeded.
      return customerInfo.entitlements.all['coins']?.isActive ?? true; 
    } catch (e) {
      print('Purchase failed: $e');
      return false;
    }
  }
}
