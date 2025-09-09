import 'package:flutter/material.dart';

import '../../../../domain/kanban/models/priority.dart';
import '../../themes/colors.dart';
import '../../themes/theme.dart';

class PriorityWidget extends StatelessWidget {
  final Priority priority;

  String get _title {
    switch (priority) {
      case Priority.low:
        return 'Низкий';
      case Priority.medium:
        return 'Средний';
      case Priority.high:
        return 'Высокий';
    }
  }

  Color get _textColor {
    switch (priority) {
      case Priority.low:
        return AppColors.green1;
      case Priority.medium:
        return AppColors.yellow1;
      case Priority.high:
        return AppColors.red1;
    }
  }

  Color get _backgroundColor {
    switch (priority) {
      case Priority.low:
        return AppColors.green2;
      case Priority.medium:
        return AppColors.yellow2;
      case Priority.high:
        return AppColors.red2;
    }
  }

  const PriorityWidget({super.key, required this.priority});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(5),
      decoration: BoxDecoration(
        color: _backgroundColor,
        borderRadius: BorderRadius.circular(15),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            _title,
            style: TextStyle(
              color: _textColor,
              fontSize: AppFontSizes.size12,
              fontWeight: AppFontWeights.bolt500,
            ),
          ),
          const SizedBox(width: 5),
          Icon(Icons.network_cell_outlined, color: _textColor, size: 12),
        ],
      ),
    );
  }
}
