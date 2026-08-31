import 'package:flutter/material.dart';
import 'dart:math';

class GlitchText extends StatefulWidget {
  final String text;
  final TextStyle style;

  const GlitchText({
    super.key,
    required this.text,
    required this.style,
  });

  @override
  State<GlitchText> createState() => _GlitchTextState();
}

class _GlitchTextState extends State<GlitchText> with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  final Random _random = Random();

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      duration: const Duration(milliseconds: 200),
      vsync: this,
    )..repeat(reverse: true);
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _controller,
      builder: (context, child) {
        final double offset1 = _random.nextDouble() * 4 - 2;
        final double offset2 = _random.nextDouble() * 4 - 2;
        final double opacity = _random.nextDouble() * 0.5 + 0.5;

        return Stack(
          alignment: Alignment.center,
          children: [
            Transform.translate(
              offset: Offset(offset1, 0),
              child: Opacity(
                opacity: opacity,
                child: Text(
                  widget.text,
                  style: widget.style.copyWith(color: Colors.red),
                ),
              ),
            ),
            Transform.translate(
              offset: Offset(offset2, 0),
              child: Opacity(
                opacity: opacity,
                child: Text(
                  widget.text,
                  style: widget.style.copyWith(color: Colors.cyan),
                ),
              ),
            ),
            Text(
              widget.text,
              style: widget.style,
            ),
          ],
        );
      },
    );
  }
}
