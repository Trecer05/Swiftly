import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../themes/colors.dart';
import '../../../themes/theme.dart';
import 'nav_item.dart';

class CustomNavigationBar extends StatefulWidget {
  final Widget child;

  const CustomNavigationBar({super.key, required this.child});

  @override
  State<CustomNavigationBar> createState() => _CustomNavigationBarState();
}

class _CustomNavigationBarState extends State<CustomNavigationBar> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Column(
        children: [
          Expanded(child: widget.child),
          Container(
            color: AppColors.white15,
            padding: const EdgeInsets.fromLTRB(16, 8, 16, 34),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                ...NavItem.values
                    .where((item) => item != NavItem.settings)
                    .map((item) => NavItemWidget(item: item)),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class NavItemWidget extends StatelessWidget {
  final NavItem? item;
  final IconData? icon;

  const NavItemWidget({super.key, this.item, this.icon})
    : assert(item != null || icon != null);

  @override
  Widget build(BuildContext context) {
    final isSelected = _isCurrentRoute(context, item!.route);
    return GestureDetector(
      behavior: HitTestBehavior.opaque,
      onTap: () => _handleTap(context),
      child: Container(
        padding: const EdgeInsets.all(8),
        margin: const EdgeInsets.only(bottom: 4),
        child: Column(
          children: [
            Icon(
              item?.icon ?? icon,
              color: isSelected ? AppColors.white : AppColors.unselectedItemMobile,
            ),
            if (item != null) ...[
              const SizedBox(height: 8),
              Text(
                item!.label,
                style: isSelected ? AppTextStyles.style16 : AppTextStyles.style15,
              ),
            ],
          ],
        ),
      ),
    );
  }

  bool _isCurrentRoute(BuildContext context, String route) {
    return GoRouterState.of(context).uri.toString() == route;
  }

  void _handleTap(BuildContext context) {
    GoRouter.of(context).go(item!.route);
  }
}