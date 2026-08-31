import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import 'package:q_love/features/island/ui/love_island_screen.dart';
import 'package:rive/rive.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

void main() {
  testWidgets('LoveIslandScreen renders correctly with BlocProvider', (tester) async {
    await tester.pumpWidget(
      const MaterialApp(
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        home: LoveIslandScreen(),
      ),
    );

    // Wait for the UI to pump
    await tester.pumpAndSettle();

    // Verify if LoveIslandScreen exists
    expect(find.byType(LoveIslandScreen), findsOneWidget);
    
    // Verify RiveAnimation is present
    expect(find.byType(RiveAnimation), findsOneWidget);
    
    // Verify Action buttons are rendered
    expect(find.byIcon(Icons.favorite), findsOneWidget);
    expect(find.byIcon(Icons.star), findsOneWidget);
    expect(find.byIcon(Icons.home), findsOneWidget);
  });
}
