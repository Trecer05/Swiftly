import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../providers/label_notifier_provider.dart';
import '../../core/themes/colors.dart';
import '../../core/themes/theme.dart';
import '../../../domain/kanban/models/card_item.dart';
import '../../core/ui/label_item_widget.dart';
import 'card_details_widgets/cart_details_widget.dart';
import 'date_widget.dart';
import 'priority_widget.dart';

class CardItemWidget extends ConsumerStatefulWidget {
  final CardItem card;
  const CardItemWidget({super.key, required this.card});

  @override
  ConsumerState<CardItemWidget> createState() => _CardItemWidgetState();
}

class _CardItemWidgetState extends ConsumerState<CardItemWidget> {
  bool isHovered = false;

  @override
  Widget build(BuildContext context) {
    final labels =
        ref
            .watch(labelNotifierProvider)
            .labels
            .where((label) => label.cardId == widget.card.id)
            .toList();
    return GestureDetector(
      onTap: () {
        showDialog(
          context: context,
          builder:
              (context) => Dialog(
                insetPadding: const EdgeInsets.all(16),
                backgroundColor: Colors.transparent,
                child: CartDetailsWidget(card: widget.card),
              ),
        );
      },
      child: MouseRegion(
        onEnter: (_) => setState(() => isHovered = true),
        onExit: (_) => setState(() => isHovered = false),
        child: Container(
          width: 300,
          padding: const EdgeInsets.all(10),
          decoration: BoxDecoration(
            color: AppColors.white15,
            borderRadius: BorderRadius.circular(15),
            border: Border.all(
              color: isHovered ? AppColors.white : AppColors.transparent,
            ),
          ),
          child: Column(
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
                  DateWidget(date: widget.card.createdAt, color: AppColors.white15)
                ],
              ),
              const SizedBox(height: 5),
              Text(widget.card.title, style: AppTextStyles.style8),
              const SizedBox(height: 5),
              Text(widget.card.description, style: AppTextStyles.style9),
              const SizedBox(height: 5),
              if (widget.card.priority != null)
                PriorityWidget(priority: widget.card.priority!),
            ],
          ),
        ),
      ),
    );
  }
}
