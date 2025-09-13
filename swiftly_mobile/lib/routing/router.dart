import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:swiftly_mobile/routing/routers.dart';
import 'package:swiftly_mobile/ui/auth/widgets/auth_screen.dart';
import 'package:swiftly_mobile/ui/chat/chat.dart';
import 'package:swiftly_mobile/ui/core/ui/custom/navigation/custom_navigation_bar.dart';
import 'package:swiftly_mobile/ui/verify_code/widgets/verify_code_screen.dart';
import 'package:swiftly_mobile/ui/core/ui/custom/navigation/custom_navigation_rail.dart';
import 'package:swiftly_mobile/ui/home/widgets/home_screen.dart';
import 'package:swiftly_mobile/ui/code/code_screen.dart';
import 'package:swiftly_mobile/ui/kanban/widgets/kanban_screen.dart';
import 'package:swiftly_mobile/ui/cloud/widgets/cloud_screen.dart';
import 'package:swiftly_mobile/ui/figma/widgets/figma_screen.dart';
import 'package:swiftly_mobile/ui/settings/widgets/settings_screen.dart';

import '../utils/responsive_layout.dart';

final GoRouter router = GoRouter(
  routes: <RouteBase>[
    ShellRoute(
      builder: (BuildContext context, GoRouterState state, Widget child) {
        return ResponsiveLayout(
          mobile: CustomNavigationBar(child: child),
          desktop: CustomNavigationRail(child: child),
        );
      },
      routes: [
        GoRoute(
          path: Routers.home,
          builder: (BuildContext context, GoRouterState state) {
            return const HomeScreen();
          },
        ),
        GoRoute(
          path: Routers.chat,
          builder: (BuildContext context, GoRouterState state) {
            return const ChatScreen();
          },
        ),
        GoRoute(
          path: Routers.code,
          builder: (BuildContext context, GoRouterState state) {
            return const CodeScreen();
          },
        ),
        GoRoute(
          path: Routers.cloud,
          builder: (BuildContext context, GoRouterState state) {
            return const CloudScreen();
          },
        ),
        GoRoute(
          path: Routers.figma,
          builder: (BuildContext context, GoRouterState state) {
            return const FigmaScreen();
          },
        ),
        GoRoute(
          path: Routers.board,
          builder: (BuildContext context, GoRouterState state) {
            return const KanbanScreen();
          },
        ),
        GoRoute(
          path: Routers.settings,
          builder: (BuildContext context, GoRouterState state) {
            return const SettingsScreen();
          },
        ),
      ],
    ),
    GoRoute(
      path: Routers.auth,
      builder: (BuildContext context, GoRouterState state) {
        return const AuthScreen();
      },
    ),
    GoRoute(
      path: Routers.login,
      builder: (BuildContext context, GoRouterState state) {
        return const Placeholder();
      },
    ),
    GoRoute(
      path: Routers.register,
      builder: (BuildContext context, GoRouterState state) {
        return const Placeholder();
      },
    ),
    GoRoute(
      path: Routers.verifyCode,
      builder: (BuildContext context, GoRouterState state) {
        return const VerifyCodeScreen();
      },
    ),
  ],
);
