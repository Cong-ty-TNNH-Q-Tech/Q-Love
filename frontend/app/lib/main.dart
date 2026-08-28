// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:app_links/app_links.dart';
import 'package:qlove/core/network/dio_client.dart';
import 'package:qlove/core/network/secure_storage_service.dart';
import 'package:qlove/features/auth/bloc/auth_bloc.dart';
import 'package:qlove/features/auth/bloc/auth_state.dart';
import 'package:qlove/features/auth/bloc/auth_event.dart';
import 'package:qlove/features/auth/data/auth_repository.dart';
import 'package:qlove/features/auth/ui/login_screen.dart';
import 'package:qlove/features/auth/ui/profile_creation_screen.dart';
import 'package:qlove/features/locket/data/widget_repository.dart';
import 'core/theme/app_theme.dart';
import 'core/bloc/app_bloc.dart';
import 'features/wingman/ui/matchmaker_popup.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await WidgetRepository().initialize();
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

  late final SecureStorageService _secureStorageService;
  late final AuthRepository _authRepository;
  late final AuthBloc _authBloc;
  late final DioClient _dioClient;

  @override
  void initState() {
    super.initState();
    _secureStorageService = SecureStorageService();
    
    _dioClient = DioClient(
      secureStorageService: _secureStorageService,
      getAccessToken: () => _authBloc.currentAccessToken,
      onTokenRefreshed: (token) {
        _authBloc.add(TokenRefreshed(token));
      },
      onLogoutRequired: () {
        _authBloc.add(LogoutRequested());
      },
    );

    _authRepository = AuthRepository(
      dio: _dioClient.dio,
      secureStorageService: _secureStorageService,
    );

    _authBloc = AuthBloc(authRepository: _authRepository);

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
    if (uri.scheme == 'qlove') {
      if (uri.host == 'wingman' && uri.path == '/match') {
        final targetId = uri.queryParameters['target'];
        if (targetId != null) {
          WidgetsBinding.instance.addPostFrameCallback((_) {
            if (_navigatorKey.currentContext != null) {
              MatchmakerPopup.show(_navigatorKey.currentContext!, targetId);
            }
          });
        }
      } else if (uri.host == 'match' && uri.pathSegments.isNotEmpty) {
        final token = uri.pathSegments.first;
        // Navigate or handle match token
        WidgetsBinding.instance.addPostFrameCallback((_) {
          if (_navigatorKey.currentContext != null) {
             ScaffoldMessenger.of(_navigatorKey.currentContext!).showSnackBar(
              SnackBar(content: Text('Received match token: $token')),
            );
          }
        });
      }
    }
  }

  @override
  void dispose() {
    _linkSubscription?.cancel();
    _authBloc.close();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return MultiRepositoryProvider(
      providers: [
        RepositoryProvider.value(value: _secureStorageService),
        RepositoryProvider.value(value: _authRepository),
      ],
      child: MultiBlocProvider(
        providers: [
          BlocProvider<AppBloc>(
            create: (context) => AppBloc()..add(AppInitialized()),
          ),
          BlocProvider<AuthBloc>.value(value: _authBloc),
        ],
        child: MaterialApp(
          navigatorKey: _navigatorKey,
          title: 'Q-Love',
          localizationsDelegates: const [
            AppLocalizations.delegate,
            GlobalMaterialLocalizations.delegate,
            GlobalWidgetsLocalizations.delegate,
            GlobalCupertinoLocalizations.delegate,
          ],
          supportedLocales: const [
            Locale('en', ''),
            Locale('vi', ''),
          ],
          theme: AppTheme.darkTheme,
          home: BlocBuilder<AuthBloc, AuthState>(
            builder: (context, state) {
              if (state is AuthAuthenticated) {
                return const Scaffold(
                  body: Center(
                    child: Text(
                      'Q-Love App Foundation (Home)',
                      style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
                    ),
                  ),
                );
              } else if (state is AuthProfileNeedsCreation) {
                return const ProfileCreationScreen();
              } else {
                return const LoginScreen();
              }
            },
          ),
        ),
      ),
    );
  }
}
