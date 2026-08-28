import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';

class TopHuntersList extends StatelessWidget {
  const TopHuntersList({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: const BoxDecoration(
        color: Color(0xFF1E1E2E),
        borderRadius: BorderRadius.vertical(top: Radius.circular(32)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(Icons.leaderboard, color: Color(0xFFFF3D6B)),
              const SizedBox(width: 8),
              Text(
                AppLocalizations.of(context)?.topHunters ?? 'Top 5 Thợ Săn',
                style: GoogleFonts.outfit(
                  color: Colors.white,
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          Expanded(
            child: ListView.builder(
              itemCount: 5,
              itemBuilder: (context, index) {
                return _buildHunterTile(index);
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildHunterTile(int index) {
    // Generate a mock hunter for UI preview
    final isTop = index == 0;
    
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: isTop ? const Color(0xFFFF3D6B).withOpacity(0.1) : Colors.white.withOpacity(0.03),
        borderRadius: BorderRadius.circular(16),
        border: isTop ? Border.all(color: const Color(0xFFFF3D6B).withOpacity(0.3)) : null,
      ),
      child: Row(
        children: [
          Container(
            width: 32,
            height: 32,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: isTop ? const Color(0xFFFF3D6B) : Colors.white12,
            ),
            child: Text(
              '${index + 1}',
              style: GoogleFonts.outfit(
                color: Colors.white,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
          const SizedBox(width: 12),
          CircleAvatar(
            radius: 20,
            backgroundColor: Colors.grey[800],
            backgroundImage: NetworkImage('https://i.pravatar.cc/100?img=${index + 10}'),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              'Hunter ${index + 1}',
              style: GoogleFonts.inter(
                color: Colors.white,
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
          Text(
            '??? Xu', // Blind auction hides exact bid
            style: GoogleFonts.outfit(
              color: const Color(0xFFFF3D6B),
              fontWeight: FontWeight.bold,
              fontSize: 16,
            ),
          ),
        ],
      ),
    );
  }
}
