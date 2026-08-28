import 'package:flutter/material.dart';
import 'package:purchases_flutter/purchases_flutter.dart';
import 'package:qlove/features/store/services/revenuecat_service.dart';
import 'package:qlove/features/store/ui/widgets/coin_package_card.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';

class CoinStoreScreen extends StatefulWidget {
  const CoinStoreScreen({Key? key}) : super(key: key);

  @override
  _CoinStoreScreenState createState() => _CoinStoreScreenState();
}

class _CoinStoreScreenState extends State<CoinStoreScreen> {
  final RevenueCatService _revenueCatService = RevenueCatService();
  List<Package> _packages = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _fetchPackages();
  }

  Future<void> _fetchPackages() async {
    final packages = await _revenueCatService.getPackages();
    setState(() {
      _packages = packages;
      _isLoading = false;
    });
  }

  Future<void> _buyPackage(Package package) async {
    setState(() => _isLoading = true);
    final success = await _revenueCatService.purchasePackage(package);
    setState(() => _isLoading = false);

    if (success) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(AppLocalizations.of(context)!.purchaseSuccess)),
      );
    } else {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(AppLocalizations.of(context)!.purchaseFailed)),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0F0F1A), // Gen-Z Dark-first
      appBar: AppBar(
        title: Text(AppLocalizations.of(context)!.storeTitle, style: const TextStyle(fontWeight: FontWeight.bold)),
        backgroundColor: Colors.transparent,
        elevation: 0,
        centerTitle: true,
      ),
      body: _isLoading 
        ? const Center(child: CircularProgressIndicator(color: Colors.pinkAccent))
        : Column(
            children: [
              const SizedBox(height: 20),
              // Header Wallet Balance
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 30, vertical: 20),
                decoration: BoxDecoration(
                  gradient: const LinearGradient(
                    colors: [Color(0xFF8A2387), Color(0xFFE94057), Color(0xFFF27121)],
                  ),
                  borderRadius: BorderRadius.circular(20),
                  boxShadow: [
                    BoxShadow(
                      color: const Color(0xFFE94057).withOpacity(0.5),
                      blurRadius: 20,
                      spreadRadius: 2,
                    )
                  ]
                ),
                child: Column(
                  children: [
                    Text(AppLocalizations.of(context)!.currentCoins, style: const TextStyle(color: Colors.white70, fontSize: 16)),
                    const SizedBox(height: 8),
                    const Text('1,250', style: TextStyle(color: Colors.white, fontSize: 36, fontWeight: FontWeight.w900)),
                  ],
                ),
              ),
              const SizedBox(height: 30),
              Text(AppLocalizations.of(context)!.selectPackage, style: const TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold)),
              const SizedBox(height: 16),
              
              Expanded(
                child: _packages.isEmpty 
                  ? Center(child: Text(AppLocalizations.of(context)!.noPackages, style: const TextStyle(color: Colors.white54)))
                  : ListView.builder(
                      padding: const EdgeInsets.symmetric(horizontal: 16),
                      itemCount: _packages.length,
                      itemBuilder: (context, index) {
                        final package = _packages[index];
                        return CoinPackageCard(
                          package: package,
                          onTap: () => _buyPackage(package),
                        );
                      },
                    ),
              ),
            ],
          ),
    );
  }
}
