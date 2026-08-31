import 'package:flutter/material.dart';
import 'package:qlove/core/theme/app_theme.dart';
import 'dart:ui';
import 'dart:async';

class PurgeRevealScreen extends StatefulWidget {
  const PurgeRevealScreen({super.key});

  @override
  State<PurgeRevealScreen> createState() => _PurgeRevealScreenState();
}

class _PurgeRevealScreenState extends State<PurgeRevealScreen> {
  int _secondsLeft = 28;
  Timer? _timer;
  bool _revealed = false;

  @override
  void initState() {
    super.initState();
    _startTimer();
  }

  void _startTimer() {
    _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
      if (_secondsLeft > 0 && !_revealed) {
        setState(() {
          _secondsLeft--;
        });
      } else {
        _timer?.cancel();
      }
    });
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF080810),
      body: Stack(
        children: [
          // Blurred background image
          Positioned.fill(
            child: TweenAnimationBuilder<double>(
              tween: Tween<double>(begin: 20.0, end: _revealed ? 0.0 : 20.0),
              duration: const Duration(seconds: 2),
              builder: (context, blurValue, child) {
                return ImageFiltered(
                  imageFilter: ImageFilter.blur(sigmaX: blurValue, sigmaY: blurValue),
                  child: Container(
                    decoration: const BoxDecoration(
                      image: DecorationImage(
                        image: AssetImage('assets/images/placeholder_avatar.png'), // placeholder
                        fit: BoxFit.cover,
                      ),
                    ),
                    // If no asset is present, fallback to a solid color placeholder
                    child: Container(color: Colors.blueGrey),
                  ),
                );
              },
            ),
          ),
          
          // Overlay
          if (!_revealed)
            Positioned.fill(
              child: Container(color: Colors.black.withOpacity(0.4)),
            ),

          SafeArea(
            child: Column(
              children: [
                const SizedBox(height: 20),
                const Text(
                  'BẠN CÓ MUỐN LỘ DIỆN?',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                    letterSpacing: 2,
                  ),
                ),
                const SizedBox(height: 10),
                Text(
                  'Cả hai phải đồng ý trong $_secondsLeft giây',
                  style: const TextStyle(
                    color: Colors.white70,
                    fontSize: 16,
                  ),
                ),
                const Spacer(),
                
                if (!_revealed)
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 40),
                    child: Row(
                      children: [
                        Expanded(
                          child: ElevatedButton(
                            onPressed: () {
                              // Cancel
                              Navigator.pop(context);
                            },
                            style: ElevatedButton.styleFrom(
                              backgroundColor: Colors.white24,
                              padding: const EdgeInsets.symmetric(vertical: 16),
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(12),
                              ),
                            ),
                            child: const Text('BỎ QUA'),
                          ),
                        ),
                        const SizedBox(width: 16),
                        Expanded(
                          child: ElevatedButton(
                            onPressed: () {
                              setState(() {
                                _revealed = true;
                              });
                            },
                            style: ElevatedButton.styleFrom(
                              backgroundColor: AppTheme.neonBlue,
                              padding: const EdgeInsets.symmetric(vertical: 16),
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(12),
                              ),
                            ),
                            child: const Text('ĐỒNG Ý'),
                          ),
                        ),
                      ],
                    ),
                  ),
                  
                if (_revealed)
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 40),
                    child: SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        onPressed: () {
                          // Continue to profile
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: AppTheme.primaryColor,
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                        ),
                        child: const Text('XEM HỒ SƠ'),
                      ),
                    ),
                  )
              ],
            ),
          )
        ],
      ),
    );
  }
}
