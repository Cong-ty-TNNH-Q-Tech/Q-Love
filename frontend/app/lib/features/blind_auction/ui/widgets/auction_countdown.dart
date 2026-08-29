import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:qlove/core/theme/app_theme.dart';

class AuctionCountdown extends StatelessWidget {
  final Duration timeRemaining;

  const AuctionCountdown({
    super.key,
    required this.timeRemaining,
  });

  @override
  Widget build(BuildContext context) {
    final hours = timeRemaining.inHours.toString().padLeft(2, '0');
    final minutes = (timeRemaining.inMinutes % 60).toString().padLeft(2, '0');
    final seconds = (timeRemaining.inSeconds % 60).toString().padLeft(2, '0');

    final isUrgent = timeRemaining.inHours < 1;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
      decoration: BoxDecoration(
        color: const Color(0xFF1E1E2E),
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          if (isUrgent)
            BoxShadow(
              color: AppTheme.neonRed.withOpacity(0.5),
              blurRadius: 20,
              spreadRadius: 2,
            ),
        ],
      ),
      child: Column(
        children: [
          Text(
            '$hours:$minutes:$seconds',
            style: GoogleFonts.dseg14Classic( // using a digital font or fallback
              fontSize: 48,
              fontWeight: FontWeight.bold,
              color: isUrgent ? AppTheme.neonRed : Colors.white,
            ).copyWith(fontFamily: 'Inter'), // Fallback if DSEG is not available
          ),
          const SizedBox(height: 8),
          Text(
            'THỜI GIAN CÒN LẠI',
            style: GoogleFonts.inter(
              fontSize: 12,
              fontWeight: FontWeight.w600,
              color: Colors.white54,
              letterSpacing: 2,
            ),
          ),
        ],
      ),
    );
  }
}
