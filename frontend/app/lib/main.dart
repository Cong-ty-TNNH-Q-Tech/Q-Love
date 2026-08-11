import 'package:flutter/material.dart';

void main() {
  runApp(const QLoveApp());
}

class QLoveApp extends StatelessWidget {
  const QLoveApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Q-Love',
      theme: ThemeData(
        brightness: Brightness.dark,
        primarySwatch: Colors.pink,
        useMaterial3: true,
      ),
      home: const Scaffold(
        body: Center(
          child: Text(
            'Q-Love App Foundation',
            style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
          ),
        ),
      ),
    );
  }
}
