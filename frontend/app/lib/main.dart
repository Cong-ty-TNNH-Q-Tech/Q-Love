// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:app_links/app_links.dart';
import 'core/theme/app_theme.dart';
import 'core/bloc/app_bloc.dart';
import 'features/wingman/ui/matchmaker_popup.dart';

void main() {
  runApp(const QLoveApp());
}

class QLoveApp extends StatefulWidget {
  const QLoveApp({super.key});

  @override
  State<QLoveApp> createState() => _QLoveAppState();
}

class _QLoveAppState extends State<QLoveApp> {
  final _navigatorKey = GlobalKey<NavigatorState>();
  late AppLinks _appLinks;
  StreamSubscription<Uri>? _linkSubscription;

  @override
  void initState() {
    super.initState();
    _initDeepLinks();
  }

  Future<void> _initDeepLinks() async {
    _appLinks = AppLinks();

    // Handle links when app is in cold state (killed)
    try {
      final initialUri = await _appLinks.getInitialLink();
      if (initialUri != null) {
        _handleDeepLink(initialUri);
      }
    } catch (e) {
      debugPrint("Failed to get initial deep link: $e");
    }

    // Handle links when app is in warm state (running in background/foreground)
    _linkSubscription = _appLinks.uriLinkStream.listen((uri) {
      _handleDeepLink(uri);
    }, onError: (err) {
      debugPrint("Failed to handle incoming deep link: $err");
    });
  }

  void _handleDeepLink(Uri uri) {
    if (uri.scheme == 'qlove' && uri.host == 'wingman' && uri.path == '/match') {
      final targetId = uri.queryParameters['target'];
      if (targetId != null) {
        // Wait for navigator to be ready
        WidgetsBinding.instance.addPostFrameCallback((_) {
          if (_navigatorKey.currentContext != null) {
            MatchmakerPopup.show(_navigatorKey.currentContext!, targetId);
          }
        });
      }
    }
  }

  @override
  void dispose() {
    _linkSubscription?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
      providers: [
        BlocProvider<AppBloc>(
          create: (context) => AppBloc()..add(AppInitialized()),
        ),
      ],
      child: MaterialApp(
        navigatorKey: _navigatorKey,
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
      ),
    );
  }
}
