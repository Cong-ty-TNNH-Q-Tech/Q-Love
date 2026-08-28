import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';

class BidInputArea extends StatefulWidget {
  final ValueChanged<int> onBidChanged;
  final VoidCallback onPlaceBid;

  const BidInputArea({
    super.key,
    required this.onBidChanged,
    required this.onPlaceBid,
  });

  @override
  State<BidInputArea> createState() => _BidInputAreaState();
}

class _BidInputAreaState extends State<BidInputArea> {
  double _bidAmount = 1000;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: const Color(0xFF1E1E2E).withOpacity(0.8),
        borderRadius: BorderRadius.circular(24),
        border: Border.all(color: Colors.white10),
      ),
      child: Column(
        children: [
          Text(
            AppLocalizations.of(context)?.yourBid ?? 'Giá thầu của bạn',
            style: GoogleFonts.inter(
              color: Colors.white54,
              fontSize: 14,
            ),
          ),
          const SizedBox(height: 16),
          TweenAnimationBuilder<double>(
            tween: Tween<double>(begin: _bidAmount, end: _bidAmount),
            duration: const Duration(milliseconds: 300),
            builder: (context, value, child) {
              return Text(
                '${value.toInt()} Xu',
                style: GoogleFonts.outfit(
                  color: Colors.white,
                  fontSize: 48,
                  fontWeight: FontWeight.bold,
                  shadows: [
                    Shadow(
                      color: const Color(0xFFFF3D6B).withOpacity(0.5),
                      blurRadius: 10,
                    ),
                  ],
                ),
              );
            },
          ),
          const SizedBox(height: 24),
          SliderTheme(
            data: SliderThemeData(
              activeTrackColor: const Color(0xFFFF3D6B),
              inactiveTrackColor: Colors.white12,
              thumbColor: Colors.white,
              overlayColor: const Color(0xFFFF3D6B).withOpacity(0.2),
              trackHeight: 8,
            ),
            child: Slider(
              value: _bidAmount,
              min: 100,
              max: 100000,
              divisions: 100,
              onChanged: (value) {
                setState(() {
                  _bidAmount = value;
                });
                widget.onBidChanged(value.toInt());
              },
            ),
          ),
          const SizedBox(height: 24),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: widget.onPlaceBid,
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFFF3D6B),
                padding: const EdgeInsets.symmetric(vertical: 16),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
                elevation: 10,
                shadowColor: const Color(0xFFFF3D6B).withOpacity(0.5),
              ),
              child: Text(
                AppLocalizations.of(context)?.placeBid ?? 'ĐẤU THẦU',
                style: GoogleFonts.inter(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  letterSpacing: 1,
                  color: Colors.white,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
