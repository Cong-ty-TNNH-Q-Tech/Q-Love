import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import 'package:qlove/core/theme/app_theme.dart';

class ExRatingResultScreen extends StatelessWidget {
  const ExRatingResultScreen({super.key});

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
          AppLocalizations.of(context)?.exRatingResultTitle ?? 'CV Tình trường',
          style: GoogleFonts.outfit(color: Colors.white, fontWeight: FontWeight.bold),
        ),
        centerTitle: true,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Center(
              child: Column(
                children: [
                  const CircleAvatar(
                    radius: 50,
                    backgroundImage: NetworkImage('https://i.pravatar.cc/500'),
                  ),
                  const SizedBox(height: 16),
                  Text(
                    'Alex',
                    style: GoogleFonts.outfit(color: Colors.white, fontSize: 24, fontWeight: FontWeight.bold),
                  ),
                  Text(
                    'Dựa trên 12 đánh giá',
                    style: GoogleFonts.inter(color: Colors.white54, fontSize: 14),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 40),
            Text(
              'Thống kê khía cạnh',
              style: GoogleFonts.outfit(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 16),
            _buildStatBar('Sự chung thủy', 0.9, AppTheme.neonGreen),
            _buildStatBar('Mức độ giao tiếp', 0.7, AppTheme.neonBlue),
            _buildStatBar('Sự quan tâm', 0.8, AppTheme.primaryColor),
            _buildStatBar('Mức độ toxic', 0.2, AppTheme.neonRed),
            
            const SizedBox(height: 40),
            Text(
              'Tags ẩn danh',
              style: GoogleFonts.outfit(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 16),
            Wrap(
              spacing: 8,
              runSpacing: 12,
              children: [
                _buildTag('#Biết lắng nghe', true),
                _buildTag('#Hài hước', true),
                _buildTag('#Trễ giờ', false),
                _buildTag('#Lãng mạn', true),
                _buildTag('#Ghoster', false),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatBar(String title, double value, Color color) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(title, style: GoogleFonts.inter(color: Colors.white, fontWeight: FontWeight.w500)),
              Text('${(value * 10).toStringAsFixed(1)}/10', style: GoogleFonts.inter(color: Colors.white70)),
            ],
          ),
          const SizedBox(height: 8),
          ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: LinearProgressIndicator(
              value: value,
              minHeight: 12,
              backgroundColor: const Color(0xFF1E1E2E),
              valueColor: AlwaysStoppedAnimation<Color>(color),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTag(String text, bool isPositive) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      decoration: BoxDecoration(
        color: isPositive ? AppTheme.neonGreen.withOpacity(0.1) : AppTheme.neonRed.withOpacity(0.1),
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: isPositive ? AppTheme.neonGreen.withOpacity(0.5) : AppTheme.neonRed.withOpacity(0.5)),
      ),
      child: Text(
        text,
        style: GoogleFonts.inter(
          color: isPositive ? AppTheme.neonGreen : AppTheme.neonRed,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }
}
