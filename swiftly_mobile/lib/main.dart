import 'package:flutter/material.dart';
import 'package:flutter_acrylic/window_effect.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/routing/router.dart';
import 'package:window_manager/window_manager.dart';

import 'package:flutter_acrylic/flutter_acrylic.dart' as acrylic;

import 'ui/core/themes/colors.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  await windowManager.ensureInitialized();
  await windowManager.setMinimumSize(const Size(700, 500));

  await acrylic.Window.initialize();
  await acrylic.Window.setEffect(
    effect: WindowEffect.transparent
  );

  runApp(
    const ProviderScope(
      child: 
      // MaterialApp(
      //   home: Scaffold(
      //     body: Stack(
      //       children: [
      //         Positioned.fill(
      //           child: Image.asset('assets/vk_logo.png', fit: BoxFit.cover),
      //         ),
              MyApp(),
      //       ],
      //     ),
      //   ),
      // ),
    ),
  );
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      debugShowCheckedModeBanner: false,
      title: 'Swiftly',
      theme: ThemeData(
        scaffoldBackgroundColor: Colors.transparent,
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
