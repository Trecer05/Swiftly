import 'package:flutter/material.dart';

import '../../../../domain/kanban/models/card_item.dart';
import '../../../../domain/models/label_item.dart';
import '../../themes/colors.dart';
import '../../themes/theme.dart';
import '../label_item_widget.dart';
import 'date_widget.dart';
import 'priority_widget.dart';

class CardItemContent extends StatelessWidget {
  final CardItem card;
  final List<LabelItem> labels;

  const CardItemContent({
    super.key,
    required this.card,
    required this.labels,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            ...labels.map(
              (label) => Padding(
                padding: const EdgeInsets.only(right: 8.0),
                child: LabelItemWidget(labelItem: label),
              ),
            ),
            const Spacer(),
            DateWidget(date: card.createdAt, color: AppColors.white15),
          ],
        ),
        const SizedBox(height: 5),
        Text(card.title, style: AppTextStyles.style8, overflow: TextOverflow.ellipsis,),
        const SizedBox(height: 5),
        Text(card.description, style: AppTextStyles.style9, overflow: TextOverflow.ellipsis,),
        const SizedBox(height: 5),
        PriorityWidget(priority: card.priority),
      ],
    );
  }
}
