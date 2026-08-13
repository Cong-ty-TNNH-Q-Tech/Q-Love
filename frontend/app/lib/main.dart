// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'core/theme/app_theme.dart';
import 'core/bloc/app_bloc.dart';

void main() {
  runApp(const QLoveApp());
}

class QLoveApp extends StatelessWidget {
  const QLoveApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
      providers: [
        BlocProvider<AppBloc>(
          create: (context) => AppBloc()..add(AppInitialized()),
        ),
      ],
      child: MaterialApp(
        title: 'Q-Love',
        theme: AppTheme.darkTheme,
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
