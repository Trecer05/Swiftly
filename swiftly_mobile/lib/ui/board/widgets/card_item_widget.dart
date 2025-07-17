import 'package:flutter/material.dart';

import '../../core/themes/colors.dart';
import '../../core/themes/theme.dart';
import 'card_item.dart';

class CardItemWidget extends StatelessWidget {
  final CardItem card;
  final VoidCallback onDelete;
  const CardItemWidget({super.key, required this.card, required this.onDelete});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onDelete,
      child: Container(
        padding: const EdgeInsets.all(10),
        decoration: BoxDecoration(
          color: AppColors.white15,
          borderRadius: BorderRadius.circular(15),
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                CategoryWidget(
                  name: card.category,
                  color: AppColors.amaranthMagenta,
                ),
                const Spacer(),
                DataWidget(date: card.createdAt),
              ],
            ),
            const SizedBox(height: 5),
            Text(card.title, style: AppTextStyles.style8),
            const SizedBox(height: 5),
            Text(card.description, style: AppTextStyles.style9),
            const SizedBox(height: 5),
            PriorityWidget(priority: card.priority),
          ],
        ),
      ),
    );
  }
}

class DataWidget extends StatelessWidget {
  final DateTime date;
  const DataWidget({super.key, required this.date});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(5),
      decoration: BoxDecoration(
        color: AppColors.white15,
        borderRadius: BorderRadius.circular(15),
      ),
      child: Text(date.toShortDate(), style: AppTextStyles.style9),
    );
  }
}

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
