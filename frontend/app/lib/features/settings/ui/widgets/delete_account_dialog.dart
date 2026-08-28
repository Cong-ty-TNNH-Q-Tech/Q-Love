import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';

class DeleteAccountDialog extends StatelessWidget {
  const DeleteAccountDialog({super.key});

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      backgroundColor: const Color(0xFF1E1E2E),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(24)),
      title: Row(
        children: [
          const Icon(Icons.warning_amber_rounded, color: Color(0xFFFF3D6B), size: 28),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              AppLocalizations.of(context)?.deleteAccountWarningTitle ?? 'Xóa tài khoản vĩnh viễn?',
              style: GoogleFonts.outfit(
                color: Colors.white,
                fontWeight: FontWeight.bold,
                fontSize: 20,
              ),
            ),
          ),
        ],
      ),
      content: Text(
        AppLocalizations.of(context)?.deleteAccountWarningDesc ?? 
        'Bạn sẽ mất toàn bộ kết nối, tin nhắn và lịch sử. Không thể hoàn tác hành động này.',
        style: GoogleFonts.inter(
          color: Colors.white70,
          fontSize: 14,
          height: 1.5,
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context, false),
          child: Text(
            AppLocalizations.of(context)?.cancel ?? 'Hủy',
            style: GoogleFonts.inter(
              color: Colors.white54,
              fontWeight: FontWeight.w500,
            ),
          ),
        ),
        ElevatedButton(
          onPressed: () {
            // API call to Soft Delete user goes here
            // e.g. context.read<AuthBloc>().add(DeleteAccountRequested());
            Navigator.pop(context, true);
          },
          style: ElevatedButton.styleFrom(
            backgroundColor: const Color(0xFFFF3D6B),
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
          ),
          child: Text(
            AppLocalizations.of(context)?.confirmDelete ?? 'Xóa vĩnh viễn',
            style: GoogleFonts.inter(
              color: Colors.white,
              fontWeight: FontWeight.bold,
            ),
          ),
        ),
      ],
    );
  }
}
