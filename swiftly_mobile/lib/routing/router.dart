import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:swiftly_mobile/routing/routers.dart';
import 'package:swiftly_mobile/ui/auth/widgets/auth_screen.dart';
import 'package:swiftly_mobile/ui/chat/chat.dart';
import 'package:swiftly_mobile/ui/cloud/file_editor_screen.dart';
import 'package:swiftly_mobile/ui/cloud/widgets/file_model.dart';
import 'package:swiftly_mobile/ui/core/ui/custom/navigation/custom_navigation_bar.dart';
import 'package:swiftly_mobile/ui/verify_code/widgets/verify_code_screen.dart';
import 'package:swiftly_mobile/ui/core/ui/custom/navigation/custom_navigation_rail.dart';
import 'package:swiftly_mobile/ui/home/widgets/home_screen.dart';
import 'package:swiftly_mobile/ui/kanban/widgets/kanban_screen.dart';
import 'package:swiftly_mobile/ui/cloud/cloud_screen.dart';
import 'package:swiftly_mobile/ui/settings/widgets/settings_screen.dart';
import 'package:swiftly_mobile/ui/chat/mob/mob_call_screen.dart';

import '../utils/responsive_layout.dart';

final GoRouter router = GoRouter(
  routes: <RouteBase>[
    GoRoute(
          path: '/call',
          builder: (context, state) {
            final String username = state.extra as String;
            return MobileCallScreen(
              username: username,
              isIncoming: false,
            );
          },
        ),
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
          path: Routers.cloud,
          name: 'cloud',
          builder: (BuildContext context, GoRouterState state) {
            return const CloudScreen();
          },
          routes: [
            GoRoute(
              path: 'path', 
              name:'fileEditor', 
              builder: (BuildContext context, GoRouterState state) {
                final file = state.extra as FileInfo;
                return FileEditorScreen(file: file);
              }
            ),
          ]
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
