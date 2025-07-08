import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../routing/routers.dart';
import '../themes/colors.dart';

enum NavItem {
  home(Routers.home, Icons.home, 'Главная'),
  chat(Routers.chat, Icons.chat, 'Чат'),
  code(Routers.code, Icons.code, 'Код'),
  cloud(Routers.cloud, Icons.cloud, 'Облако'),
  figma(Routers.figma, Icons.pan_tool, 'Фигма'),
  board(Routers.board, Icons.task, 'Задачи'),
  settings(Routers.settings, Icons.settings, 'Настройки');

  final String route;
  final IconData icon;
  final String label;

  const NavItem(this.route, this.icon, this.label);

  static final _routeMap = {for (var item in NavItem.values) item.route: item};

  static NavItem? fromRoute(String route) => _routeMap[route];
}

class NavigationState {
  final bool showLabels;

  const NavigationState({required this.showLabels});
}

class CustomNavigationRail extends StatefulWidget {
  final Widget child;

  const CustomNavigationRail({super.key, required this.child});

  @override
  State<CustomNavigationRail> createState() => _CustomNavigationRailState();
}

class _CustomNavigationRailState extends State<CustomNavigationRail> {
  final ValueNotifier<NavigationState> _stateNotifier = ValueNotifier(
    const NavigationState(showLabels: true),
  );
  final double expandedWidth = 200;

  @override
  void dispose() {
    _stateNotifier.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Row(
        children: [
          ValueListenableBuilder<NavigationState>(
            valueListenable: _stateNotifier,
            builder: (_, state, __) {
              return Container(
                width: state.showLabels ? expandedWidth : null,
                padding: const EdgeInsets.symmetric(
                  vertical: 10,
                  horizontal: 10,
                ),
                decoration: const BoxDecoration(
                  color: Color.fromARGB(255, 3, 13, 55),
                ),
                child: Column(
                  children: [
                    NavItemWidget(
                      icon: Icons.sort,
                      isController: true,
                      stateNotifier: _stateNotifier,
                    ),
                    ...NavItem.values.where((item) => item != NavItem.settings)
                        .map((item) => NavItemWidget(
                              item: item,
                              stateNotifier: _stateNotifier,
                            )),
                    const Spacer(),
                    NavItemWidget(
                      item: NavItem.settings,
                      stateNotifier: _stateNotifier,
                    ),
                  ],
                ),
              );
            },
          ),
          Expanded(child: widget.child),
        ],
      ),
    );
  }
}

class NavItemWidget extends StatelessWidget {
  final NavItem? item;
  final IconData? icon;
  final bool isController;
  final ValueNotifier<NavigationState> stateNotifier;

  const NavItemWidget({
    super.key,
    this.item,
    this.icon,
    this.isController = false,
    required this.stateNotifier,
  }) : assert(item != null || icon != null);

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<NavigationState>(
      valueListenable: stateNotifier,
      builder: (_, state, __) {
        final isSelected =
            !isController &&
            item != null &&
            _isCurrentRoute(context, item!.route);

        return GestureDetector(
          onTap: () => _handleTap(context),
          child: Container(
            padding: const EdgeInsets.all(8),
            margin: const EdgeInsets.only(bottom: 4),
            decoration: BoxDecoration(
              gradient: isSelected ? AppColors.gradient_4 : null,
              borderRadius: BorderRadius.circular(12),
            ),
            child: Row(
              children: [
                Icon(
                  item?.icon ?? icon,
                  color: isSelected ? AppColors.white : AppColors.white128,
                ),
                if (state.showLabels && item != null) ...[
                  const SizedBox(width: 8),
                  Text(
                    item!.label,
                    style: TextStyle(
                      color: isSelected ? AppColors.white : AppColors.white128,
                    ),
                  ),
                ],
              ],
            ),
          ),
        );
      },
    );
  }

  bool _isCurrentRoute(BuildContext context, String route) {
    return GoRouterState.of(context).uri.toString() == route;
  }

  void _handleTap(BuildContext context) {
    if (isController) {
      stateNotifier.value = NavigationState(
        showLabels: !stateNotifier.value.showLabels,
      );
    } else if (item != null) {
      GoRouter.of(context).go(item!.route);
    }
  }
}
