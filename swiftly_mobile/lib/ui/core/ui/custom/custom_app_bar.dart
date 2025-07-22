import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/core/themes/theme.dart';

import '../../themes/colors.dart';

class CustomAppBar extends StatelessWidget {
  final String title;
  final int? quantity;
  final List<Widget> buttons;
  const CustomAppBar({super.key, required this.title, this.quantity, required this.buttons});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.only(right: 16, left: 16, bottom: 16, top: 48),
      color: AppColors.black168,
      width: double.infinity,
      height: 104,
      child: Row(
        children: [
          Stack(
            clipBehavior: Clip.none,
            children: [
              Text(title, style: AppTextStyles.style6),
              if (quantity != null && quantity! > 0) Positioned(
                right: -30,
                child: Text('($quantity)', style: AppTextStyles.style14,)),
            ],
          ),
          const Spacer(),
          ...buttons.expand((button) => [button, const SizedBox(width: 5)]).toList()..removeLast(),
        ],
      ),
    );
  }
}