import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter_acrylic/window_effect.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/providers/card_notifier_provider.dart';
import 'package:swiftly_mobile/providers/user_notifier_provider.dart';
import 'package:swiftly_mobile/routing/router.dart';
import 'package:window_manager/window_manager.dart';

import 'package:flutter_acrylic/flutter_acrylic.dart' as acrylic;

import 'package:flutter/foundation.dart' show kIsWeb, defaultTargetPlatform;
import 'dart:io' show Platform;

import 'domain/user/models/user.dart';
import 'providers/current_user_provider.dart';
import 'ui/core/themes/colors.dart';

bool isDesktopPlatform() {
  if (kIsWeb) return false;
  switch (defaultTargetPlatform) {
    case TargetPlatform.macOS:
    case TargetPlatform.windows:
    case TargetPlatform.linux:
      return true;
    default:
      return false;
  }
}

bool isWindows11() {
  if (!Platform.isWindows) return false;
  final version = Platform.operatingSystemVersion;
  return version.contains(RegExp(r'Build 2[2-9]\d{3}'));
}

class AcrylicEffectLifecycleObserver extends WidgetsBindingObserver {
  @override
  void didChangeAppLifecycleState(AppLifecycleState state) async {
    if (Platform.isWindows) {
      if (state == AppLifecycleState.inactive || state == AppLifecycleState.paused) {
        await acrylic.Window.setEffect(
          effect: isWindows11() ? WindowEffect.acrylic : WindowEffect.aero,
          color: const Color(0x88000000),
        );
      } else if (state == AppLifecycleState.resumed) {
        await acrylic.Window.setEffect(
          effect: isWindows11() ? WindowEffect.acrylic : WindowEffect.aero,
          color: const Color(0xCC000000),
        );
      }
    }
  }
}

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  // команда для дебага ui
  debugPaintSizeEnabled = false;

  if (isDesktopPlatform()) {
    try {
      await windowManager.ensureInitialized();

      await acrylic.Window.initialize();
      WindowOptions windowOptions = const WindowOptions(
        size: Size(1500, 1000),
        minimumSize: Size(1024, 500),
        backgroundColor: Colors.transparent,
        
        skipTaskbar: false,
        titleBarStyle: TitleBarStyle.normal, 
      );

      windowManager.waitUntilReadyToShow(windowOptions, () async {
        
        if (Platform.isWindows) {
          await acrylic.Window.setEffect(
            effect: isWindows11() ? WindowEffect.acrylic : WindowEffect.aero,
            color: const Color(0xCC000000),
          );
        } else {
          await acrylic.Window.setEffect(
            effect: WindowEffect.hudWindow,
            color: Colors.transparent,
          );
        }

        await windowManager.show();
        await windowManager.focus();
      });

    } catch (e) {
      debugPrint('ERROR: $e');
    }
  }

  WidgetsBinding.instance.addObserver(AcrylicEffectLifecycleObserver());
  runApp(const ProviderScope(child: MyApp()));
}

class MyApp extends ConsumerWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final bool onDesktop = isDesktopPlatform();

    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(userNotifierProvider.notifier).loadUsers();
      ref.read(cardNotifierProvider.notifier).loadCards(ref);
      final user = User.create(
        id: 'aaa',
        name: 'Иван',
        image:
            'https://givotniymir.ru/wp-content/uploads/2016/05/enot-poloskun-obraz-zhizni-i-sreda-obitaniya-enota-poloskuna-1.jpg',
      );
      ref.read(userNotifierProvider.notifier).addUser(user);
      ref.read(currentUserProvider.notifier).state = user;
    });

    return MaterialApp.router(
      debugShowCheckedModeBanner: false,
      title: 'Swiftly',
      theme: ThemeData(
        scaffoldBackgroundColor: onDesktop ? Colors.transparent : Colors.black,
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
