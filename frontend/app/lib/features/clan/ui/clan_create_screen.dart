import 'package:flutter/material.dart';
import 'package:qlove/features/clan/services/clan_service.dart';
import 'package:qlove/features/clan/services/profanity_service.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import 'package:qlove/features/clan/ui/clan_detail_screen.dart';

class ClanCreateScreen extends StatefulWidget {
  const ClanCreateScreen({Key? key}) : super(key: key);

  @override
  _ClanCreateScreenState createState() => _ClanCreateScreenState();
}

class _ClanCreateScreenState extends State<ClanCreateScreen> {
  final TextEditingController _nameController = TextEditingController();
  final ProfanityService _profanityService = ProfanityService();
  final ClanService _clanService = ClanService();
  
  bool _isLoading = false;
  String? _errorText;

  Future<void> _handleCreateClan() async {
    final name = _nameController.text.trim();
    if (name.isEmpty) return;

    if (_profanityService.containsProfanity(name)) {
      setState(() {
        _errorText = AppLocalizations.of(context)!.profanityWarning;
      });
      return;
    }

    setState(() {
      _errorText = null;
      _isLoading = true;
    });

    final success = await _clanService.createClan(name);
    
    if (!mounted) return;
    setState(() {
      _isLoading = false;
    });

    if (success) {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(
          builder: (context) => ClanDetailScreen(clanName: name),
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0F0F1A), // Gen-Z Dark-first
      appBar: AppBar(
        title: Text(AppLocalizations.of(context)!.createClan, style: const TextStyle(fontWeight: FontWeight.bold)),
        backgroundColor: Colors.transparent,
        elevation: 0,
        centerTitle: true,
      ),
      body: Padding(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const Icon(Icons.shield, size: 80, color: Colors.pinkAccent),
            const SizedBox(height: 24),
            Text(
              AppLocalizations.of(context)!.clanName,
              style: const TextStyle(color: Colors.white70, fontSize: 16),
            ),
            const SizedBox(height: 8),
            TextField(
              controller: _nameController,
              style: const TextStyle(color: Colors.white),
              decoration: InputDecoration(
                filled: true,
                fillColor: Colors.white.withOpacity(0.05),
                hintText: AppLocalizations.of(context)!.clanNameHint,
                hintStyle: TextStyle(color: Colors.white.withOpacity(0.3)),
                errorText: _errorText,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(16),
                  borderSide: BorderSide.none,
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(16),
                  borderSide: const BorderSide(color: Colors.pinkAccent, width: 2),
                ),
              ),
              onChanged: (_) {
                if (_errorText != null) {
                  setState(() => _errorText = null);
                }
              },
            ),
            const Spacer(),
            SizedBox(
              height: 56,
              child: ElevatedButton(
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.pinkAccent,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(16),
                  ),
                ),
                onPressed: _isLoading ? null : _handleCreateClan,
                child: _isLoading
                    ? const CircularProgressIndicator(color: Colors.white)
                    : Text(
                        AppLocalizations.of(context)!.createClan,
                        style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                      ),
              ),
            ),
            const SizedBox(height: 24),
          ],
        ),
      ),
    );
  }

  @override
  void dispose() {
    _nameController.dispose();
    super.dispose();
  }
}
