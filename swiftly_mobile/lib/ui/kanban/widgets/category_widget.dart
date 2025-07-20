import 'package:flutter/material.dart';

import '../../core/themes/theme.dart';

class CategoryWidget extends StatelessWidget {
  final String name;
  final Color color;
  const CategoryWidget({super.key, required this.name, required this.color});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(5),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.2),
        borderRadius: BorderRadius.circular(15),
      ),
      child: Text(
        name,
        style: TextStyle(
          color: color,
          fontSize: AppFontSizes.size12,
          fontWeight: AppFontWeights.bolt500,
        ),
      ),
    );
  }
}