import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/routing/router.dart';
import 'package:window_manager/window_manager.dart';

import 'ui/core/themes/colors.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  await windowManager.ensureInitialized();
  await windowManager.setMinimumSize(const Size(700, 500));
  runApp(const ProviderScope(child: MyApp()));
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'Sloyka bakery',
      theme: ThemeData(
        // scaffoldBackgroundColor: Colors.transparent,
        navigationRailTheme: const NavigationRailThemeData(
          selectedIconTheme: IconThemeData(color: AppColors.white),
          selectedLabelTextStyle: TextStyle(color: AppColors.white),
          unselectedIconTheme: IconThemeData(color: AppColors.white128),
          unselectedLabelTextStyle: TextStyle(color: AppColors.white128),
        ),
        textSelectionTheme: TextSelectionThemeData(cursorColor: Colors.white),
      ),
      routerConfig: router,
    );
  }
}
