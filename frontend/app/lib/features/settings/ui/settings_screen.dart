import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';

import 'widgets/delete_account_dialog.dart';

class SettingsScreen extends StatelessWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0D0D14),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.white),
          onPressed: () => Navigator.pop(context),
        ),
        title: Text(
          AppLocalizations.of(context)?.settingsTitle ?? 'Cài đặt',
          style: GoogleFonts.outfit(color: Colors.white, fontWeight: FontWeight.bold),
        ),
        centerTitle: true,
      ),
      body: ListView(
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
        children: [
          _buildSettingsGroup(
            context,
            title: AppLocalizations.of(context)?.settingsTitle ?? 'Chung',
            children: [
              _buildListTile(
                icon: Icons.notifications_none_rounded,
                title: AppLocalizations.of(context)?.notifications ?? 'Thông báo',
                onTap: () {},
              ),
              _buildListTile(
                icon: Icons.language_rounded,
                title: AppLocalizations.of(context)?.language ?? 'Ngôn ngữ',
                onTap: () {},
              ),
            ],
          ),
          const SizedBox(height: 32),
          _buildSettingsGroup(
            context,
            title: 'Quyền riêng tư',
            children: [
              _buildListTile(
                icon: Icons.shield_outlined,
                title: 'Chính sách bảo mật',
                onTap: () {},
              ),
              _buildListTile(
                icon: Icons.delete_outline_rounded,
                title: AppLocalizations.of(context)?.deleteAccountBtn ?? 'Xóa tài khoản',
                titleColor: const Color(0xFFFF3D6B),
                iconColor: const Color(0xFFFF3D6B),
                onTap: () async {
                  final result = await showDialog<bool>(
                    context: context,
                    builder: (context) => const DeleteAccountDialog(),
                  );
                  if (result == true) {
                    // Soft delete triggered, navigate out
                    // Navigator.pushAndRemoveUntil(...)
                  }
                },
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildSettingsGroup(BuildContext context, {required String title, required List<Widget> children}) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.only(left: 8.0, bottom: 8.0),
          child: Text(
            title.toUpperCase(),
            style: GoogleFonts.inter(
              color: Colors.white54,
              fontSize: 12,
              fontWeight: FontWeight.bold,
              letterSpacing: 1.2,
            ),
          ),
        ),
        Container(
          decoration: BoxDecoration(
            color: const Color(0xFF1E1E2E),
            borderRadius: BorderRadius.circular(20),
          ),
          child: Column(
            children: children,
          ),
        ),
      ],
    );
  }

  Widget _buildListTile({
    required IconData icon,
    required String title,
    required VoidCallback onTap,
    Color titleColor = Colors.white,
    Color iconColor = Colors.white70,
  }) {
    return ListTile(
      onTap: onTap,
      leading: Icon(icon, color: iconColor),
      title: Text(
        title,
        style: GoogleFonts.inter(
          color: titleColor,
          fontWeight: FontWeight.w500,
        ),
      ),
      trailing: const Icon(Icons.chevron_right_rounded, color: Colors.white24),
    );
  }
}
