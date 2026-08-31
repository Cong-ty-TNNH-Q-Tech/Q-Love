import 'package:flutter/material.dart';
import 'package:qlove/core/theme/app_theme.dart';
import 'package:qlove/widgets/glitch_text.dart';
import 'package:qlove/widgets/radar_pulse.dart';

class PurgeLobbyScreen extends StatelessWidget {
  const PurgeLobbyScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF080810),
      body: Stack(
        children: [
          // Dark eerie background with subtle gradient
          Positioned.fill(
            child: Container(
              decoration: const BoxDecoration(
                gradient: RadialGradient(
                  center: Alignment.center,
                  radius: 1.5,
                  colors: [
                    Color(0xFF151020),
                    Color(0xFF080810),
                  ],
                ),
              ),
            ),
          ),
          
          SafeArea(
            child: Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const GlitchText(
                    text: 'THE PURGE',
                    style: TextStyle(
                      fontSize: 48,
                      fontWeight: FontWeight.bold,
                      letterSpacing: 4,
                      color: Colors.white,
                    ),
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Sự kiện ghép đôi hỗn loạn',
                    style: TextStyle(
                      color: Colors.white70,
                      fontSize: 16,
                    ),
                  ),
                  const SizedBox(height: 60),
                  
                  // Radar pulse
                  const SizedBox(
                    height: 250,
                    child: RadarPulse(
                      size: 250,
                      color: AppTheme.neonRed,
                    ),
                  ),
                  
                  const SizedBox(height: 40),
                  const Text(
                    '👁️ 1,247 người đang chờ ghép đôi',
                    style: TextStyle(
                      color: Colors.white54,
                      fontSize: 14,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                  
                  const SizedBox(height: 40),
                  
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 40),
                    child: SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        onPressed: () {
                          // Navigate to chat
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: AppTheme.primaryColor,
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                          elevation: 10,
                          shadowColor: AppTheme.primaryColor.withOpacity(0.5),
                        ),
                        child: const Text(
                          '⚡ THAM GIA NGAY',
                          style: TextStyle(
                            fontSize: 18,
                            fontWeight: FontWeight.bold,
                            color: Colors.white,
                          ),
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
