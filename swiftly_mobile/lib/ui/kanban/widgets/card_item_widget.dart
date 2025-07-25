import 'package:flutter/material.dart';

import '../../core/themes/colors.dart';
import '../../core/themes/theme.dart';
import '../../../domain/kanban/models/card_item.dart';
import '../../core/ui/label_item_widget.dart';
import 'card_details_widgets/cart_details_widget.dart';
import 'date_widget.dart';
import 'priority_widget.dart';

class CardItemWidget extends StatelessWidget {
  final CardItem card;
  const CardItemWidget({super.key, required this.card});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: (){
        showDialog(
      context: context,
      builder: (_) => CartDetailsWidget(card: card),
    );
      },
      child: Container(
        width: 300,
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
                if (card.category != null) LabelItemWidget(
                  labelItem: card.category!,
                ),
                const Spacer(),
                DateWidget(date: card.createdAt, color: AppColors.white15,),
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

