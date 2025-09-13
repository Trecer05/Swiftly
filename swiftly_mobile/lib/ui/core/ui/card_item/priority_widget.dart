import 'package:flutter/material.dart';

import '../../../../domain/kanban/models/priority.dart';
import '../../themes/theme.dart';

class PriorityWidget extends StatelessWidget {
  final Priority priority;

  const PriorityWidget({super.key, required this.priority});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(5),
      decoration: BoxDecoration(
        color: priority.color.withValues(alpha: 0.2),
        borderRadius: BorderRadius.circular(15),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            priority.title,
            style: TextStyle(
              color: priority.color,
              fontSize: AppFontSizes.size12,
              fontWeight: AppFontWeights.bolt500,
            ),
          ),
          const SizedBox(width: 5),
          Icon(Icons.network_cell_outlined, color: priority.color, size: 12),
        ],
      ),
    );
  }
}
