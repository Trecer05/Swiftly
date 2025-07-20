import 'package:flutter/material.dart';

import '../../core/themes/colors.dart';
import '../../core/themes/theme.dart';
import '../../../domain/kanban/models/card_item.dart';
import 'category_widget.dart';
import 'date_widget.dart';
import 'priority_widget.dart';

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
                if (card.category != null) CategoryWidget(
                  name: card.category!,
                  color: AppColors.amaranthMagenta,
                ),
                const Spacer(),
                DateWidget(date: card.createdAt),
              ],
            ),
            const SizedBox(height: 5),
            Text(card.title, style: AppTextStyles.style8),
            const SizedBox(height: 5),
            Text(card.description, style: AppTextStyles.style9),
            const SizedBox(height: 5),
            if (card.priority != null) PriorityWidget(priority: card.priority!),
          ],
        ),
      ),
    );
  }
}

