import 'package:flutter/material.dart' hide Card;
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/providers/card_notifier_provider.dart';
import 'package:swiftly_mobile/routing/router.dart';
import 'package:window_manager/window_manager.dart';

import 'ui/core/themes/colors.dart';
void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  await windowManager.ensureInitialized();
  await windowManager.setMinimumSize(const Size(700, 500));
  runApp(const ProviderScope(child: MyApp()));
}

class MyApp extends ConsumerWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(cardNotifierProvider.notifier).loadCarts();
    });
    return MaterialApp.router(
      title: 'Swiftly',
      theme: ThemeData(
        // scaffoldBackgroundColor: Colors.transparent,
        navigationRailTheme: const NavigationRailThemeData(
          selectedIconTheme: IconThemeData(color: AppColors.white),
          selectedLabelTextStyle: TextStyle(color: AppColors.white),
          unselectedIconTheme: IconThemeData(color: AppColors.white128),
          unselectedLabelTextStyle: TextStyle(color: AppColors.white128),
        ),
        textSelectionTheme: const TextSelectionThemeData(
          cursorColor: Colors.white,
        ),
      ),
      routerConfig: router,
    );
  }
}
