import 'package:flutter/material.dart';
import 'package:qlove/core/theme/app_theme.dart';

class CandleData {
  final double open;
  final double high;
  final double low;
  final double close;

  CandleData({
    required this.open,
    required this.high,
    required this.low,
    required this.close,
  });
}

class CandlestickChart extends StatelessWidget {
  final List<CandleData> data;
  final double height;
  final double width;

  const CandlestickChart({
    super.key,
    required this.data,
    this.height = 200,
    this.width = double.infinity,
  });

  @override
  Widget build(BuildContext context) {
    if (data.isEmpty) {
      return SizedBox(
        height: height,
        width: width,
        child: const Center(child: Text("No data available")),
      );
    }
    
    return Container(
      height: height,
      width: width,
      decoration: BoxDecoration(
        color: AppTheme.scaffoldBackground.withOpacity(0.5),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: Colors.white12),
      ),
      child: CustomPaint(
        painter: _CandlestickPainter(
          data: data,
          bullColor: AppTheme.neonGreen,
          bearColor: AppTheme.neonRed,
        ),
      ),
    );
  }
}

class _CandlestickPainter extends CustomPainter {
  final List<CandleData> data;
  final Color bullColor;
  final Color bearColor;

  _CandlestickPainter({
    required this.data,
    required this.bullColor,
    required this.bearColor,
  });

  @override
  void paint(Canvas canvas, Size size) {
    if (data.isEmpty) return;

    double minPrice = data.first.low;
    double maxPrice = data.first.high;

    for (var candle in data) {
      if (candle.low < minPrice) minPrice = candle.low;
      if (candle.high > maxPrice) maxPrice = candle.high;
    }

    final double priceRange = maxPrice - minPrice;
    // Prevent division by zero if all prices are the same
    final double safeRange = priceRange == 0 ? 1.0 : priceRange;
    
    final double candleWidth = size.width / data.length;
    final double bodyWidth = candleWidth * 0.6;
    final double bodyPadding = (candleWidth - bodyWidth) / 2;

    final Paint bullPaint = Paint()
      ..color = bullColor
      ..style = PaintingStyle.fill;
    
    final Paint bullLinePaint = Paint()
      ..color = bullColor
      ..strokeWidth = 1.5
      ..style = PaintingStyle.stroke;

    final Paint bearPaint = Paint()
      ..color = bearColor
      ..style = PaintingStyle.fill;
      
    final Paint bearLinePaint = Paint()
      ..color = bearColor
      ..strokeWidth = 1.5
      ..style = PaintingStyle.stroke;

    for (int i = 0; i < data.length; i++) {
      final candle = data[i];
      final isBull = candle.close >= candle.open;
      
      final paint = isBull ? bullPaint : bearPaint;
      final linePaint = isBull ? bullLinePaint : bearLinePaint;

      final leftX = i * candleWidth + bodyPadding;
      final centerX = leftX + (bodyWidth / 2);

      final topY = size.height - ((candle.high - minPrice) / safeRange * size.height);
      final bottomY = size.height - ((candle.low - minPrice) / safeRange * size.height);
      
      final openY = size.height - ((candle.open - minPrice) / safeRange * size.height);
      final closeY = size.height - ((candle.close - minPrice) / safeRange * size.height);
      
      final bodyTop = isBull ? closeY : openY;
      final bodyBottom = isBull ? openY : closeY;

      canvas.drawLine(
        Offset(centerX, topY),
        Offset(centerX, bottomY),
        linePaint,
      );

      final rectHeight = (bodyBottom - bodyTop).clamp(1.0, double.infinity);
      canvas.drawRect(
        Rect.fromLTWH(leftX, bodyTop, bodyWidth, rectHeight),
        paint,
      );
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => true;
}
