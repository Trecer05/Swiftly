import 'package:flutter/material.dart';

import '../../../../domain/kanban/models/priority.dart';
import '../../../core/themes/colors.dart';

class PrioritySelector extends StatelessWidget {
  final Priority selected;
  final void Function(Priority) onSelected;

  const PrioritySelector({
    super.key,
    required this.selected,
    required this.onSelected,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: Priority.values.map((priority) {
        final isSelected = selected == priority;
        return GestureDetector(
          onTap: () => onSelected(priority),
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 5),
            margin: const EdgeInsets.only(right: 8),
            decoration: BoxDecoration(
              color: isSelected ? priority.color.withValues(alpha: 0.2) : AppColors.white15,
              borderRadius: BorderRadius.circular(10),
              border: isSelected
                  ? Border.all(color: priority.color)
                  : null,
            ),
            child: Row(
              children: [
                Text(
                  priority.title,
                  style: TextStyle(
                    color: priority.color,
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(width: 5),
                Icon(Icons.network_cell_outlined, size: 12, color: priority.color),
              ],
            ),
          ),
        );
      }).toList(),
    );
  }
}
