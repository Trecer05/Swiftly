import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/core/themes/theme.dart';

import '../../themes/colors.dart';

class CustomAppBarMobile extends StatelessWidget {
  final String title;
  final int? quantity;
  final List<Widget> buttons;
  const CustomAppBarMobile({
    super.key,
    required this.title,
    this.quantity,
    required this.buttons,
  });

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: double.infinity,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
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
          const SizedBox(height: 10),
          ...buttons
              .expand((button) => [button, const SizedBox(width: 5)])
              .toList()
            ..removeLast(),
        ],
      ),
    );
  }
}

class SegmentedControlWidgetMobile extends StatefulWidget {
  const SegmentedControlWidgetMobile({super.key});

  @override
  State<SegmentedControlWidgetMobile> createState() => _SegmentedControlWidgetMobileState();
}

class _SegmentedControlWidgetMobileState extends State<SegmentedControlWidgetMobile> {
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
            // Expanded(
            //   child: 
              GestureDetector(
                onTap: () => _onSelect(i),
                child: SegmentWidgetMobile(
                  title: items[i],
                  selected: selectedIndex == i,
                ),
              ),
            // ),
            if (i < items.length - 1) const SizedBox(width: 4),
          ],
        ],
      ),
    );
  }
}

class SegmentWidgetMobile extends StatelessWidget {
  final String title;
  final bool selected;
  const SegmentWidgetMobile({super.key, required this.title, required this.selected});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 8),
      decoration: BoxDecoration(
        color: selected ? AppColors.white15 : null,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Text(
        title,
        style: AppTextStyles.style17,
        // textAlign: TextAlign.center,
        overflow: TextOverflow.ellipsis,
      ),
    );
  }
}
