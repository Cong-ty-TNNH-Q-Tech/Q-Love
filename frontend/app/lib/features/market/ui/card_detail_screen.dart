// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter/material.dart';
import 'package:qlove/core/theme/app_theme.dart';
import 'package:qlove/widgets/candlestick_chart.dart';
import 'dart:math';

class CardDetailScreen extends StatefulWidget {
  final String heroTag;

  const CardDetailScreen({super.key, required this.heroTag});

  @override
  State<CardDetailScreen> createState() => _CardDetailScreenState();
}

class _CardDetailScreenState extends State<CardDetailScreen> {
  double _tiltX = 0;
  double _tiltY = 0;

  final List<CandleData> _mockData = List.generate(
    20,
    (index) {
      final open = 120.0 + Random().nextDouble() * 20;
      final close = open + (Random().nextDouble() * 10 - 5);
      final high = max(open, close) + Random().nextDouble() * 5;
      final low = min(open, close) - Random().nextDouble() * 5;
      return CandleData(open: open, high: high, low: low, close: close);
    },
  );

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppTheme.scaffoldBackground,
      appBar: AppBar(
        title: const Text(
          "#NVA — Nguyễn Văn A",
          style: TextStyle(fontWeight: FontWeight.bold, fontSize: 18),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.share),
            onPressed: () {},
          )
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.symmetric(horizontal: 20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const SizedBox(height: 20),
            // 3D Hologram Card
            GestureDetector(
              onPanUpdate: (details) {
                setState(() {
                  _tiltX += details.delta.dy * -0.01;
                  _tiltY += details.delta.dx * 0.01;
                  _tiltX = _tiltX.clamp(-0.5, 0.5);
                  _tiltY = _tiltY.clamp(-0.5, 0.5);
                });
              },
              onPanEnd: (_) {
                setState(() {
                  _tiltX = 0;
                  _tiltY = 0;
                });
              },
              child: Transform(
                transform: Matrix4.identity()
                  ..setEntry(3, 2, 0.001)
                  ..rotateX(_tiltX)
                  ..rotateY(_tiltY),
                alignment: FractionalOffset.center,
                child: Container(
                  height: MediaQuery.of(context).size.height * 0.45,
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(24),
                    gradient: const LinearGradient(
                      begin: Alignment.topLeft,
                      end: Alignment.bottomRight,
                      colors: [Color(0xFF2C2C3E), Color(0xFF161622)],
                    ),
                    boxShadow: [
                      BoxShadow(
                        color: AppTheme.primaryColor.withOpacity(0.3),
                        blurRadius: 30,
                        offset: Offset(_tiltY * -20, _tiltX * -20),
                      )
                    ],
                    border: Border.all(
                      color: Colors.white.withOpacity(0.1),
                      width: 1,
                    ),
                  ),
                  child: Stack(
                    children: [
                      // Image Placeholder
                      ClipRRect(
                        borderRadius: BorderRadius.circular(24),
                        child: Center(
                          child: Icon(
                            Icons.person,
                            size: 150,
                            color: Colors.white.withOpacity(0.1),
                          ),
                        ),
                      ),
                      // Hologram Glare
                      Positioned.fill(
                        child: Container(
                          decoration: BoxDecoration(
                            borderRadius: BorderRadius.circular(24),
                            gradient: LinearGradient(
                              begin: Alignment(-1.0 + _tiltY, -1.0 + _tiltX),
                              end: Alignment(1.0 + _tiltY, 1.0 + _tiltX),
                              colors: [
                                Colors.white.withOpacity(0.0),
                                Colors.white.withOpacity(0.2),
                                Colors.white.withOpacity(0.0),
                              ],
                              stops: const [0.0, 0.5, 1.0],
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
            const SizedBox(height: 30),
            
            // Time Selector
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: ['1H', '1D', '1W', '1M'].map((time) {
                return Container(
                  padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                  decoration: BoxDecoration(
                    color: time == '1H' ? AppTheme.surfaceColor : Colors.transparent,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Text(
                    time,
                    style: TextStyle(
                      color: time == '1H' ? Colors.white : Colors.white54,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                );
              }).toList(),
            ),
            
            const SizedBox(height: 20),
            
            // Candlestick Chart
            CandlestickChart(data: _mockData, height: 180),
            
            const SizedBox(height: 20),
            
            // Stats
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: AppTheme.surfaceColor,
                borderRadius: BorderRadius.circular(16),
              ),
              child: Column(
                children: [
                  _buildStatRow('Giá hiện tại', '147.5 Xu'),
                  const Divider(color: Colors.white12, height: 24),
                  _buildStatRow('Biến động', '+4.3 (+3%)', valueColor: AppTheme.neonGreen),
                  const Divider(color: Colors.white12, height: 24),
                  _buildStatRow('Lượt Match', '12 (hôm nay)'),
                  const Divider(color: Colors.white12, height: 24),
                  _buildStatRow('Clan Score', '85 điểm'),
                  const Divider(color: Colors.white12, height: 24),
                  _buildStatRow('Người Sưu Tầm', '23 người'),
                ],
              ),
            ),
            
            const SizedBox(height: 30),
            
            // Action Buttons
            Row(
              children: [
                Expanded(
                  child: ElevatedButton(
                    onPressed: () {},
                    style: ElevatedButton.styleFrom(
                      backgroundColor: AppTheme.neonGreen,
                      foregroundColor: Colors.black,
                      padding: const EdgeInsets.symmetric(vertical: 16),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    child: const Text('🛒 MUA', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: ElevatedButton(
                    onPressed: () {},
                    style: ElevatedButton.styleFrom(
                      backgroundColor: AppTheme.neonRed,
                      foregroundColor: Colors.white,
                      padding: const EdgeInsets.symmetric(vertical: 16),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    child: const Text('📤 BÁN', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 40),
          ],
        ),
      ),
    );
  }

  Widget _buildStatRow(String label, String value, {Color? valueColor}) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: const TextStyle(color: Colors.white54, fontSize: 14),
        ),
        Text(
          value,
          style: TextStyle(
            color: valueColor ?? Colors.white,
            fontWeight: FontWeight.bold,
            fontSize: 14,
          ),
        ),
      ],
    );
  }
}

