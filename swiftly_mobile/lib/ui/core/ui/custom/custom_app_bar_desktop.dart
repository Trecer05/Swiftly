import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/core/themes/colors.dart';
import 'package:swiftly_mobile/ui/core/themes/theme.dart';

class CustomAppBarDesktop extends StatelessWidget {
  final String title;
  final int? quantity;
  final List<Widget> buttons;
  const CustomAppBarDesktop({
    super.key,
    required this.title,
    this.quantity,
    required this.buttons,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.only(right: 16, left: 16, bottom: 16, top: 48),
      width: double.infinity,
      // height: 104,
      child: Row(
        children: [
          Stack(
            clipBehavior: Clip.none,
            children: [
              Text(title, style: AppTextStyles.style6),
              if (quantity != null && quantity! > 0)
                Positioned(
                  right: -30,
                  child: Text('($quantity)', style: AppTextStyles.style14),
                ),
            ],
          ),
          const Spacer(),
          ...buttons
              .expand((button) => [button, const SizedBox(width: 5)])
              .toList()
            ..removeLast(),
        ],
      ),
    );
  }
}

class SegmentedControlWidgetDesktop extends StatefulWidget {
  const SegmentedControlWidgetDesktop({super.key});

  @override
  State<SegmentedControlWidgetDesktop> createState() =>
      _SegmentedControlWidgetDesktopState();
}

class _SegmentedControlWidgetDesktopState
    extends State<SegmentedControlWidgetDesktop> {
  int selectedIndex = 0;

  void _onSelect(int index) {
    setState(() {
      selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    final items = ['Фильтры', 'Сортировать', 'Мои задачи'];

    return Container(
      padding: const EdgeInsets.all(4),
      height: 40,
      decoration: BoxDecoration(
        color: AppColors.white15,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          for (int i = 0; i < items.length; i++) ...[
            GestureDetector(
              onTap: () => _onSelect(i),
              child: SegmentWidgetDesktop(
                title: items[i],
              ),
            ),
            if (i < items.length - 1)
              Container(margin: const EdgeInsets.symmetric(horizontal: 8), width: 1, height: 16, color: AppColors.white128),
          ],
        ],
      ),
    );
  }
}

class SegmentWidgetDesktop extends StatefulWidget {
  final String title;
  const SegmentWidgetDesktop({
    super.key,
    required this.title,
  });

  @override
  State<SegmentWidgetDesktop> createState() => _SegmentWidgetDesktopState();
}

class _SegmentWidgetDesktopState extends State<SegmentWidgetDesktop> {
  bool isHovered = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => isHovered = true),
      onExit: (_) => setState(() => isHovered = false),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        decoration: BoxDecoration(
          color: isHovered ? AppColors.white15 : null,
          borderRadius: BorderRadius.circular(8),
        ),
        child: Text(
          widget.title,
          style: AppTextStyles.style17.copyWith(height: 1.0),
          overflow: TextOverflow.ellipsis,
        ),
      ),
    );
  }
}
