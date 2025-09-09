import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../themes/colors.dart';
import 'nav_item.dart';

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
                color: AppColors.transparent168,
                width: state.showLabels ? expandedWidth : null,
                padding: const EdgeInsets.symmetric(
                  vertical: 10,
                  horizontal: 10,
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

class NavItemWidget extends StatefulWidget {
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
  State<NavItemWidget> createState() => _NavItemWidgetState();
}

class _NavItemWidgetState extends State<NavItemWidget> {
  bool isHovered = false;

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<NavigationState>(
      valueListenable: widget.stateNotifier,
      builder: (_, state, __) {
        final isSelected = !widget.isController &&
            widget.item != null &&
            _isCurrentRoute(context, widget.item!.route);

        return MouseRegion(
          cursor: SystemMouseCursors.click,
          onEnter: (_) => setState(() => isHovered = true),
          onExit: (_) => setState(() => isHovered = false),
          child: GestureDetector(
            onTap: () => _handleTap(context),
            child: Container(
              padding: const EdgeInsets.all(8),
              margin: const EdgeInsets.only(bottom: 4),
              decoration: BoxDecoration(
                gradient: isSelected ? AppColors.gradient_4 : null,
                color: isHovered && !widget.isController && !isSelected ? AppColors.white15 : null,
                borderRadius: BorderRadius.circular(12),
              ),
              child: Row(
                children: [
                  Icon(
                    widget.item?.icon ?? widget.icon,
                    color:
                        isSelected ? AppColors.white : AppColors.white128,
                  ),
                  if (state.showLabels && widget.item != null) ...[
                    const SizedBox(width: 8),
                    Text(
                      widget.item!.label,
                      style: TextStyle(
                        color: isSelected
                            ? AppColors.white
                            : AppColors.white128,
                      ),
                    ),
                  ],
                ],
              ),
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
    if (widget.isController) {
      widget.stateNotifier.value = NavigationState(
        showLabels: !widget.stateNotifier.value.showLabels,
      );
    } else if (widget.item != null) {
      GoRouter.of(context).go(widget.item!.route);
    }
  }
}
