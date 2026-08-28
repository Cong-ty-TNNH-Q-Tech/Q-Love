import 'package:flutter/material.dart';
import 'package:qlove/features/clan/services/clan_service.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';

class ClanDetailScreen extends StatefulWidget {
  final String clanName;

  const ClanDetailScreen({Key? key, required this.clanName}) : super(key: key);

  @override
  _ClanDetailScreenState createState() => _ClanDetailScreenState();
}

class _ClanDetailScreenState extends State<ClanDetailScreen> {
  final ClanService _clanService = ClanService();
  bool _isInviting = false;

  Future<void> _inviteFriend() async {
    setState(() => _isInviting = true);
    final success = await _clanService.inviteMember('mock-clan-id', 'mock-user-id');
    setState(() => _isInviting = false);

    if (success && mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(AppLocalizations.of(context)!.inviteSuccess)),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0F0F1A),
      appBar: AppBar(
        title: Text(widget.clanName, style: const TextStyle(fontWeight: FontWeight.bold)),
        backgroundColor: Colors.transparent,
        elevation: 0,
        centerTitle: true,
      ),
      body: Padding(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          children: [
            // Clan Shield/Logo
            Container(
              height: 120,
              width: 120,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                gradient: const LinearGradient(
                  colors: [Color(0xFF8A2387), Color(0xFFE94057), Color(0xFFF27121)],
                ),
                boxShadow: [
                  BoxShadow(
                    color: const Color(0xFFE94057).withOpacity(0.5),
                    blurRadius: 20,
                    spreadRadius: 2,
                  )
                ]
              ),
              child: const Icon(Icons.shield, size: 60, color: Colors.white),
            ),
            const SizedBox(height: 24),
            Text(
              widget.clanName,
              style: const TextStyle(color: Colors.white, fontSize: 28, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 8),
            Text(
              AppLocalizations.of(context)!.clanPoints('0'),
              style: const TextStyle(color: Colors.pinkAccent, fontSize: 16, fontWeight: FontWeight.w600),
            ),
            const SizedBox(height: 40),
            
            // Invite Button
            SizedBox(
              width: double.infinity,
              height: 56,
              child: ElevatedButton.icon(
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.white.withOpacity(0.1),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(16),
                    side: const BorderSide(color: Colors.pinkAccent),
                  ),
                ),
                icon: _isInviting 
                  ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(color: Colors.pinkAccent, strokeWidth: 2))
                  : const Icon(Icons.person_add, color: Colors.pinkAccent),
                label: Text(
                  AppLocalizations.of(context)!.inviteMembers,
                  style: const TextStyle(fontSize: 16, color: Colors.white),
                ),
                onPressed: _isInviting ? null : _inviteFriend,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
