import 'package:flutter/material.dart';

import '../../../domain/models/label_item.dart';
import '../themes/theme.dart';

class LabelItemWidget extends StatelessWidget {
  final LabelItem labelItem;

  const LabelItemWidget({super.key, required this.labelItem});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(5),
      decoration: BoxDecoration(
        color: labelItem.color.withValues(alpha: 0.2),
        borderRadius: BorderRadius.circular(15),
      ),
      child: Text(
        labelItem.title,
        style: TextStyle(
          color: labelItem.color,
          fontSize: AppFontSizes.size12,
          fontWeight: AppFontWeights.bolt500,
        ),
      ),
    );
  }
}