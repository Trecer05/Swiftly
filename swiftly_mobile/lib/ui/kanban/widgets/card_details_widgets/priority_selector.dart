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
              color: isSelected ? _backgroundColor(priority) : AppColors.white15,
              borderRadius: BorderRadius.circular(10),
              border: isSelected
                  ? Border.all(color: _textColor(priority))
                  : null,
            ),
            child: Row(
              children: [
                Text(
                  _title(priority),
                  style: TextStyle(
                    color: _textColor(priority),
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(width: 5),
                Icon(Icons.network_cell_outlined, size: 12, color: _textColor(priority)),
              ],
            ),
          ),
        );
      }).toList(),
    );
  }

  String _title(Priority priority) {
    switch (priority) {
      case Priority.low:
        return 'Низкий';
      case Priority.medium:
        return 'Средний';
      case Priority.high:
        return 'Высокий';
    }
  }

  Color _textColor(Priority priority) {
    switch (priority) {
      case Priority.low:
        return AppColors.green1;
      case Priority.medium:
        return AppColors.yellow1;
      case Priority.high:
        return AppColors.red1;
    }
  }

  Color _backgroundColor(Priority priority) {
    switch (priority) {
      case Priority.low:
        return AppColors.green2;
      case Priority.medium:
        return AppColors.yellow2;
      case Priority.high:
        return AppColors.red2;
    }
  }
}
