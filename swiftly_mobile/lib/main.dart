import 'package:flutter/material.dart';
import 'package:flutter_acrylic/window_effect.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/domain/models/label_item.dart';
import 'package:swiftly_mobile/providers/card_notifier_provider.dart';
import 'package:swiftly_mobile/providers/label_notifier_provider.dart';
import 'package:swiftly_mobile/providers/user_notifier_provider.dart';
import 'package:swiftly_mobile/routing/router.dart';
import 'package:window_manager/window_manager.dart';

import 'package:flutter_acrylic/flutter_acrylic.dart' as acrylic;

import 'package:flutter/foundation.dart' show kIsWeb;
import 'domain/user/models/user.dart';
import 'providers/current_user_provider.dart';
import 'ui/core/themes/colors.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  await acrylic.Window.initialize();
  await acrylic.Window.setEffect(
    effect: WindowEffect.hudWindow,
    color: Colors.transparent,
  );

  // acrylic.Window.makeWindowFullyTransparent();

  if (!kIsWeb) {
    try {
      await windowManager.ensureInitialized();
      await windowManager.setMinimumSize(const Size(700, 500));
    } catch (e) {
      debugPrint('Ошибка инициализации windowManager: $e');
    }
  }
  runApp(const ProviderScope(child: MyApp()));
}

class MyApp extends ConsumerWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(userNotifierProvider.notifier).loadUsers();
      ref.read(cardNotifierProvider.notifier).loadCards(ref);
      // ref.read(labelNotifierProvider.notifier).loadLabels();
      final user = User.create(
        id: 'aaa',
        name: 'Иван',
        image:
            'https://givotniymir.ru/wp-content/uploads/2016/05/enot-poloskun-obraz-zhizni-i-sreda-obitaniya-enota-poloskuna-1.jpg',
        role: LabelItem.create(cardId: '1', userId: 'a', title: 'flutter', color: AppColors.amaranthMagenta),
      );
      ref.read(userNotifierProvider.notifier).addUser(user);
      ref.read(currentUserProvider.notifier).state = user;
    });
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
