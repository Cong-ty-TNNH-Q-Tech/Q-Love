import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';

class UnmatchRatingBottomSheet extends StatefulWidget {
  const UnmatchRatingBottomSheet({super.key});

  static void show(BuildContext context) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => const UnmatchRatingBottomSheet(),
    );
  }

  @override
  State<UnmatchRatingBottomSheet> createState() => _UnmatchRatingBottomSheetState();
}

class _UnmatchRatingBottomSheetState extends State<UnmatchRatingBottomSheet> {
  int _rating = 0;
  final Set<String> _selectedTags = {};

  final List<String> _availableTags = [
    'Toxic',
    'Bạo lực',
    'Lừa dối',
    'Vô tâm',
    'Ghoster',
    'Hài hước',
    'Lắng nghe',
    'Ga lăng',
  ];

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: const BoxDecoration(
        color: Color(0xFF1E1E2E),
        borderRadius: BorderRadius.vertical(top: Radius.circular(32)),
      ),
      padding: EdgeInsets.only(
        top: 32,
        left: 24,
        right: 24,
        bottom: MediaQuery.of(context).viewInsets.bottom + 32,
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          Container(
            width: 48,
            height: 4,
            decoration: BoxDecoration(
              color: Colors.white24,
              borderRadius: BorderRadius.circular(4),
            ),
          ),
          const SizedBox(height: 32),
          Text(
            AppLocalizations.of(context)?.unmatchRatingTitle ?? 'Đánh giá Ex',
            style: GoogleFonts.outfit(color: Colors.white, fontSize: 24, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 8),
          Text(
            AppLocalizations.of(context)?.unmatchRatingSubtitle ?? 'Để lại đánh giá ẩn danh cho người này.',
            textAlign: TextAlign.center,
            style: GoogleFonts.inter(color: Colors.white54, fontSize: 14),
          ),
          const SizedBox(height: 32),
          
          // Rating Stars
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: List.generate(5, (index) {
              return IconButton(
                icon: Icon(
                  index < _rating ? Icons.star_rounded : Icons.star_outline_rounded,
                  color: index < _rating ? const Color(0xFFFFD700) : Colors.white24,
                  size: 40,
                ),
                onPressed: () {
                  setState(() {
                    _rating = index + 1;
                  });
                },
              );
            }),
          ),
          
          const SizedBox(height: 32),
          Wrap(
            spacing: 8,
            runSpacing: 12,
            alignment: WrapAlignment.center,
            children: _availableTags.map((tag) {
              final isSelected = _selectedTags.contains(tag);
              return FilterChip(
                label: Text(tag),
                selected: isSelected,
                onSelected: (selected) {
                  setState(() {
                    if (selected) {
                      _selectedTags.add(tag);
                    } else {
                      _selectedTags.remove(tag);
                    }
                  });
                },
                backgroundColor: Colors.white.withOpacity(0.05),
                selectedColor: const Color(0xFFFF3D6B).withOpacity(0.2),
                checkmarkColor: const Color(0xFFFF3D6B),
                labelStyle: GoogleFonts.inter(
                  color: isSelected ? const Color(0xFFFF3D6B) : Colors.white70,
                  fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                ),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(20),
                  side: BorderSide(
                    color: isSelected ? const Color(0xFFFF3D6B).withOpacity(0.5) : Colors.transparent,
                  ),
                ),
              );
            }).toList(),
          ),
          
          const SizedBox(height: 40),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: _rating > 0
                  ? () {
                      // Submit rating
                      Navigator.pop(context);
                    }
                  : null,
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFFF3D6B),
                disabledBackgroundColor: Colors.white12,
                padding: const EdgeInsets.symmetric(vertical: 16),
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
              ),
              child: Text(
                AppLocalizations.of(context)?.submitRating ?? 'Gửi đánh giá',
                style: GoogleFonts.inter(
                  color: _rating > 0 ? Colors.white : Colors.white30,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
