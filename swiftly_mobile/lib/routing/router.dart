import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:swiftly_mobile/routing/routers.dart';
import 'package:swiftly_mobile/ui/core/ui/custom_navigation_rail.dart';
import 'package:swiftly_mobile/ui/auth/widgets/auth_screen.dart';
import 'package:swiftly_mobile/ui/verify_code/widgets/verify_code_screen.dart';

final GoRouter router = GoRouter(
  routes: <RouteBase>[
    GoRoute(
      path: Routers.auth,
      builder: (BuildContext context, GoRouterState state) {
        return AuthScreen();
      },
    ),
    GoRoute(
      path: Routers.login,
      builder: (BuildContext context, GoRouterState state) {
        return Placeholder();
      },
    ),
    GoRoute(
      path: Routers.register,
      builder: (BuildContext context, GoRouterState state) {
        return Placeholder();
      },
    ),
    GoRoute(
      path: Routers.verifyCode,
      builder: (BuildContext context, GoRouterState state) {
        return VerifyCodeScreen();
      },
    ),
    GoRoute(
      path: Routers.home,
      builder: (BuildContext context, GoRouterState state) {
        return CustomNavigationRail();
      },
    ),
    GoRoute(
      path: Routers.chat,
      builder: (BuildContext context, GoRouterState state) {
        return Placeholder();
      },
    ),
    GoRoute(
      path: Routers.code,
      builder: (BuildContext context, GoRouterState state) {
        return Placeholder();
      },
    ),
    GoRoute(
      path: Routers.cloud,
      builder: (BuildContext context, GoRouterState state) {
        return Placeholder();
      },
    ),
    GoRoute(
      path: Routers.figma,
      builder: (BuildContext context, GoRouterState state) {
        return Placeholder();
      },
    ),
    GoRoute(
      path: Routers.board,
      builder: (BuildContext context, GoRouterState state) {
        return Placeholder();
      },
    ),
  ],
);